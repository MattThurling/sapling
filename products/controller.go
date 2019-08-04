package products

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
func Show(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	p, err := One(r)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	if p.Description == "" {
		fmt.Println("Empty product")
		http.Redirect(w, r, "/product/create", 303)
	}

	config.TPL.ExecuteTemplate(w, "product.gohtml", p)
}

//Show the form for creating a new product
func Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	config.TPL.ExecuteTemplate(w, "product-create.gohtml", "")
}

//Store a new product
func Store(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	config.TPL.ExecuteTemplate(w, "product-create.gohtml", annotations)

}
