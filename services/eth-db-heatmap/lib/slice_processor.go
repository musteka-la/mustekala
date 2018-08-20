package lib

import (
	"bytes"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

// SliceProcessor wraps the goque stack, enabling the adding of specific
// methods for dealing with the state trie.
type SliceProcessor struct {
	stateRoot common.Hash
	db        *ethdb.LDBDatabase
	trie      *trie.SecureTrie
	trieDB    *trie.Database
	fileDir   string
}

// SliceProcessorConfig contains tweakable properties
type SliceProcessorConfig struct {
	DatabaseCache   int
	TrieCache       int
	MaxTrieCacheGen uint16
}

// NewSliceProcessor initializes the traversal stack, and finds the canonical
// block header, returning the TrieStack wrapper for further instructions
func NewSliceProcessor(dbPath string, blockNumber uint64, filedir string) *SliceProcessor {
	var err error
	sliceProcessor := &SliceProcessor{}

	log.Println("Init parameters and opening the DB")

	// init config parameterss
	config := &SliceProcessorConfig{
		DatabaseCache:   768, // taken from geth defaults at eth/config
		TrieCache:       256, // taken from geth defaults at eth/config
		MaxTrieCacheGen: 120, // taken from core/state/database.go
	}

	// setup the storage directory
	sliceProcessor.fileDir = filedir

	// open the database
	db, err := ethdb.NewLDBDatabase(dbPath, config.DatabaseCache, config.TrieCache)
	if err != nil {
		panic("can't open ethdb!")
	}
	sliceProcessor.db = db

	// we need to get the state root from the database
	// will query for its canonical block number
	blockHash := sliceProcessor.getCanonicalHash(blockNumber)
	headerRLP := sliceProcessor.getHeaderRLP(blockHash, blockNumber)
	header := new(types.Header)
	if err := rlp.Decode(bytes.NewReader(headerRLP), header); err != nil {
		panic(err)
	}

	log.Printf("For the given canonical block number: %d, found this root: %x",
		blockNumber, header.Root[:6])
	sliceProcessor.stateRoot = header.Root

	// init the trie
	sliceProcessor.trieDB = trie.NewDatabase(db)
	tr, err := trie.NewSecure(header.Root, sliceProcessor.trieDB, config.MaxTrieCacheGen)
	if err != nil {
		panic(err)
	}
	sliceProcessor.trie = tr

	log.Println("trie loaded")

	return sliceProcessor
}
