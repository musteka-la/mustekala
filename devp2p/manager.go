package devp2p

import (
	"bufio"
	"fmt"
	"os"
	"time"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
	logging "github.com/ipfs/go-log"

	"github.com/metamask/eth-ipld-hot-importer/bridge/message"
)

var log = logging.Logger("devp2p")

// Manager is the object that sets the node as a devp2p peer,
// managing its connected peers, and communicating with the bridge
// through the assigned channels.
type Manager struct {
	config *Config

	server *p2p.Server

	toDevP2PChan   chan message.Message
	fromDevP2PChan chan message.Message

	peerstore *peerStore

	p2pLibLogger  gethlog.Logger
	networkStatus *networkStatus
}

// Config is the configuration object for DevP2P
type Config struct {
	// bootnodes file
	BootnodesPath string

	// bootnodes slice
	bootnodes []*discover.Node

	// node database path. Must be appointed outside this package
	NodeDatabasePath string

	// we can find the client's private key here
	PrivateKeyFilePath string

	// if you want to log the status of the network
	NetworkStatusPath string

	// we can log what is going on with the go-ethereum/p2p library
	LibP2PDebug bool

	// Incoming and Outgoing Channels
	IncomingChan chan message.Message
	OutgoingChan chan message.Message
}

// NewDevP2P returns a DevP2P Manager object
//
// * defines bootnodes, and
// * defines node database, to be used by the go-ethereum/p2p library.
//
// * assigns the bridge channels: to and from devp2p.
//
// * sets up the peerstore, which , among other things,
//   keeps track of the best available peer to send a message.
//
// * sets up the network status object, we can retrieve some stats here
//   for further analysis.
//
// * sets up a custom p2p lib logger, to be able to get stats on dialing, encrypted
//   and proto handshakes, as well as catching and printing those logs.
//
// * sets up the server. See manager.protocolHandler() and handlerIncomingMsg()
//   to understand a peer's lifecycle and the receiving of
//   devp2p messages and sending to the bridge.
func NewManager(config *Config) *Manager {
	var err error

	manager := &Manager{}

	manager.config = config

	config.bootnodes, err = parseBootnodesFile(config.BootnodesPath)
	if err != nil {
		log.Error("processBootnodesFile", err)
		os.Exit(1)
	}

	if config.NodeDatabasePath == "" {
		log.Error("node database path must be appointed outside this package")
		os.Exit(1)
	}

	manager.toDevP2PChan = config.IncomingChan
	manager.fromDevP2PChan = config.OutgoingChan

	manager.peerstore = newPeerStore()

	manager.networkStatus = newNetworkStatus(config.NetworkStatusPath)

	manager.p2pLibLogger = &p2pLibLogger{mgr: manager}

	manager.server = manager.newServer(config)

	return manager
}

// Start should be run as a goroutine
func (m *Manager) Start() {
	log.Info("starting devp2p node")

	go func() {
		if err := m.server.Start(); err != nil {
			log.Error("devp2p server", err)
			os.Exit(1)
		}
	}()

	log.Info("launching fromDevP2P channel consumer")

	go m.consumeToDevP2PChan()

	if m.config.NodeDatabasePath != "" {
		log.Info("launching network status csv writer")

		go func() {
			for {
				m.networkStatus.dumpStatus()

				time.Sleep(60 * time.Second)
			}
		}()
	}
}

// Stop terminates the server
func (m *Manager) Stop() {
	m.server.Stop()
}

// parseBootnodesFile parses the bootnodes file to be included in the
// devp2p server.
func parseBootnodesFile(filePath string) ([]*discover.Node, error) {
	nodes := []*discover.Node{}

	if filePath == "" {
		return nil, fmt.Errorf("A bootnodes file must be defined!")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		nodeUrl := scanner.Text()
		node, err := discover.ParseNode(nodeUrl)
		if err != nil {
			log.Error("add bootstrap node error", nodeUrl, err)
		}
		nodes = append(nodes, node)
		log.Debug("added bootstrap node", nodeUrl)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nodes, nil
}
