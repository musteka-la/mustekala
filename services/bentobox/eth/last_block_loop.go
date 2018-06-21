package eth

import (
	"database/sql"
	"log"
	"time"

	"github.com/metamask/mustekala/services/bentobox/db"
)

// we put this here for aesthetic purposes
// EXPLAIN:
// We could not just used "GROUP BY", because we needed the latest timestamp
//   per highest block number found.
// This query,
// * Gets the latest 100 block number found into tmp0
// * Adds a pos column, which ranks inside each block number per timestamp
//   (this is the GROUP BY analogy)
// * Gets from this new tmp1, the top ranked timestamps per number_id
// * And finally returns the top block id.
const lastBlockSQLQuery = `
SELECT number_id, inserted_ts
FROM (
	SELECT number_id, inserted_ts, rank()
		OVER (
			PARTITION BY tmp0.number_id
			ORDER BY tmp0.inserted_ts DESC
		) AS pos
	FROM (
			SELECT inserted_ts, number_id
			FROM lastblock
			ORDER BY inserted_ts DESC
			LIMIT 100
	) tmp0
) tmp1
WHERE pos <= 1
ORDER BY number_id DESC
LIMIT 1;
`

// LastBlockLoop is a simple loop that polls the eth json rpc
// for the latest block (height of the network), getting a
// block number in return, which stores alongside the timestamp
// of such received number.
// Eventually this service will be horizontally scaled,
// meaning that we must retrieve several "last block" values
// from the shared db and determine whether we can count with
// an actual number, or we should poll again.
func (e *EthManager) LastBlockLoop() {
	var err error

	log.Printf("Starting LastBlockLoop")

	for {
		needToPoll := false

		// what do we have here?
		var lastDbBlock *db.LastBlock
		err = e.dbMap.SelectOne(&lastDbBlock, lastBlockSQLQuery)
		if err != nil {
			if err == sql.ErrNoRows {
				// no rows? OK, let's poll for one
				log.Printf("Lastblock tuples not found in the DB, issuing a poll")
				needToPoll = true
			} else {
				log.Printf("Error on SQL query for last block: %v", err)
				time.Sleep(500 * time.Millisecond)
				continue
			}
		} else {
			// found a value
			// find out _when_ we inserted this tuple in the DB
			now := time.Now().UnixNano()

			if now-lastDbBlock.InsertedTS > int64(e.pollIntervalMS)*1000*1000 {
				// time difference exceeded configured poll interval, issuing a poll
				// we are not logging as this event should be very frequent
				needToPoll = true
			}
		}

		// so, we poll of the flag was activated above
		if needToPoll {
			// do the actual query here
			response, err := e.getNetworkHeight()
			if err != nil {
				log.Printf("There was an error requesting the last block, %v", err)
				time.Sleep(500 * time.Millisecond)
				continue
			}

			// avoid below storage code, if is the same last block as in memory
			if response == lastDbBlock.NumberId {
				time.Sleep(500 * time.Millisecond)
				continue
			}

			// store the response in the database
			lastBlockTuple := db.LastBlock{
				InsertedTS: time.Now().UnixNano(),
				NumberId:   response,
			}

			log.Printf("Inserting new block found: %v", response)

			if err := e.dbMap.Insert(&lastBlockTuple); err != nil {
				log.Printf("Error inserting last block tuple %v: %v", lastBlockTuple, err)
			}

			// we good, keep going
		}

		// Avoid the dreaded all-devouring loop
		time.Sleep(e.pollIntervalMS * time.Millisecond)
	}
}
