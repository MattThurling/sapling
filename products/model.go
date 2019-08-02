package products

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
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

//func All() ([]Product, error) {
//	ps := []Product{}
//	err := config.Products.Find(config.CTX, bson.D{})
//
//	//err := config.Products.Find(config.CTX, bson.D{}).All(&ps)
//
//	if err != nil {
//		return nil, err
//	}
//	return ps, nil
//}

func One(r *http.Request) (Product, error) {
	p := Product{}

	gtin := r.FormValue("gtin")
	if gtin == "" {
		return p, errors.New("400. Bad Request.")
	}

	//Does the Gtin exist in the database?
	err := config.Products.FindOne(config.CTX, bson.M{"gtin": gtin}).Decode(&p)
	fmt.Println(err)
	if err.Error() == "mongo: no documents in result" {
		// No. Let's ask Tesco about it...
		p = CallApi(gtin)
	}

	return p, nil
}

//func Put(r *http.Request) (Product, error) {
//	// get form values
//	p := Product{}
//	p.Gtin = r.FormValue("gtin")
//
//	// insert values
//	err := config.Products.InsertOne(config.CTX, p)
//	if err != nil {
//		return p, errors.New("500. Internal Server Error." + err.Error())
//	}
//	return p, nil
//}
//
//func Update(r *http.Request) (Product, error) {
//	// get form values
//	p := Product{}
//	p.Gtin = r.FormValue("gtin")
//
//	// update values
//	err := config.Products.Update(bson.M{"gtin": p.Gtin}, &p)
//	if err != nil {
//		return p, err
//	}
//	return p, nil
//}
//
//func Delete(r *http.Request) error {
//	gtin := r.FormValue("gtin")
//	if gtin == "" {
//		return errors.New("400. Bad Request.")
//	}
//
//	err := config.Products.Remove(bson.M{"gtin": gtin})
//	if err != nil {
//		return errors.New("500. Internal Server Error")
//	}
//	return nil
//}


