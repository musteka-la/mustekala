package main

import "fmt"

func main() {
	// get the flags
	cfg := ParseFlags()

	// Load SQL Database
	/*
		dbOpts := database.Options{
			User:     *optionsDBUser,
			Password: *optionsDBPassword,
			DBName:   *optionsDBName,
		}
		dbmap := database.InitDb(dbOpts)
		defer dbmap.Db.Close()
	*/

	// PLACEHOLDER
	_ = cfg
	fmt.Printf("OK")
	// PLACEHOLDER
}
