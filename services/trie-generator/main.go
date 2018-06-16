package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

func main() {
	// get the flags
	cfg := ParseFlags()

	// get your redisDB
	redisdb := NewRedisDB(cfg.DatabaseConn)

	// create a new trie structure
	trieDB := trie.NewDatabase(redisdb)
	t, err := trie.New(common.Hash{}, trieDB)
	checkError(err)

	// welcome message
	fmt.Println("Trie Generator")

	start := time.Now()

	// input Loop
	for i := 0; i < cfg.NumberOfAccounts; i++ {
		// create your account
		idBalance := fmt.Sprintf("1%06d000000", i) // two birds...
		account := newAccount(idBalance)

		// do the transformations
		value, _ := rlp.EncodeToBytes(account)
		secureKey := crypto.Keccak256([]byte(idBalance))

		// add it into the trie
		err := t.TryUpdate(secureKey, value)
		checkError(err)
	}

	// commit your work into the DB
	_, err = t.Commit(nil)
	checkError(err)

	err = trieDB.Commit(t.Hash(), false)
	checkError(err)

	// how much did it take?
	duration := time.Since(start)
	fmt.Printf("Done. Added %d accounts in %.f miliseconds\n", cfg.NumberOfAccounts, duration.Seconds()*1000)

	// return the root
	fmt.Printf("\nYour root is %x\n", t.Hash())
}

func checkError(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(0)
	}
}
