package main

import (
	"fmt"

	"github.com/metamask/mustekala/services/bentobox/db"
	"github.com/metamask/mustekala/services/bentobox/eth"
)

func main() {
	// get the flags
	cfg := ParseFlags()

	// load PostgreSQL Database
	dbOpts := db.Options{
		User:     cfg.DbUser,
		Password: cfg.DbPassword,
		DBName:   cfg.DbName,
	}
	dbmap := db.InitDb(dbOpts)
	defer dbmap.Db.Close()

	// PLACEHOLDER
	// Just do a single last block query
	eth.GetNetworkHeight(cfg.EthHost)

	fmt.Printf("OK!\n")
	// PLACEHOLDER
}
