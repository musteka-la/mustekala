package bridge

import (
	logging "github.com/ipfs/go-log"

	"github.com/metamask/eth-ipld-hot-importer/bridge/message"
)

// Setup the logger as a package variable.
var log = logging.Logger("bridge")

// Bridge is the center of our application.
// By starting and maintaining its channels, it provides a way
// for a node to pass the information it receives from the
// different networks where is is connected as peer.
type Bridge struct {
	Channels BridgeChannels
}

// BridgeChannels is the wrapping of the channels to and from the
// networks where the node is connected as peer. These channels
// must be initialized at construction time in the NewBridge function.
type BridgeChannels struct {
	FromDevP2P chan message.Message
	ToDevP2P   chan message.Message

	// TODO
	// FromLibP2P chan message.Message
	// ToLibP2P chan message.Message
}

// NewBridge sets up the store, initializes the channels and return the bridge object
func NewBridge(channelCapacity int) *Bridge {
	bridge := &Bridge{}

	/*
		bridge.Channels = BridgeChannels{
			FromDevP2P: make(chan Message, channelCapacity),
			ToDevP2P:   make(chan Message, channelCapacity),
		}

		log.Info("bridge set up. channels created. hot importer object allocated.")
	*/

	return bridge
}

// Start kicks off the loops that make the bridge alive
func (b *Bridge) Start() {
	log.Info("starting Bridge")

	// go b.consumeFromDevP2PChan()
	// go b.launchHotImporter()
}

// Stop finalizes the execution of the loops and gracefully shuts down the store
func (b *Bridge) Stop() {
	log.Info("stopping Bridge")

	// TODO
	// finalize the execution of llops
	// gracefully shut down the store
}
