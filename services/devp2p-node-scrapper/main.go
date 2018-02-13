package main

import (
	"github.com/metamask/mustekala/devp2p"
)

func main() {
	// get the flags
	cfg := ParseFlags()

	// TODO
	// set up database

	// setup the devp2p server
	devp2pConfig := &devp2p.Config{
		BootnodesPath:        cfg.BootnodesPath,
		NodeDatabasePath:     cfg.NodeDatabasePath,
		LibP2PDebug:          cfg.DevP2PLibDebug,
		IsPeerScrapperActive: true, // the main point of this service
		// TODO
		// database conn
	}

	devp2pServer := devp2p.NewManager(devp2pConfig)

	// start the devp2p server
	go devp2pServer.Start()

	// TODO
	// a nice SIGINT manager
	select {}
}
