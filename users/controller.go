package users

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"sapling/config"
	//"sapling/products"
)

//Register shows the form for registering a new user
func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	config.TPL.ExecuteTemplate(w, "register.gohtml", "")
}

//Login shows the form for logging in an existing user and handles the form submission
func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params){

	if Auth(r) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	config.TPL.ExecuteTemplate(w, "login.gohtml", "")
}


//Store adds a new user to the database and creates a new session for them
//and writes the session id to a cookie
func Store(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	//Create the user
	u, err := Put(r)

	if err != nil {
		fmt.Fprintln(w, err)
	}
	//Create session for the user
	s, err := CreateSession(r, u.Id)

	if err != nil {
		fmt.Fprintln(w, err)
	}

	// Store session id in a cookie
	c := &http.Cookie{
		Name:  "session",
		Value: s.Id,
	}
	http.SetCookie(w, c)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

}

//PostLogin processes the login form submission
func PostLogin (w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	if r.Method == http.MethodPost {
		e := r.FormValue("email")
		p := r.FormValue("password")

		u := User{}

		if e == "" || p == "" {
			http.Error(w,"400. Bad Request.", http.StatusBadRequest)
			return
		}

		//Does the user email exist in the database?
		row := config.Db.QueryRow("SELECT * FROM users WHERE email = $1", e)
		// Whack the JSON into a string
		err := row.Scan(&u.Id, &u.Name, &u.Email, &u.Password)

		// does the entered password match the stored password?
		err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))

		if err != nil {
			http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			return
		}

		// create session

		s, err := CreateSession(r, u.Id)

		c := &http.Cookie{
			Name:  "session",
			Value: s.Id,
		}
		http.SetCookie(w, c)

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
}

//Dashboard shows the user dashboard
//func Dashboard(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
//	var ps []products.Product
//
//	rows, err := config.Db.Query("SELECT gtin FROM products")
//
//	if err != nil {
//		http.Error(w, http.StatusText(500), 500)
//		return
//	}
//	defer rows.Close()
//
//	for rows.Next() {
//		p := products.Product{}
//		err := rows.Scan(&p.Gtin)
//		if err != nil {
//			http.Error(w, http.StatusText(500), 500)
//			return
//		}
//		ps = append(ps, p)
//	}
//	if err = rows.Err(); err != nil {
//		http.Error(w, http.StatusText(500), 500)
//		return
//	}
//
//	config.TPL.ExecuteTemplate(w, "dashboard.gohtml", ps)
//}