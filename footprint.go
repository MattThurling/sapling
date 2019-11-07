package main

import (
	"fmt"
	"sapling/config"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func calcFootprint(p Product) string {
	//Get the sub region code
	row := config.Db.QueryRow("SELECT region_code FROM countries WHERE id = $1", p.CountryId)
	var s int
	_ = row.Scan(&s)

	//var columnName string
	switch s {
	case 150:
		//columnName = "ghg_europe"
	case 40:
	default:
		// Do Something
		break;
	}

	row = config.Db.QueryRow("SELECT ghg_global FROM categories WHERE id = $1", p.CategoryId)
	var ghg int
	err := row.Scan(&ghg)

	fmt.Println(err)

	pghg := ghg * p.Quantity / 1000

	f := message.NewPrinter(language.English)
	footprint := f.Sprintf("%d\n", pghg)
	//fallback if there is no value for the specified region

	return footprint
}

