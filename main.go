package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"sapling/config"
	"sapling/products"
)



type Product struct {
	id uint
	code string
	carbon int
}

func home(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	config.TPL.ExecuteTemplate(w, "home.gohtml", "")
}

func scan(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	config.TPL.ExecuteTemplate(w, "scan.gohtml", "")
}

func iframe(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	config.TPL.ExecuteTemplate(w, "iframe.gohtml", "")
}

func register(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	config.TPL.ExecuteTemplate(w, "register.gohtml", "")
}

func login(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	config.TPL.ExecuteTemplate(w, "login.gohtml", "")
}


//Go application entrypoint
func main() {

	router := httprouter.New()

	router.GET("/", home)
	router.GET("/scan", scan)
	router.GET("/iframe", iframe)
	router.GET("/product", products.Show)
	router.GET("/register", register)
	router.GET("/login", login)
	router.GET("/product/create", products.Create)
	router.POST("/product/create", products.Store)

	router.ServeFiles("/static/*filepath", http.Dir("static"))

	log.Fatal(fmt.Println(http.ListenAndServe(":8080", router)))

}
