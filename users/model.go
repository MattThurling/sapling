package users

import (
	"fmt"
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
	if err != nil {
		fmt.Println(err)
	}
	return u, nil
}

