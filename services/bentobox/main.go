package main

import (
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

	// setup the eth manager
	ethManager := eth.NewManager(cfg.EthHost, cfg.PollInterval, dbmap)

	// start network height (last block) loop
	go ethManager.LastBlockLoop()

	// start the eth query dispatcher loop
	//   reads the wanted from devp2p table
	//   and sends queries
	// TODO

	// start the ipfs loader loop
	//  reads the eth data table, find the elements
	//  not already added, to include them
	// TODO

	// TODO
	// we don't have proper metrics yet
	// shame on you herman

	// TODO
	// Make it a graceful shutdown with SIGINT
	select {}
}
