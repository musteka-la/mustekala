## Optimizing geth for the slicer (for fun and profit)

### Overview

We need to understand all we can about the way go-ethereum (geth from now)
optimizes its memory to manage its data. We should leverage from that,
and build our slice memory management overlay, so we can avoid as much
as possible interactions with the disk (which is way slower) at the time of
retrieving a slice (and the subsequent _deltas_ we want to get).

### Current State of the Art

#### debug.metrics

https://github.com/ethereum/go-ethereum/wiki/Metrics-and-Monitoring

(WIP)

#### pprof

(WIP)

#### the cache mentioned in the logs

	INFO [08-23|02:10:45.930] Imported new chain segment               blocks=1  txs=47   mgas=7.988   elapsed=278.571ms mgasps=28.674  number=6197430 hash=4d2fc5…e848dc cache=105.35mB

* https://github.com/ethereum/go-ethereum/blob/f34f361ca6635690f6dd81d6f3bddfff498e9fd6/core/blockchain.go#L1193
* https://github.com/ethereum/go-ethereum/blob/f34f361ca6635690f6dd81d6f3bddfff498e9fd6/trie/database.go#L731

```go
// Size returns the current storage size of the memory cache in front of the
// persistent database layer.
func (db *Database) Size() (common.StorageSize, common.StorageSize) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	// db.nodesSize only contains the useful data in the cache, but when reporting
	// the total memory consumption, the maintenance metadata is also needed to be
	// counted. For every useful node, we track 2 extra hashes as the flushlist.
	var flushlistSize = common.StorageSize((len(db.nodes) - 1) * 2 * common.HashLength)
	return db.nodesSize + flushlistSize, db.preimagesSize
}
```

### Historic links

#### The PR where they introduced the intermediate cache

* https://github.com/ethereum/go-ethereum/pull/15857

### Cheatsheet

#### geth command line option to add pprof and metrics

	--pprof --pprofaddr=0.0.0.0 --metrics

#### pprof command

	go tool pprof http://localhost:6060/debug/pprof/heap

#### memsize

(This is expensive, handle with care)

	curl -L -X POST http://localhost:6060/memsize/scan?root=node

#### bonus

	curl http://localhost:6060/debug/pprof/goroutine?debug=1