package main

import (
	"github.com/metamask/eth-ipld-hot-importer/bridge"
	"github.com/metamask/eth-ipld-hot-importer/devp2p"
)

func main() {
	// get the flags
	cfg := ParseFlags()

	// setup the bridge
	// the bridge is the main piece of code of this importer, it controls
	// the incoming and outgoing channels, and houses the importing algorithms
	br := bridge.NewBridge(cfg.BridgeChannelCapacity)

	// setup the devp2p server
	// code has been written to help abstract the bridge developer of the
	// go-ethereum devp2p functionality:
	// * discovery and maintaining of the devp2p peers.
	// * messaging is managed using an
	//   incoming and an outgoing channel, using the peer requesting some
	//   piece of data, or "the best" one, based on determined criteria of the
	//   peerstore.
	devp2pConfig := &devp2p.Config{
		BootnodesPath:     cfg.DevP2PBootnodesPath,
		NodeDatabasePath:  cfg.DevP2PNodeDatabasePath,
		NetworkStatusPath: cfg.DevP2PNetworkStatusPath,
		LibP2PDebug:       cfg.DevP2PLibDebug,
		IncomingChan:      br.Channels.ToDevP2P,
		OutgoingChan:      br.Channels.FromDevP2P,
	}

	devp2pServer := devp2p.NewManager(devp2pConfig)

	// start the bridge
	// initiates the consuming of the bridge incoming channels
	go br.Start()

	// start the devp2p server
	// * initiates the consuming of the devp2p incoming channel.
	// * starts the go-ethereum p2p server, making this node, an ethereum peer.
	//   * peers that pass the ethereum handshake will be added to a peerstore.
	//   * messages from the peers will be delivered to the outgoing channel.
	go devp2pServer.Start()

	// TODO
	// a nice SIGINT manager
	select {}
}
