package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var Db *sql.DB

//Set up connection to the database
func init() {
	var err error
	Db, err = sql.Open(
		"postgres",
		//"postgres://postgres:21satoshi@localhost:5433/sapling?sslmode=disable"
		"postgres://postgres:HJ7hyalql2mkMFp7@/sapling?host=/cloudsql/sapling:europe-west1:sapling",
		)
	if err != nil {
		panic(err)
	}


	err = Db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to database")

}
