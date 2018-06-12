package devp2p

import (
	"bufio"
	"fmt"
	"os"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/garyburd/redigo/redis"
	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("devp2p")

// Manager is the object that sets up the node as a devp2p peer,
// managing its connected peers, and communicating with the bridge
// through the assigned channels.
type Manager struct {
	config *Config

	// bootnodes slice
	bootnodes []*discover.Node

	// the actual peer connecting to devp2p
	server *p2p.Server

	// each service process has its own peerstore:
	// it is used to coordinate the data requests to devp2p
	peerstore *peerStore

	// this is our wrapper so we can get logging information
	// from github.com/ethereum/go-ethereum/p2p
	p2pLibLogger gethlog.Logger

	// mustekala services database connection pool
	dbPool *redis.Pool

	// this one give us block headers
	deliverHeaderCh chan deliverHeaderMsg

	// the blck header syncer
	syncer *Syncer
}

// Config is the configuration object for DevP2P
type Config struct {
	// bootnodes file
	BootnodesPath string

	// node database path. Must be appointed outside this package
	NodeDatabasePath string

	// we can log what is going on with the go-ethereum/p2p library
	LibP2PDebug bool

	// Passing labels from one object to another
	DbPool *redis.Pool

	// activate this value to send updates on discovered and connected
	// peers to the database
	IsPeerScrapperActive bool

	// activate this value to start the block header syncing
	IsSyncBlockHeaderActive bool

	// the client's private key here
	PrivateKeyFilePath string
}

// NewManager returns a DevP2P Manager object
//
// * defines bootnodes, and
// * defines node database, to be used by the
//   discovery functions of go-ethereum/p2p library.
//
// * sets up the mustekala services database connection.
//
// * sets up the peerstore, which , among other things,
//   keeps track of the best available peer to send a message.
//
// * sets up a custom p2p lib logger, to be able to get stats
//   on dialing, encrypted and proto handshakes,
//   as well as catching and printing those logs.
//
// * sets up the server. See manager.protocolHandler() and handlerIncomingMsg()
//   to understand a peer's lifecycle and the receiving of
//   devp2p messages and sending to the bridge.
//
// * sets up the block header syncer, if the option is enabled.
func NewManager(config *Config) *Manager {
	var err error

	manager := &Manager{}

	manager.config = config

	manager.bootnodes, err = parseBootnodesFile(config.BootnodesPath)
	if err != nil {
		log.Error("processBootnodesFile", err)
		os.Exit(1)
	}

	if config.NodeDatabasePath == "" {
		log.Error("node database path must be appointed outside this package")
		os.Exit(1)
	}

	manager.dbPool = config.DbPool

	manager.peerstore = newPeerStore()

	manager.p2pLibLogger = &p2pLibLogger{mgr: manager}

	manager.server = manager.newServer()

	if config.IsSyncBlockHeaderActive {
		manager.deliverHeaderCh = make(chan deliverHeaderMsg, 1)
		manager.syncer = manager.NewSyncer()
	}

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

	if m.config.IsSyncBlockHeaderActive {
		log.Info("starting block header sync")

		go m.syncer.Start()
	}
}

// Stop terminates the server
func (m *Manager) Stop() {
	m.server.Stop()
}

// BestPeer makes available the best peer from the store
func (m *Manager) BestPeer() *Peer {
	return m.peerstore.bestPeer()
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
