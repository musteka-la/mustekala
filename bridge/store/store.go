package store

import (
	logging "github.com/ipfs/go-log"
)

// Setup the logger as a package variable.
var log = logging.Logger("store")

// TODO
// everything will be in memory for now.
// soon we want to implement persistence of data here.
//
// we probably will use the IPFS client storage for content addressed data
//   in a such way that the libp2p component of the bridge uses the same store.
// and a redis database for location addressed data (i.e. coselector indexes).

// newstore
// getHash
// setHash
// ...
