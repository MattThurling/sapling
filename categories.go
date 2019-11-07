package main

import (
	"fmt"
	"sapling/config"
)

type Category struct {
	Id int
	Name string
}

func getCategories() []Category {

	rows, err := config.Db.Query("SELECT id, name FROM categories")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	cs := make([]Category, 0)
	for rows.Next() {
		c := Category{}
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
