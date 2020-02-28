package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"sapling/config"
	"time"
)

//Info is the data on the product from the Tesco API, stored in the database as JSON
type Info struct {
	Ean			string `json:"ean"`
	Gtin        string `json:"gtin"`
	Tpnb        string `json:"tpnb"`
	Tpnc        string `json:"tpnc"`
	Description string `json:"description"`
	Brand       string `json:"brand"`
	QtyContents struct {
		Quantity      float64 `json:"quantity"`
		TotalQuantity float64 `json:"totalQuantity"`
		QuantityUom   string  `json:"quantityUom"`
		NetContents   string  `json:"netContents"`
	} `json:"qtyContents"`
	ProductCharacteristics struct {
		IsFood             bool   `json:"isFood"`
		IsDrink            bool   `json:"isDrink"`
		HealthScore        int    `json:"healthScore"`
		IsHazardous        bool   `json:"isHazardous"`
		StorageType        string `json:"storageType"`
		IsAnalgesic        bool   `json:"isAnalgesic"`
		ContainsLoperamide bool   `json:"containsLoperamide"`
	} `json:"productCharacteristics"`
	Gda struct {
		GdaRefs []struct {
			GdaDescription string   `json:"gdaDescription"`
			Headers        []string `json:"headers"`
			Footers        []string `json:"footers"`
			Values         []struct {
				Name    string   `json:"name"`
				Values  []string `json:"values"`
				Percent string   `json:"percent"`
				Rating  string   `json:"rating,omitempty"`
			} `json:"values"`
		} `json:"gdaRefs"`
	} `json:"gda"`
	CalcNutrition struct {
		Per100Header     string `json:"per100Header"`
		PerServingHeader string `json:"perServingHeader"`
		CalcNutrients    []struct {
			Name            string `json:"name"`
			ValuePer100     string `json:"valuePer100"`
			ValuePerServing string `json:"valuePerServing"`
			QualPerServing  string `json:"qualPerServing,omitempty"`
		} `json:"calcNutrients"`
	} `json:"calcNutrition"`
	Storage       []string `json:"storage"`
	MarketingText string   `json:"marketingText"`
	PkgDimensions []struct {
		No           int     `json:"no"`
		Height       float64 `json:"height"`
		Width        float64 `json:"width"`
		Depth        float64 `json:"depth"`
		DimensionUom string  `json:"dimensionUom"`
		Weight       float64 `json:"weight"`
		WeightUom    string  `json:"weightUom"`
		Volume       float64 `json:"volume"`
		VolumeUom    string  `json:"volumeUom"`
	} `json:"pkgDimensions"`
	ProductAttributes []struct {
		Category []struct {
			Lifestyle []struct {
				Lifestyle struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				} `json:"lifestyle"`
			} `json:"lifestyle"`
		} `json:"category"`
	} `json:"productAttributes"`
}

//Product contains Sapling's structured data about the product, as well as the unstructured Tesco-style Info
type Product struct {
	Ean string
	Name string
	Info Info
	CountryId int
	CategoryId int
	UnitId int
	Quantity int
	Virgin bool
	Footprint string
}

//findProduct searches the db by ean for a product and if it can't find one, calls the Tesco API
//Returns virgin value depending on whether the product has been scanned before
func findProduct(r *http.Request, ean string) (Product, error) {
	i := Info{}
	p := Product{}

	if ean == "" {
		return p, errors.New("400. Bad Request.")
	}

	p.Ean = ean
	p.Info = i
	p.Virgin = true

	row := config.Db.QueryRow("SELECT info, name, country_id, category_id, unit_id, quantity FROM products WHERE ean = $1", ean)

	// Whack the JSON into a string before unmarshalling it into the struct
	var col string
	err := row.Scan(&col, &p.Name, &p.CountryId, &p.CategoryId, &p.UnitId, &p.Quantity)

	if err == nil {
		//Product does already exist
		p.Virgin = false
		err = json.Unmarshal([]byte(col), &p.Info)
	} else if err.Error() == "sql: no rows in result set" {
		// Product doesn't exist. Let's ask Tesco about it...
		i = callTescoApi(ean)
		// TODO - some checking on whether the ean was valid in the first place
		p.Ean = ean
		p.Info = i
		p.Name = i.Description

		// Persist to the db, even if it's an empty product
		pbs, _ := json.Marshal(i)
		_, err = config.Db.Exec("INSERT INTO products (ean, info, name) VALUES ($1, $2, $3)",ean, pbs, p.Name )

		if err != nil {
			fmt.Println(err)
		}
	}

	// Is the user authenticated?
	u, _ := authUser(r)
	//Is the product unclaimed?
	row = config.Db.QueryRow("SELECT user_id FROM products_users WHERE product_ean = $1", ean)
	var pu string
	_ = row.Scan(&pu)
	if u.Email != "" && pu == "" {
		// Yes. Link this product to the user
		p.Virgin = true
		_, err = config.Db.Exec("INSERT INTO products_users (product_ean, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4)", ean, u.Id, time.Now(), time.Now())
	}

	return p, nil
}


//showProduct shows a single product
func showProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	ean := ps.ByName("ean")
	p, err := findProduct(r, ean)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a := auth(r)
	countries := getCountries()
	categories := getCategories()
	f := calcFootprint(p)
	data := struct{
		Auth bool
		P Product
		Countries []Country
		Categories []Category
		F string
		}{
			a,
			p,
			countries,
			categories,
			f,
	}

	err = config.TPL.ExecuteTemplate(w, "product.gohtml", data)
	if err != nil {
		fmt.Println(err)
		return
	}
}

//updateProduct stores user submitted data about a product
func updateProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	ean:= ps.ByName("ean")
	n := r.FormValue("name")
	co := r.FormValue("country-id")
	ca := r.FormValue("category-id")
	q := r.FormValue("quantity")
	u := r.FormValue("unit-id")

	_, err := config.Db.Exec(
		"UPDATE products SET name = $1, country_id = $2, category_id = $3, quantity = $4, unit_id = $5 WHERE ean = $6",
		n, co, ca, q, u, ean)

	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, "/product/" + ean, http.StatusSeeOther)
}

type ApiResponse struct {
	Products []Info
}

//callTescoApi calls the Tesco API
func callTescoApi(g string) Info {
	// Call the Tesco API
	client := &http.Client{}
	u := "https://dev.tescolabs.com/product/?gtin=" + g
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("Ocp-Apim-Subscription-Key", "fb708b6003e94a32861a6c0556601af4")
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	// Parse the response from JSON to Go
	bs := []byte(body)
	a := ApiResponse{}

	err := json.Unmarshal(bs, &a)
	if err != nil {
		fmt.Println(err)
	}

	// default Info in case none get returned
	i := Info{}

	// Check to see if there were any products returned
	if len(a.Products) > 0 {
		i = a.Products[0]
		// Product exists, call our footprint API
		if err != nil {
			fmt.Println(err)
		}
	}

	return i
}

func callFootprintApi(g string) string {
	// Call the Tesco API
	client := &http.Client{}
	u := "https://footprint-dot-sapling.appspot.com/using-tesco/gtin/" + g
	req, _ := http.NewRequest("GET", u, nil)
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(body)
	return string(body)

}