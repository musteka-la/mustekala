package lib

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var logFmt1 = "%22s | %24s | %27s | %s | %25s | %32s"
var logFmt2 = "%22s | %2s %10s %10s | %10s %16s |  %s/%s | %10s %5s %8s | %10s %10s %10s"

func init() {
	gob.Register(SliceMetadataToFile{})
}

func logSliceHeader() {
	log.Printf(logFmt1,
		"id",
		"stem (#, K & B time)",
		"state (t, B/E/L #)",
		"final",
		"storage (t, #sc, #nodes)",
		"bytes (stem, state, storage)",
	)
}

func logSliceData(sliceData *SliceData) {
	log.Printf(
		logFmt2,
		sliceData.id,
		fmt.Sprintf("%d", len(sliceData.stemKeys)),
		sliceData.stats["fetch-stem-keys"],
		sliceData.stats["fetch-stem-blobs"],
		sliceData.stats["fetch-slice"],
		fmt.Sprintf("%s/%s/%s %s",
			sliceData.stats["state-branches"],
			sliceData.stats["state-extensions"],
			sliceData.stats["state-leaves"],
			sliceData.stats["state-total-nodes"],
		),
		sliceData.stats["is-final"],
		sliceData.stats["max-depth"],
		sliceData.stats["fetch-smart-contract-trie-nodes"],
		sliceData.stats["smart-contracts"],
		sliceData.stats["smart-contract-trie-nodes"],
		sliceData.stats["bytes-stem"],
		sliceData.stats["bytes-slice"],
		sliceData.stats["bytes-storage"],
	)
}

type SliceMetadataToFile struct {
	NodePathId          string
	TotalStateNodes     string
	TotalStorageNodes   string
	TotalLeaves         string
	TotalSmartContracts string
	BytesStem           string
	BytesState          string
	BytesStorage        string
	Final               string
	MaxDepth            string
}

func (sp *SliceProcessor) storeSliceMetadata(path string, sliceData *SliceData) {
	var err error

	metadata := &SliceMetadataToFile{
		NodePathId:          sliceData.id,
		TotalStateNodes:     sliceData.stats["state-total-nodes"],
		TotalStorageNodes:   sliceData.stats["smart-contract-trie-nodes"],
		TotalLeaves:         sliceData.stats["state-leaves"],
		TotalSmartContracts: sliceData.stats["smart-contracts"],
		BytesStem:           sliceData.stats["bytes-stem"],
		BytesState:          sliceData.stats["bytes-slice"],
		BytesStorage:        sliceData.stats["bytes-storage"],
		Final:               sliceData.stats["is-final"],
		MaxDepth:            sliceData.stats["max-depth"],
	}

	// serialize the struct
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err = e.Encode(metadata)
	if err != nil {
		fmt.Println(`failed gob Encode`, err)
	}

	// create the filepath
	err = os.MkdirAll(sp.fileDir, 0755)
	if err != nil {
		panic(err)
	}

	// save the file
	err = ioutil.WriteFile(filepath.Join(sp.fileDir, path), b.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}

func (sp *SliceProcessor) querySliceData(path string) *SliceData {
	response := &SliceData{}

	f, err := os.OpenFile(filepath.Join(sp.fileDir, path), os.O_RDONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			panic(err)
		}
	}

	sliceMetadata := &SliceMetadataToFile{}
	gdec := gob.NewDecoder(f)
	err = gdec.Decode(&sliceMetadata)
	if err != nil {
		panic(fmt.Sprintf("failed go decode: %v", err))
	}

	response.id = sliceMetadata.NodePathId
	response.stats = make(map[string]string)
	response.stats["state-total-nodes"] = sliceMetadata.TotalStateNodes
	response.stats["smart-contract-trie-nodes"] = sliceMetadata.TotalStorageNodes
	response.stats["state-leaves"] = sliceMetadata.TotalLeaves
	response.stats["smart-contracts"] = sliceMetadata.TotalSmartContracts
	response.stats["bytes-stem"] = sliceMetadata.BytesStem
	response.stats["bytes-slice"] = sliceMetadata.BytesState
	response.stats["bytes-storage"] = sliceMetadata.BytesStorage
	response.stats["is-final"] = sliceMetadata.Final
	response.stats["max-depth"] = sliceMetadata.MaxDepth

	f.Close()

	return response
}

// GetHeatMap just takes the stored data and agregates it
// into a json file, for visualization purposes
func (sp *SliceProcessor) GetHeatMap(mode string) {
	if mode != "txt" {
		// more modes in the future
		return
	}

	log.Printf("generating heatmap data in mode %s", mode)

	// open the file or create it
	f, err := os.OpenFile(
		filepath.Join(sp.fileDir, "heatmap-data.txt"),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0x0000; i < 0x10000; i++ {
		path := fmt.Sprintf("%04x", i)

		// query the DB for this slice's file (keys and metadata)
		sliceData := sp.querySliceData(path)
		if sliceData == nil {
			log.Printf("can't find data for path %s. Run the slicer again")
			break
		}

		// append the data into the file
		line := fmt.Sprintf("%s %s%s %s %s %s %s %s %s\n",
			sliceData.id[:4],
			sliceData.stats["is-final"],
			sliceData.stats["max-depth"],
			sliceData.stats["state-total-nodes"],
			sliceData.stats["state-leaves"],
			sliceData.stats["smart-contracts"],
			sliceData.stats["smart-contract-trie-nodes"],
			sliceData.stats["bytes-slice"],
			sliceData.stats["bytes-storage"],
		)

		if _, err := f.Write([]byte(line)); err != nil {
			log.Fatal(err)
		}
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
