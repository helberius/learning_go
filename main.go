package main

import (
	"./owid"

	"./usgs"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {

	// method to retrieve the quakes occurred in the last hour from the usgs databas using a url
	usgs.GetLastHourQuakes()

	// testing how to query in a postgres DB
	owid.IsDataOlderThanDays(5)

}
