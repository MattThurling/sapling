package main

import (
	"database/sql"
	"fmt"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var templates *template.Template
var db *sql.DB

type Product struct {
	id uint
	code string
	carbon int
}

func home(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	templates.ExecuteTemplate(w, "home.gohtml", "")
}

func scan(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	templates.ExecuteTemplate(w, "scan.gohtml", "")
}

func forest(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	templates.ExecuteTemplate(w, "forest.gohtml", "")
}

func showProduct(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	// Look up product
	u := "https://api.barcodelookup.com/v2/products?barcode=" + ps.ByName("code") + "&formatted=y&key=uq257lzb7kcgvpu5yfijlgfyz2dsf4"
	fmt.Println(u)
	resp, _ := http.Get(u)

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	bs := string(body)

	//p := Product{}
	//row := db.QueryRow("SELECT * FROM products WHERE id = $1", c)
	//err := row.Scan(&p.id, &p.code, &p.carbon)
	//if err != nil {
	//	panic(err)
	//}
	//
	//switch {
	//case err == sql.ErrNoRows:
	//	http.NotFound(w, r)
	//	return
	//case err != nil:
	//	http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	//	return
	//}
	templates.ExecuteTemplate(w, "product.gohtml", bs)
}

//Set up connection to the database
func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:21satoshi@localhost:5433/sapling?sslmode=disable")
	if err != nil {
		panic(err)
	}


	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to database")
}


//Go application entrypoint
func main() {

	//We tell Go exactly where we can find our html file. We ask Go to parse the html file (Notice
	// the relative path). We wrap it in a call to template.Must() which handles any errors and halts if there are fatal errors

	templates = template.Must(template.ParseGlob("templates/*.gohtml"))


	router := httprouter.New()

	router.GET("/", home)
	router.GET("/scan", scan)
	router.GET("/forest", forest)
	router.GET("/product/:code", showProduct)

	router.ServeFiles("/static/*filepath", http.Dir("static"))

	log.Fatal(fmt.Println(http.ListenAndServe(":8080", router)))

	db.Close()
}
