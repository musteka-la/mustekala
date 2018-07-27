package lib

import (
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/ethereum/go-ethereum/common"
)

func (sd *SliceData) fetchStemKeys() {
	_start := time.Now().UnixNano()

	// the actual fetching
	stem := sd.it.StemKeys()

	sd.stats["fetch-stem-keys"] = humanizeTime(time.Now().UnixNano() - _start)

	for _, item := range stem {
		sd.stemKeys = append(sd.stemKeys, item)
	}
}

func (sd *SliceData) fetchStemBlobs() {
	_start := time.Now().UnixNano()

	// the actual fetching
	stem := sd.it.StemBlobs()

	sd.stats["fetch-stem-blobs"] = humanizeTime(time.Now().UnixNano() - _start)

	byteCount := 0

	for _, item := range stem {
		sd.stemBlobs = append(sd.stemBlobs, item)
		byteCount += len(item) + 32
	}

	sd.stats["bytes-stem"] = humanize.Bytes(uint64(byteCount))
}

func (sd *SliceData) fetchSlice(depth int) {
	_start := time.Now().UnixNano()

	// the actual fetching
	keys, blobs := sd.it.Slice(depth, true)

	sd.stats["fetch-slice"] = humanizeTime(time.Now().UnixNano() - _start)

	byteCount := 0

	for i, _ := range keys {

		sd.sliceKeys[i] = make([]common.Hash, 0)
		sd.sliceBlobs[i] = make([][]byte, 0)

		for j, _ := range keys[i] {
			sd.sliceKeys[i] = append(sd.sliceKeys[i], keys[i][j])
			sd.sliceBlobs[i] = append(sd.sliceBlobs[i], blobs[i][j])

			byteCount += len(keys[i][j]) + len(blobs[i][j])
		}
	}

	sd.stats["bytes-slice"] = humanize.Bytes(uint64(byteCount))
}
