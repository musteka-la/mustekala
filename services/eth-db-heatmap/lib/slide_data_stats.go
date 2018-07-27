package lib

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

var emptyCodeHash = crypto.Keccak256(nil)

type Account struct {
	Nonce    uint64
	Balance  *big.Int
	Root     common.Hash // merkle root of the storage trie
	CodeHash []byte
}

// getBELstats counts the number of branches, extensions and leaves
// in your slice
func (sd *SliceData) getBELstats() {
	branches := 0
	extensions := 0
	leaves := 0

	for i, _ := range sd.sliceBlobs {
		for j, _ := range sd.sliceBlobs[i] {
			switch sd.identifyTrieNode(sd.sliceBlobs[i][j]) {
			case "branch":
				branches++
			case "extension":
				extensions++
			case "leaf":
				leaves++
			}
		}
	}

	sd.stats["state-branches"] = fmt.Sprintf("%d", branches)
	sd.stats["state-extensions"] = fmt.Sprintf("%d", extensions)
	sd.stats["state-leaves"] = fmt.Sprintf("%d", leaves)
	sd.stats["state-total-nodes"] = fmt.Sprintf("%d", branches+extensions+leaves)
	sd.stats["smart-contracts"] = fmt.Sprintf("%d", len(sd.storageRootKeys))
}

// isFinal tells you whether the slice have children, by returning false
// and which is its maximum depth.
func (sd *SliceData) isFinal() {
	var maxDepth int
	foundBEinLastRow := false

	// find the maxDepth first
	for i := len(sd.sliceBlobs) - 1; i >= 0; i-- {
		if l := len(sd.sliceBlobs[i]); l > 0 {
			maxDepth = i
			break
		}
	}

	for _, node := range sd.sliceBlobs[maxDepth] {
		if t := sd.identifyTrieNode(node); t != "leaf" {
			foundBEinLastRow = true
			break
		}
	}

	if foundBEinLastRow {
		sd.stats["is-final"] = "T"
	} else {
		sd.stats["is-final"] = "F"
	}

	sd.stats["max-depth"] = fmt.Sprintf("%02d", maxDepth)
}

func (sd *SliceData) identifyTrieNode(input []byte) string {
	var i []interface{}
	var account Account

	err := rlp.DecodeBytes(input, &i)
	if err != nil {
		// zero tolerance
		panic(err)
	}

	switch len(i) {
	case 2:
		first := i[0].([]byte)
		last := i[1].([]byte)

		switch first[0] / 16 {
		case '\x00':
			fallthrough
		case '\x01':
			return "extension"
		case '\x02':
			fallthrough
		case '\x03':
			err = rlp.DecodeBytes(last, &account)
			if err != nil {
				panic(err)
			}
			if !bytes.Equal(account.CodeHash, emptyCodeHash) {
				sd.storageRootKeys = append(sd.storageRootKeys, account.Root)
			}
			return "leaf"
		default:
			panic("unknown hex prefix on trie node")
		}

	case 17:
		return "branch"

	default:
		panic("unknown trie node type")
	}
}

func (sd *SliceData) getStorageTrieStats(trieDB *trie.Database) {
	totalStorageNodeCount := 0
	byteCount := 0

	// we count the bulk operation
	_start := time.Now().UnixNano()

	for _, root := range sd.storageRootKeys {
		storageNodeCount := 0

		tr, err := trie.NewSecure(root, trieDB, 0)
		if err != nil {
			panic(err)
		}

		for it := trie.NewSliceIterator(tr, nil); it.Next(true); {
			if !(it.Hash() == common.Hash{}) {
				storageNodeCount++
				byteCount += len(it.Blob(it.Hash())) + 32
			}
		}

		totalStorageNodeCount += storageNodeCount
	}

	sd.stats["fetch-smart-contract-trie-nodes"] = humanizeTime(time.Now().UnixNano() - _start)
	sd.stats["smart-contract-trie-nodes"] = strconv.Itoa(totalStorageNodeCount)
	sd.stats["bytes-storage"] = humanize.Bytes(uint64(byteCount))
}
