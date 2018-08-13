package lib

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/ethdb"
)

/*
	This file contains functions to retrieve values stored in geth's levelDB
	following certain specific key prefixes used
*/

type GethDB struct {
	db *ethdb.LDBDatabase
}

func NewGethDB(dbFilePath string) *GethDB {
	db, err := ethdb.NewLDBDatabase(dbFilePath, 768, 256)
	if err != nil {
		panic("can't open ethdb!")
	}

	return &GethDB{db: db}
}

func (g *GethDB) DB() *ethdb.LDBDatabase {
	return g.db
}

// GetCanonicalHash returns the stored CHT Hash for a given number
func (g *GethDB) GetCanonicalHash(number uint64) []byte {
	headerPrefix := []byte("h")
	numSuffix := []byte("n")
	encodedNumber := make([]byte, 8)
	binary.BigEndian.PutUint64(encodedNumber, number)

	key := append(append(headerPrefix, encodedNumber...), numSuffix...)
	val, _ := g.db.Get(key)

	return val
}

// GetHeaderRLP returns the RLP of the block header
// for a pair (hash, number) as key
func (g *GethDB) GetHeaderRLP(hash []byte, number uint64) []byte {
	headerPrefix := []byte("h")
	encodedNumber := make([]byte, 8)
	binary.BigEndian.PutUint64(encodedNumber, number)

	key := append(append(headerPrefix, encodedNumber...), hash...)

	val, _ := g.db.Get(key)
	return val
}
