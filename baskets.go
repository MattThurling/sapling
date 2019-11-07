package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"sapling/config"
)

type Basket struct {
	Id string
	PurchaserId string
	VendorId string
	Items []Item
}

type Item struct {
	Ean string
	Quantity uint
}

//showBasket shows details of a unique transaction and its items
func showBasket(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	enableCors(&w)
	id := ps.ByName("id")
	b := Basket{
		id,
		"abc",
		"def",
		[]Item{},
	}

	//Get all the products in the basket
	items := []Item{}
	rows, _ := config.Db.Query("SELECT ean, quantity FROM items WHERE basket_id = $1", id )

	defer rows.Close()

	for rows.Next() {
		i := Item{}
		_ = rows.Scan(&i.Ean, &i.Quantity)

		items = append(items, i)
	}
	b.Items = items
	//Return the basket as JSON
	json.NewEncoder(w).Encode(b)
}

//createBasket creates a unique record of a transaction and its items
func createBasket(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//TODO: restrict this route to authenticated user
	// Is the user authenticated?
	u, _ := authUser(r)
	//Create a unique basket for this transaction
	id, _ := uuid.NewV4()
	b := Basket{
		id.String(),
		r.FormValue("purchaser_id"),
		r.FormValue("vendor_id"),
		[]Item{},
	}
	_, err := config.Db.Exec("INSERT INTO baskets (id, purchaser_id, vendor_id) VALUES ($1, $2, $3)", b.Id, b.Id, u.Id)

	if err != nil {
		fmt.Println(err)
	}

	//Add the items
	bs := []byte(r.FormValue("details"))
	items := []Item{}

	err = json.Unmarshal(bs, &items)

	fmt.Println(items)

	for _, item := range items {
		_, err = config.Db.Exec("INSERT INTO items (ean, quantity, basket_id) VALUES ($1, $2, $3)", item.Ean, item.Quantity, b.Id)
	}

	if err != nil {
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(items)

}
