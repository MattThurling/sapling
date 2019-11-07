package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"sapling/config"
)


func home(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	config.TPL.ExecuteTemplate(w, "home.gohtml", "")
}

func scan(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	config.TPL.ExecuteTemplate(w, "scan.gohtml", "")
}

func terms(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	config.TPL.ExecuteTemplate(w, "terms.gohtml", "")
}

func privacy(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	config.TPL.ExecuteTemplate(w, "privacy.gohtml", "")
}

func iframe(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	config.TPL.ExecuteTemplate(w, "iframe.gohtml", "")
}

// Caution: this opens the API to anyone
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

//Go application entrypoint
func main() {

	router := httprouter.New()

	router.GET("/", home)
	router.GET("/scan", scan)
	router.GET("/iframe", iframe)
	router.GET("/product/:ean", showProduct)
	router.POST("/product/:ean", updateProduct)
	router.GET("/register", register)
	router.POST("/register", createUser)
	router.GET("/login", login)
	router.POST("/login", postLogin)
	router.GET("/logout", logout)
	router.GET("/ocr", ocrForm)
	router.POST("/ocr", ocrUpload)
	router.GET("/dashboard", dashboard)
	router.GET("/leaderboard", leaderboard)
	router.GET("/terms", terms)
	router.GET("/privacy", privacy)

	router.POST("/api/basket", createBasket)
	router.GET("/api/basket/:id", showBasket)
	router.GET("/api/product/:ean", apiShowProduct)


	router.ServeFiles("/static/*filepath", http.Dir("static"))

	//log.Fatal(fmt.Println(http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", router))) // production, https

	log.Fatal(fmt.Println(http.ListenAndServe(":8080", router))) // local

}
