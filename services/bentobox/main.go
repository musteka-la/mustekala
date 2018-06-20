package main

import (
	"fmt"

	"github.com/metamask/mustekala/services/bentobox/db"
)

func main() {
	// get the flags
	cfg := ParseFlags()

	// load PostgreSQL Database
	dbOpts := db.Options{
		User:     *cfg.DbUser,
		Password: *cfg.DbPassword,
		DBName:   *cfg.DbName,
	}
	dbmap := db.InitDb(dbOpts)
	defer dbmap.Db.Close()

	// PLACEHOLDER
	fmt.Printf("OK")
	// PLACEHOLDER
}
