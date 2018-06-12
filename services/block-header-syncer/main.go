package main

import (
	"fmt"
	"os"

	"github.com/metamask/mustekala/services/lib/db"
	"github.com/metamask/mustekala/services/lib/devp2p"
)

func main() {
	// get the flags
	cfg := ParseFlags()

	// set up database
	dbPool := db.NewPool(cfg.DatabaseConn)
	if _, err := dbPool.Get().Do("PING"); err != nil {
		fmt.Printf("Database Error: %v\n", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	// setup the devp2p server
	devp2pConfig := &devp2p.Config{
		BootnodesPath:           cfg.BootnodesPath,
		NodeDatabasePath:        cfg.NodeDatabasePath,
		LibP2PDebug:             cfg.DevP2PLibDebug,
		DbPool:                  dbPool,
		IsPeerScrapperActive:    false,
		IsSyncBlockHeaderActive: true,
	}

	devp2pServer := devp2p.NewManager(devp2pConfig)

	// setup the syncer
	// syncer := NewSyncer(devp2pServer, dbPool)

	// start the devp2p server
	go devp2pServer.Start()

	// start the block header syncing loop
	// go syncer.Start()

	// TODO
	// a nice SIGINT manager
	// * closing the devp2p server
	// * nicely turning off the syncer
	// * and closing the redis pool
	select {}
}
