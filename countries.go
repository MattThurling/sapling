package main

import (
	"fmt"
	"sapling/config"
)

type Country struct {
	Id int
	Name string
}

func getCountries() []Country {

	rows, err := config.Db.Query("SELECT id, country FROM countries")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	cs := make([]Country, 0)
	for rows.Next() {
		c := Country{}
		err := rows.Scan(&c.Id, &c.Name)
		if err != nil {
			fmt.Println(err)
		}
		cs = append(cs, c)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err)
	}

	return cs
}
