package products

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sapling/config"
)


type Product struct {
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

//One searches the db by gtin for a product and if it can't find one, calls the Tesco API
func One(r *http.Request) (Product, error) {
	p := Product{}

	g := r.FormValue("gtin")
	if g == "" {
		return p, errors.New("400. Bad Request.")
	}

	//Does the Gtin exist in the database?
	//For some reason, Tesco add a leading 0 to the gtin
	tg := "0" + g

	row := config.Db.QueryRow("SELECT info FROM products WHERE gtin = $1", tg)
	// Whack the JSON into a string
	var col string
	err := row.Scan(&col)

	if err == nil {
		err = json.Unmarshal([]byte(col), &p)
		if err != nil {
			fmt.Println(err)
		} else {
			return p, nil
		}
	}

	if err.Error() == "sql: no rows in result set" {
		// No. Let's ask Tesco about it...
		p = CallApi(g)
	}

	return p, nil
}

func Put(r *http.Request) (Product, error) {
	// get form values
	p := Product{}
	//p.Gtin = r.FormValue("gtin")

	// insert values
	//i, err := config.Products.InsertOne(context.TODO(), p)
	//if err != nil {
	//	return p, errors.New("500. Internal Server Error." + err.Error())
	//}
	//fmt.Println(i.InsertedID)
	return p, nil
}



