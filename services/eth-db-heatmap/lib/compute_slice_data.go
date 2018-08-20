package lib

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/trie"
)

// SliceData has all the juicy bits you need about computing
// an ethereum slice.
type SliceData struct {
	it *trie.SliceIterator

	// how i am identifying this slice for all purposes
	id string

	// from state root to slice head
	stemKeys  []common.Hash
	stemBlobs [][]byte

	// the slice itself. nodes are grouped by depth
	sliceKeys  [][]common.Hash
	sliceBlobs [][][]byte

	// for storage metadata purposes
	storageRootKeys []common.Hash

	// such as retrieval times; branches, extensions, and leaves; sizes.
	stats map[string]string
}

// getSliceData gets your slice and all its related info
func (sp *SliceProcessor) getSliceData(head string, depth int) (*SliceData, error) {
	// validate params
	if !isSliceHeadValid(head) {
		return nil, errors.New(
			fmt.Sprintf("expected `blank` or a sequence of lowercase nibbles, got `%s` instead", head))
	}
	sliceHead := sliceHeadToKeyBytes(head)

	// init the output object

	output := &SliceData{
		id:              getSliceId(sliceHead, depth, sp.stateRoot),
		stemKeys:        make([]common.Hash, 0),
		stemBlobs:       make([][]byte, 0),
		sliceKeys:       make([][]common.Hash, depth+1),
		sliceBlobs:      make([][][]byte, depth+1),
		storageRootKeys: make([]common.Hash, 0),
		stats:           make(map[string]string),
	}

	// get the stem
	output.it = trie.NewSliceIterator(sp.trie, sliceHead)
	output.it.Next(true)
	output.fetchStemKeys()
	output.fetchStemBlobs()

	// get the slice
	output.it = trie.NewSliceIterator(sp.trie, sliceHead)
	output.fetchSlice(depth)

	// branches/extensions/leaves
	output.getBELstats()

	// is it final? (do we need to get its children?)
	// what is its max depth?
	output.isFinal()

	// count how many storage tries
	// (in the slices with smart contracts)
	output.getStorageTrieStats(sp.trieDB)

	// we are done here, to the pneumatic tube!
	return output, nil
}
