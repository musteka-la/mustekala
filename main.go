package main

import (
	"github.com/metamask/eth-ipld-hot-importer/bridge"
	"github.com/metamask/eth-ipld-hot-importer/devp2p"
)

func main() {
	// Get the flags
	cfg := ParseFlags()

	// Setup the bridge
	// The bridge is the main piece of code of this importer, it controls
	// the incoming and outgoing channels, and houses the importing algorithms
	br := bridge.NewBridge(cfg.BridgeChannelCapacity)

	// Setup the devp2p server
	// Code has been written to help abstract the bridge developer of the
	// go-ethereum devp2p functionality:
	// * Discovery and maintaining of the devp2p peers.
	// * Messaging is managed using an
	//   incoming and an outgoing channel, using the peer requesting some
	//   piece of data, or "the best" one, based on determined criteria of the
	//   peerstore.
	devp2pConfig := &devp2p.Config{
		BootnodesPath:     cfg.DevP2PBootnodesPath,
		NodeDatabasePath:  cfg.DevP2PNodeDatabasePath,
		NetworkStatusPath: cfg.DevP2PNetworkStatusPath,
		LibP2PDebug:       cfg.DevP2PLibDebug,
	}

	devp2pServer := devp2p.NewManager(br, devp2pConfig)

	// Start the bridge
	// initiates the consuming of the bridge incoming channels
	go br.Start()

	// Start the devp2p server
	// * initiates the consuming of the devp2p incoming channel.
	// * starts the go-ethereum p2p server, making this node, an ethereum peer.
	//   * peers that pass the ethereum handshake will be added to a peerstore.
	//   * messages from the peers will be delivered to the outgoing channel.
	go devp2pServer.Start()

	// Start the importing of ethereum data
	// * runs the block header sync algorithm (TODO)
	// * runs the 4 level (65536) current block state trie importing algorithm (TODO)
	// * runs the on-demand requesting algorithm for state + storage data algorithm (TODO)
	// TODO

	// PLACEHOLDER
	// We run the main importing loop here, to block the main function
	// We should delete this line after we issue the start of the hot importer for loop.
	select {}
	// PLACEHOLDER
}
