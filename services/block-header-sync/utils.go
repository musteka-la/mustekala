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
	// Bridge / Overall Options
	BridgeChannelCapacity int
	BridgeDebug           bool

	// DevP2P Options
	DevP2PBootnodesPath     string
	DevP2PNodeDatabasePath  string
	DevP2PNetworkStatusPath string
	DevP2PLibDebug          bool
}

// ParseFlags gets those command line options
// and set them into a nice Config struct.
func ParseFlags() *Config {
	cfg := &Config{}

	// Bridge / Overall Flags
	flag.IntVar(&cfg.BridgeChannelCapacity, "bridge-channel-cap", 100, "channel capacity of messages")
	flag.BoolVar(&cfg.BridgeDebug, "bridge-debug", false, "set this variable to have logging level DEBUG")

	// DevP2P Flags
	flag.StringVar(&cfg.DevP2PBootnodesPath, "devp2p-bootnodes", "", "Location of devp2p bootnodes file")
	flag.StringVar(&cfg.DevP2PNodeDatabasePath, "devp2p-nodes-database", "", "Location of the devp2p node database")
	flag.StringVar(&cfg.DevP2PNetworkStatusPath, "devp2p-network-status", "", "Location of the devp2p network status file")
	flag.BoolVar(&cfg.DevP2PLibDebug, "devp2p-lib-debug", false, "set this variable if you really like logs")
	flag.Parse()

	if cfg.BridgeDebug {
		logging.SetDebugLogging()
	} else {
		logging.SetAllLoggers(whylogging.INFO)
	}

	if cfg.DevP2PNodeDatabasePath == "" {
		homeDir := os.Getenv("HOME")
		cfg.DevP2PNodeDatabasePath = filepath.Join(homeDir, ".mustekala", "devp2p", "nodes")
	}

	return cfg
}
