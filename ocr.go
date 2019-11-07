package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sapling/config"
)

//createProduct shows the form for uploading a picture
func ocrForm(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	config.TPL.ExecuteTemplate(w, "ocr.gohtml", "")
}

//storeProduct handles the upload of an image and sends it for OCR
func ocrUpload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	annotations, err := DetectText(w, path)
	config.TPL.ExecuteTemplate(w, "ocr.gohtml", annotations[0].Description)

}
