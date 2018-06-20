package main

import "flag"

// Config has all the options you defined at the command line.
type Config struct {
	DbUser       string
	DbPassword   string
	DbName       string
	EthHost      string
	IpfsHost     string
	PollInterval int
}

// ParseFlags gets those command line options
// and set them into a nice Config struct.
func ParseFlags() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.DbUser, "dbuser", "postgres", "database username")
	flag.StringVar(&cfg.DbPassword, "dbpassword", "mysecretpassword", "database password")
	flag.StringVar(&cfg.DbName, "dbname", "bentobox", "database name")

	flag.StringVar(&cfg.EthHost, "eth-host", "http://127.0.0.1:8545", "URL of the ethereum node RPC")
	flag.StringVar(&cfg.IpfsHost, "ipfs-host", "http://127.0.0.1:5001", "URL of the IPFS HTTP API")

	flag.IntVar(&cfg.PollInterval, "last-block-polling-interval", 1, "Iteration interval for last block querying")

	return cfg
}
