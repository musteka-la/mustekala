package bridge

import "sync"

// Here will be the logic to the actual hot importing
// from main.go
// * runs the block header sync algorithm (TODO)
// * runs the 4 level (65536) current block state trie importing algorithm (TODO)
// * runs the on-demand requesting algorithm for state + storage data algorithm (TODO)
//
// Basically, We build a single function wrapping that launches these four goroutines.
// In a first stage, all the data will be in memory. As we grow a beard,
// We manage in-memory data, as well as a connection to the libp2p network (IPFS),
// So we become good citizens and reply to the requests for data of our fellow devp2p peers.

const (
	MIN_BLOCK_NUMBER = 0 // Initially is the genesis block. Later we may want to have it at byzantium (4,370,000).
	//HEADER_LOOP_ACTIVE_REQUESTS     = 100 // Maximum number of requests to do on a header loop iteration.
	//HEADER_LOOP_HEADERS_PER_REQUEST = 21  // How many headers I ask a peer for?
	//HEADER_LOOP_CONFIRMATIONS_PIVOT = 3   // If I get N blocks over a valid one, I'll say it is confirmed.
	//HEADER_LOOP_RESET               = 100 // How many blocks do we want to go back to, if the sync stalled.
)

// headerSyncStatus works the following way:
//
// * requested block header is one we asked the networks to fetch for us, example, block 156,034.
// * arrived block header means that we have it, useful for the loop to evaluate it.
// * valid means that the block header information is consistent with its hash, and its parent corresponds
//   to what we have.
// * confirmed is equivalent to the popular blockchain UX, and is associated to the number of blocks
//   that have this block header as an ancestor.
type headerSyncStatus byte

const (
	requested = headerSyncStatus(iota)
	arrived
	valid
	confirmed
)

// hotImporter is the data structure for the Hot Importer
type hotImporter struct {
	lock        sync.RWMutex
	headers     map[int]header // All the arrived block headers mapped by their block number
	pivotHeader int
}

// header represents a block header
type header struct {
	syncStatus       headerSyncStatus
	numberOfConfirms int
	blockHash        []byte
	parentHash       []byte
	timeStamp        int
}

// launchHotImporterLoops is the heart of the Hot Importer Application
func (b *Bridge) launchHotImporter() {
	log.Info("launching hot importer")

	// init the hot importer variables
	b.hotImporter = &hotImporter{
		headers:     make(map[int]header),
		pivotHeader: MIN_BLOCK_NUMBER,
	}

	// get the blockchain headers
	go b.headerLoop()
}
