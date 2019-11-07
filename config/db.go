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
	//datastoreName := os.Getenv("POSTGRES_CONNECTION")
	Db, err = sql.Open(
		"postgres",
		//"postgres://postgres:21satoshi@localhost:5433/sapling?sslmode=disable", // running Postgres locally
		"postgres://postgres:HJ7hyalql2mkMFp7@localhost:5434/sapling?sslmode=disable", // running Postgres via proxy
		//datastoreName,
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
