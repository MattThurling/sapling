package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"sapling/config"
)

type User struct {
	Id string
	Name string
	Email string
	Password string
}

//Put creates a new user record in the database
func Put(r *http.Request) (User, error) {

	id, _ := uuid.NewV4()
	pw := r.FormValue("password")
	h, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.MinCost)

	u := User{
		id.String(),
		r.FormValue("name"),
		r.FormValue("email"),
		string(h),
	}
	_, err := config.Db.Exec("INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4)", u.Id, u.Name, u.Email, u.Password)

	return u, err
}

//Register shows the form for registering a new user
func register(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	if auth(r) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	config.TPL.ExecuteTemplate(w, "register.gohtml", "")
}

//Login shows the form for logging in an existing user and handles the form submission
func login(w http.ResponseWriter, r *http.Request, _ httprouter.Params){

	if auth(r) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
	config.TPL.ExecuteTemplate(w, "login.gohtml", "")
}


//Store adds a new user to the database and creates a new session for them
//and writes the session id to a cookie
func createUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	s, err := createSession(r, u.Id)

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
func postLogin (w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

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
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		// create session

		s, err := createSession(r, u.Id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		c := &http.Cookie{
			Name:  "session",
			Value: s.Id,
		}
		http.SetCookie(w, c)

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}
}

//logout removes the current session and cookie
func logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !auth(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	c, _ := r.Cookie("session")
	// delete the session
	_, _ = config.Db.Exec("DELETE FROM sessions WHERE id = $1", c.Value)
	// remove the cookie
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

//dashboard shows the user dashboard
func dashboard(w http.ResponseWriter, r *http.Request, _ httprouter.Params){

	data := struct{
		Ps []Product
		U User
		}{}

	data.U, _ = authUser(r)

	var ps []Product

	rows, err := config.Db.Query("SELECT name, ean FROM products JOIN products_users ON ean = product_ean WHERE user_id = $1", data.U.Id)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	for rows.Next() {

		//Create an empty Product
		p := Product{}
		err := rows.Scan(&p.Name, &p.Ean)

		if err != nil {
			fmt.Println(err)
		}

		ps = append(ps, p)
	}

	data.Ps = ps

	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	config.TPL.ExecuteTemplate(w, "dashboard.gohtml", data)
}

//leaderboard shows the users ranked by the number of scans
func leaderboard(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type leader struct{
		Name string
		Score int
	}
	leaders := []leader{}

	rows, _ := config.Db.Query("SELECT name, count(*) FROM users JOIN products_users ON id = user_id GROUP BY id ORDER BY count(*) DESC")

	defer rows.Close()

	for rows.Next() {
		l := leader{}
		_ = rows.Scan(&l.Name, &l.Score)

		leaders = append(leaders, l)
	}

	config.TPL.ExecuteTemplate(w, "leaderboard.gohtml", leaders)
}

type Session struct {
	Id string
	UserId string
}

func createSession(r *http.Request, uId string) (Session, error) {

	id, _ := uuid.NewV4()

	s := Session{
		id.String(),
		uId,
	}

	_, err := config.Db.Exec("INSERT INTO sessions (id, user_id) VALUES ($1, $2)", s.Id, s.UserId)

	return s, err
}


//Auth checks the cookie session value against the sessions table in the database
func auth(req *http.Request) bool {
	c, err := req.Cookie("session")
	if err != nil {
		return false
	}

	//Check a session exists
	row := config.Db.QueryRow("SELECT user_id FROM sessions WHERE id = $1", c.Value)
	var col string
	err = row.Scan(&col)

	if err != nil {
		return false
	}

	//Check a user exists
	row = config.Db.QueryRow("SELECT email FROM users WHERE id = $1", col)
	err = row.Scan(&col)

	if err != nil {
		return false
	}

	fmt.Println(col)

	return true
}

//AuthUser returns the authenticated user
func authUser(req *http.Request) (User, error) {

	u := User{}

	c, err := req.Cookie("session")
	if err != nil {
		return u, err
	}

	//Check a session exists
	row := config.Db.QueryRow("SELECT user_id FROM sessions WHERE id = $1", c.Value)
	var col string
	err = row.Scan(&col)

	if err != nil {
		return u, err
	}

	//Fetch the user
	row = config.Db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", col)
	err = row.Scan(&u.Id, &u.Name, &u.Email)

	if err != nil {
		return u, err
	}

	fmt.Println(u)

	return u, nil
}

