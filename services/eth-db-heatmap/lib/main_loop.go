package lib

import (
	"fmt"
	"log"
)

// SliceTheTrie is the main loop of this program.
// One of its optimizations is that the program will try to
// see whether the slice has been computer before and stored
// in the DB so we spend O(1) on get it again.
func (sp *SliceProcessor) SliceTheTrie() {
	var err error

	log.Println("init the main loop: Slice The Trie!")

	logSliceHeader()

	// we won't do R-04 as it is currently 69,905 uniformily
	// distributed branches

	for i := 0x0000; i < 0x10000; i++ {
		path := fmt.Sprintf("%04x", i)

		// query the DB for this slice's file (keys and metadata)
		sliceData := sp.querySliceData(path)
		if sliceData != nil {
			log.Printf("found slice for path %s, skipping", path)
			continue
		}

		// compute and store the slice data from the DB
		sliceData, err = sp.getSliceData(path, 10)
		if err != nil {
			// let's make sure we stop to check what happened
			panic(err)
		}

		// print a descriptive line for this slice
		logSliceData(sliceData)

		// store the metadata in its file
		sp.storeSliceMetadata(path, sliceData)
	}
}
