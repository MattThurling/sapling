package sessions

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

func Put(r *http.Request, u_id string) (Session, error) {

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


//HasActiveSession checks the cookie session value against the sessions table in the database
func HasActiveSession(req *http.Request) bool {
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
