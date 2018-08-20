package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/metamask/mustekala/services/eth-db-heatmap/lib"
)

/*

## ETHEREUM DATABASE HEATMAP

See the README.md

### EXAMPLE USAGE

./build/bin/eth-db-heatmap \
	--block-number 6035070 \
	--geth-db-filepath /Users/hj/Documents/_ethHeatmap/geth/chaindata

*/

func main() {
	var (
		blockNumber uint64
		dbFilePath  string
	)

	// command line options
	flag.Uint64Var(&blockNumber, "block-number", 6035070, "Canonical number of the block state to import")
	flag.StringVar(&dbFilePath, "geth-db-filepath", "", "Path to the Go-Ethereum Database")
	flag.Parse()

	// we just hardcode this one
	fileDir := fmt.Sprintf("/tmp/slicedata-%d", blockNumber)

	// welcome message
	log.Printf("Ethereum Database Heatmap Generator")

	// init the slice processor
	sliceProcessor := lib.NewSliceProcessor(dbFilePath, blockNumber, fileDir)

	// launch main loop:
	// we take the slices, and store the results in the /tmp dir
	// for further recovery of the data
	sliceProcessor.SliceTheTrie()

	// we compute the heatmap here
	// output goes to the /tmp directory as well.
	// you should be able to parse it with little work into a nice graph
	sliceProcessor.GetHeatMap("txt")

	// done
	log.Println("we are done here [END]")
}
