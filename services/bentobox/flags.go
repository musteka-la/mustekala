package main

import "flag"

// Config has all the options you defined at the command line.
type Config struct {
	DbUser       *string
	DbPassword   *string
	DbName       *string
	EthHost      *string
	IpfsHost     *string
	PollInterval *int
}

// ParseFlags gets those command line options
// and set them into a nice Config struct.
func ParseFlags() *Config {
	cfg := &Config{}

	cfg.DbUser = flag.String("dbuser", "postgres", "database username")
	cfg.DbPassword = flag.String("dbpassword", "mysecretpassword", "database password")
	cfg.DbName = flag.String("dbname", "bentobox", "database name")

	cfg.EthHost = flag.String("eth-host", "http://127.0.0.1:8545", "URL of the ethereum node RPC")
	cfg.IpfsHost = flag.String("ipfs-host", "http://127.0.0.1:5001", "URL of the IPFS HTTP API")

	cfg.PollInterval = flag.Int("last-block-polling-interval", 1, "Iteration interval for last block querying")

	return cfg
}
