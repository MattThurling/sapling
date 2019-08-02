package products

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sapling/config"
)

type ApiResponse struct {
	Products []Product
	}

func CallApi(g string) Product {

	client := &http.Client{}
	u := "https://dev.tescolabs.com/product/?gtin=" + g
	fmt.Println(u)
	req, _ := http.NewRequest("GET", u, nil)

	req.Header.Set("Ocp-Apim-Subscription-Key", "fb708b6003e94a32861a6c0556601af4")
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	bs := []byte(body)
	a := ApiResponse{}

	err := json.Unmarshal(bs, &a)

	p := a.Products[0]

	// Persist to the db
	err = config.Products.Insert(p)

	if err != nil {
		fmt.Println(err)
	}

	return p
}