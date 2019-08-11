package users

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"sapling/config"
)

type Session struct {
	Id string
	User_id string
}

func CreateSession(r *http.Request, u_id string) (Session, error) {

	id, _ := uuid.NewV4()

	s := Session{
		id.String(),
		u_id,
	}

	_, err := config.Db.Exec("INSERT INTO sessions (id, user_id) VALUES ($1, $2)", s.Id, s.User_id)
	if err != nil {
		fmt.Println(err)
	}
	return s, nil
}


//Auth checks the cookie session value against the sessions table in the database
func Auth(req *http.Request) bool {
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
func AuthUser(req *http.Request) (User, error){

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

