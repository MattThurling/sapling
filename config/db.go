package config

import (
	"fmt"
	"github.com/globalsign/mgo"
)

// database
var DB *mgo.Database

// collections
var Products *mgo.Collection

func init() {
	// get a mongo session
	s, err := mgo.Dial("mongodb://sapling:password@localhost/sapling")
	if err != nil {
		panic(err)
	}

	if err = s.Ping(); err != nil {
		panic(err)
	}

	DB = s.DB("sapling")
	Products = DB.C("products")

	fmt.Println("You connected to your mongo database.")
}
