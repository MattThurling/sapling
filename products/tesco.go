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

	p := Product{}

	// Check to see if there were any products returned
	if len(a.Products) > 0 {
		p = a.Products[0]
		// Persist to the db
		pbs, err := json.Marshal(p)
		_, err = config.Db.Exec("INSERT INTO products (gtin, info) VALUES ($1, $2)","0" + g, pbs )
		if err != nil {
			fmt.Println(err)
		}
	}

	return p
}