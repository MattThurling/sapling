package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"sapling/config"
)

func apiShowProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	enableCors(&w)
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	ean := ps.ByName("ean")

	// Is there a footprint value for this product?
	row := config.Db.QueryRow("SELECT co2 FROM footprints WHERE ean = $1", ean)
	var co2 int
	err := row.Scan(&co2)


	//Return the product info as JSON
	json.NewEncoder(w).Encode(co2)
	if err != nil {
		fmt.Println(err)
		return
	}
}
