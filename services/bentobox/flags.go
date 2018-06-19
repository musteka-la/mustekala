package main

// Config has all the options you defined at the command line.
type Config struct {
	Debug bool
}

// ParseFlags gets those command line options
// and set them into a nice Config struct.
func ParseFlags() *Config {
	cfg := &Config{}

	// optionsDBUser := flag.String("dbuser", "postgres", "database username")
	// optionsDBPassword := flag.String("dbpassword", "mysecretpassword", "database password")
	// optionsDBName := flag.String("dbname", "bentobox", "database name")

	// optionsHost := flag.String("host", "http://127.0.0.1:8545", "URL of the ethereum node RPC")
	// optionsMaxProcessingQueries := flag.Int("maxprocessingqueries", 100,
	//	"Maximum of concurrent queries to the RPC Node")
	// optionsLoopTimeMs := flag.Int("looptimems", 5000, "Iteration interval for block querying")

	return cfg
}
