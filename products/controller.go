package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sapling/cloudvision"
	"sapling/config"
)

//Show a single product
func ShowProduct(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	v, err := One(r)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	err = config.TPL.ExecuteTemplate(w, "product.gohtml", v)
	if err != nil {
		fmt.Println(err)
		return
	}
}

//Show the form for creating a new product
func CreateProduct(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	config.TPL.ExecuteTemplate(w, "product-create.gohtml", "")
}

//Store a new product
func StoreProduct(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	mf, fh, err := r.FormFile("image")
	if err != nil {
		fmt.Println(err)
	}
	defer mf.Close()
	// Store the file here on the server
	wd, err := os.Getwd()
	path := filepath.Join(wd, "public", "pics", fh.Filename)
	nf, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
	}
	defer nf.Close()
	// copy
	mf.Seek(0,0)
	io.Copy(nf, mf)

	annotations, err := cloudvision.DetectText(w, path)
	config.TPL.ExecuteTemplate(w, "product-create.gohtml", annotations[0].Description)

}
