package main

import (
	"flag"
	"fmt"
	"os"
)

// Config has all the options you defined at the command line.
type Config struct {
	DatabaseConn     string
	NumberOfAccounts int
}

// ParseFlags gets those command line options
// and set them into a nice Config struct.
func ParseFlags() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.DatabaseConn, "redis-conn", ":6379", "redis DB connection string")
	flag.IntVar(&cfg.NumberOfAccounts, "accounts", 0, "number of accounts to create in the secure state trie")
	flag.Parse()

	if cfg.NumberOfAccounts < 1 {
		fmt.Println("You want to set up the --accounts option to a number greater than zero")
		os.Exit(1)
	}

	return cfg
}
