package main

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	emptyCodeHash = crypto.Keccak256(nil)
	emptyRoot     = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
)

type Account struct {
	Nonce    uint64
	Balance  *big.Int
	Root     common.Hash // merkle root of the storage trie
	CodeHash []byte
}

func newAccount(balance string) *Account {
	balanceBigInt, _ := new(big.Int).SetString(balance, 10)

	return &Account{
		Nonce:    5,
		Balance:  balanceBigInt,
		Root:     emptyRoot,
		CodeHash: emptyCodeHash,
	}
}
