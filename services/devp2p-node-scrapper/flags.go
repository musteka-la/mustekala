package main

import (
	"flag"
	"os"
	"path/filepath"

	logging "github.com/ipfs/go-log"
	whylogging "github.com/whyrusleeping/go-logging"
)

// Config has all the options you defined at the command line.
type Config struct {
	Debug            bool
	BootnodesPath    string
	NodeDatabasePath string
	DevP2PLibDebug   bool
}

// ParseFlags gets those command line options
// and set them into a nice Config struct.
func ParseFlags() *Config {
	cfg := &Config{}

	flag.BoolVar(&cfg.Debug, "debug", false, "set this variable to have logging level DEBUG")
	flag.StringVar(&cfg.BootnodesPath, "devp2p-bootnodes", "", "Location of devp2p bootnodes file")
	flag.StringVar(&cfg.NodeDatabasePath, "devp2p-nodes-database", "", "Location of the devp2p node database")
	flag.BoolVar(&cfg.DevP2PLibDebug, "devp2p-lib-debug", false, "set this variable if you really like logs (p2p lib logs)")
	flag.Parse()

	if cfg.Debug {
		logging.SetDebugLogging()
	} else {
		logging.SetAllLoggers(whylogging.INFO)
	}

	if cfg.NodeDatabasePath == "" {
		homeDir := os.Getenv("HOME")
		cfg.NodeDatabasePath = filepath.Join(homeDir, ".mustekala", "devp2p", "nodes")
	}

	return cfg
}
