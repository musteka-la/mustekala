package eth

import (
	"database/sql"
	"log"
	"time"

	"github.com/metamask/mustekala/services/bentobox/db"
)

// we put this here for aesthetic purposes
// EXPLAIN:
//
const wantedElementsSQLQuery = `
UPDATE wantfromdevp2p
SET last_request_ts = $1
WHERE inserted_ts IN (
	SELECT inserted_ts
	FROM wantfromdevp2p
	WHERE
		$1-last_request_ts>=$2
		AND
		success_ts=0
	LIMIT $3
	FOR UPDATE SKIP LOCKED
)
RETURNING kind, key;
`

const updateSuccessTSSQLQuery = `
UPDATE wantfromdevp2p
SET success_ts = $3
WHERE
	kind = $1
	AND
	key = $2;
`

// RpcDispatcherLoop obtains ethereum data from the clients
// by reading the "wantfromdevp2p" table and dispatching
// queries. Succesful results are stored into the "ethdata"
// table, for further processing
func (e *EthManager) RpcDispatcherLoop() {
	var err error

	log.Printf("Starting RpcDispatcherLoop")

	for {
		wantedElementsCount := e.maxQueries - e.qm.getQueueCount()

		if wantedElementsCount < 0 {
			// wait until this clears
			time.Sleep(500 * time.Millisecond)
			continue
		}

		// query for wanted elements in the table
		// criteria:
		// - last_request was made after "redo time" seconds
		// - haven't had success
		var wantedElements []*db.WantFromDevp2p

		_, err = e.dbMap.Select(&wantedElements,
			wantedElementsSQLQuery,
			time.Now().UnixNano(),
			e.redoQueryTime*1000*1000*1000,
			wantedElementsCount)

		if err != nil {
			if err != sql.ErrNoRows {
				log.Printf("Error on SQL query for last block: %v", err)
			}

			time.Sleep(500 * time.Millisecond)
			continue
		}

		// dispatch in parallel
		for _, _wantedItem := range wantedElements {
			// we need to clone this value when doing async voodoo
			wantedItem := _wantedItem

			go e.dispatcher(wantedItem.Kind, wantedItem.Key)
		}

		// avoid the dreaded all-devouring loop
		time.Sleep(500 * time.Millisecond)
	}
}

// dispatcher encapsulates the rpc call and the subsequent
// process of the obtained value
func (e *EthManager) dispatcher(kind, key string) {
	// add it to our query manager
	qmKey := "kind" + "_" + "key"
	e.qm.addQuery(qmKey)

	value, err := e.rpcCall(kind, key)
	if err != nil {
		log.Printf("Error on RPC Call (%v) (%v)", kind, key)
		return
	}

	if err = e.processEthData(kind, key, value); err != nil {
		log.Printf("Error on Eth Data processing (%v) (%v)", kind, key)
		return
	}

	// we good?
	// remove the kind/key from the manager
	e.qm.removeQuery(qmKey)

	// write the sucess timestamp into the DB
	_, err = e.dbMap.Exec(
		updateSuccessTSSQLQuery,
		kind,
		key,
		time.Now().UnixNano())
	if err != nil {
		log.Printf("Erorr updating success_ts in (%v) (%v)", kind, key)
	}
}

// rpcCall switches by kind to get the data from the ethereum client
func (e *EthManager) rpcCall(kind, key string) (string, error) {
	// DEBUG
	log.Printf("DEBUG: Invoking the calls for %v %v", kind, key)
	// DEBUG

	return "", nil
}

// processEthDAta switches by kind of element to store the
// obtained content in the DB, for further processing.
// (i.e. making it available to IPFS)
func (e *EthManager) processEthData(kind, key, value string) error {
	// DEBUG
	log.Printf("DEBUG: Processing the result for %v %v", kind, key)
	// DEBUG

	return nil
}
