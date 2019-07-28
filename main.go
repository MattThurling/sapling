package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var templates *template.Template

func home(w http.ResponseWriter, r *http.Request){
	templates.ExecuteTemplate(w, "home.gohtml", "")
}

func scan(w http.ResponseWriter, r *http.Request){
	templates.ExecuteTemplate(w, "scan.gohtml", "")
}

func forest(w http.ResponseWriter, r *http.Request){
	templates.ExecuteTemplate(w, "forest.gohtml", "")
}


//Go application entrypoint
func main() {

	//We tell Go exactly where we can find our html file. We ask Go to parse the html file (Notice
	// the relative path). We wrap it in a call to template.Must() which handles any errors and halts if there are fatal errors

	templates = template.Must(template.ParseGlob("templates/*.gohtml"))

	//Our HTML comes with CSS that go needs to provide when we run the app. Here we tell go to create
	// a handle that looks in the static directory, go then uses the "/static/" as a url that our
	//html can refer to when looking for our css and other files.

	http.Handle("/static/", //final url can be anything
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static")))) //Go looks in the relative "static" directory first using http.FileServer(), then matches it to a
	//url of our choice as shown in http.Handle("/static/"). This url is what we need when referencing our css files
	//once the server begins. Our html code would therefore be <link rel="stylesheet"  href="/static/stylesheet/...">
	//It is important to note the url in http.Handle can be whatever we like, so long as we are consistent.


	http.HandleFunc("/", home)
	http.HandleFunc("/scan", scan)
	http.HandleFunc("/forest", forest)

	//Start the web server, set the port to listen to 8080. Without a path it assumes localhost
	//Print any errors from starting the webserver using fmt
	fmt.Println("Listening")
	fmt.Println(http.ListenAndServe(":8080", nil))
}
