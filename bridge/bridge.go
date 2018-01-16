package bridge

import (
	logging "github.com/ipfs/go-log"
)

// Setup the logger as a package variable.
var log = logging.Logger("bridge")

// Bridge is the center of our application.
// By starting and maintaining its channels, it provides a way
// for a node to pass the information it receives from the
// different networks where is is connected as peer.
type Bridge struct {
	Channels BridgeChannels

	hotImporter *hotImporter
}

// BridgeChannels is the wrapping of the channels to and from the
// networks where the node is connected as peer. These channels
// must be initialized at construction time in the NewBridge function.
type BridgeChannels struct {
	FromDevP2P chan Message
	ToDevP2P   chan Message
}

// Message is the structure to pass information between the peer
// component of the node in a network and the bridge.
// The network components should import this library in order to be
// able to send and receive messages over the channels.
type Message struct {
	Header  string
	Payload interface{}
}

// NewBridge initializes the channels and returns a bridge
func NewBridge(channelCapacity int) *Bridge {
	bridge := &Bridge{}

	bridge.Channels = BridgeChannels{
		FromDevP2P: make(chan Message, channelCapacity),
		ToDevP2P:   make(chan Message, channelCapacity),
	}

	log.Info("bridge set up. channels created. hot importer object allocated.")

	return bridge
}

// Start kicks off the go routines that consume the incoming channels
// (in plural, as we are thinking on the libp2p abstraction channels in the future)
func (b *Bridge) Start() {
	log.Info("starting Bridge")

	go b.consumeFromDevP2PChan()
	go b.launchHotImporter()
}
