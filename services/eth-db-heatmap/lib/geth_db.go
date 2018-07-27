package lib

import "encoding/binary"

/*
	This file contains functions to retrieve values stored in geth's levelDB
	following certain specific key prefixes used
*/

// get returns the value associated to that key in the DB
func (sp *SliceProcessor) get(key []byte) ([]byte, error) {
	return sp.db.Get(key)
}

// getCanonicalHash returns the stored CHT Hash for a given number
func (sp *SliceProcessor) getCanonicalHash(number uint64) []byte {
	headerPrefix := []byte("h")
	numSuffix := []byte("n")
	encodedNumber := make([]byte, 8)
	binary.BigEndian.PutUint64(encodedNumber, number)

	key := append(append(headerPrefix, encodedNumber...), numSuffix...)
	val, _ := sp.db.Get(key)

	return val
}

// getHeaderRLP returns the RLP of the block header
// for a pair (hash, number) as key
func (sp *SliceProcessor) getHeaderRLP(hash []byte, number uint64) []byte {
	headerPrefix := []byte("h")
	encodedNumber := make([]byte, 8)
	binary.BigEndian.PutUint64(encodedNumber, number)

	key := append(append(headerPrefix, encodedNumber...), hash...)

	val, _ := sp.db.Get(key)
	return val
}

// getBodyRLP returns the RLP of the block header, plus ommer list
// transactions for a pair (hash, number) as key
func (sp *SliceProcessor) getBodyRLP(hash []byte, number uint64) []byte {
	bodyPrefix := []byte("b")
	encodedNumber := make([]byte, 8)
	binary.BigEndian.PutUint64(encodedNumber, number)

	key := append(append(bodyPrefix, encodedNumber...), hash...)

	val, _ := sp.db.Get(key)
	return val
}
