package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	gorp "gopkg.in/gorp.v1"
)

// Options keep db conn parameters
type Options struct {
	User     string
	Password string
	DBName   string
}

// LastBlock tracks the canonical chain.
// We index using the last inserted value, as there are
// always chances to rewrites.
type LastBlock struct {
	InsertedTS int64 `db:"inserted_ts"` // doubles as PK
	NumberId   int64 `db:"number_id"`
}

// WantFromDevp2p is the least of wanted data for our agents
// to retrieve from different devp2p clients
type WantFromDevp2p struct {
	InsertedTS    int64  `db:"inserted_ts"`     // doubles as PK
	Kind          string `db:"kind"`            // block body, tx receipt, etc
	Key           string `db:"key"`             // id or hash required
	LastRequestTS int64  `db:"last_request_ts"` // last time we sent a req for this key
	SuccessTS     int64  `db:"success_ts"`      // so we know to not ask again for it
}

// EthData is the data retrieved from the devp2p clients
// we try to keep for the longest time possible our data in
// here, to give the chance to agents to add it into several libp2p
// availabilities, however in the future we may want to prune
// this table.
type EthData struct {
	InsertedTS    int64  `db:"inserted_ts"`
	Kind          string `db:"kind"`             // block body, tx receipt, etc
	Hash          string `db:"hash"`             // ethereum hash id
	CID           string `db:"cid"`              // ipld cid, useful, so we compute it once
	Value         string `db:"value"`            // stored in stringed hex
	LastIPFSAddTS int64  `db:"last_ipfs_add_ts"` // last time a client tried to add it into IPFS
	IPFSSuccessTS int64  `db:"ipfs_success_ts"`  // so we know that we have add it at least once

}

// BlockTx is useful to find out whether we have all the
// transactions of a block, so we can prepare their
// representation as trie elements
type BlockTX struct {
	InsertedTS int64  `db:"inserted_ts"`
	BlockID    string `db:"block_id"`
	TxId       string `db:"tx_id"`
}

// BlockNumberofTx give us how many transactions a block has
type BlockNumberofTx struct {
	InsertedTS  int64  `db:"inserted_ts"`
	BlockID     string `db:"block_id"`
	NumberOfTxs int64  `db:"number_of_txs"`
}

// TxReceipt maps transactions against their receipts
type TxReceipts struct {
	InsertedTS   int64  `db:"inserted_ts"`
	TxId         string `db:"tx_id"`
	TxReceiptsId string `db:"tx_receipts_id"`
}

func InitDb(options Options) *gorp.DbMap {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		options.User, options.Password, options.DBName)

	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatalf("sql.Open failed %v", err)
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	dbmap.AddTableWithName(LastBlock{}, "lastblock")
	dbmap.AddTableWithName(WantFromDevp2p{}, "wantfromdevp2p")
	dbmap.AddTableWithName(EthData{}, "ethdata")
	dbmap.AddTableWithName(BlockTX{}, "blocktx")
	dbmap.AddTableWithName(BlockNumberofTx{}, "blocknumberoftx")
	dbmap.AddTableWithName(TxReceipts{}, "txreceipts")

	return dbmap
}
