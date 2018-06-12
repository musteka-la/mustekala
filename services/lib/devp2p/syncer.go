package devp2p

import (
	"time"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"

	"github.com/metamask/mustekala/services/lib/db"
	"github.com/metamask/mustekala/services/lib/devp2p/downloader"
)

// Syncer is the coordinator of the block header syncer service.
type Syncer struct {
	manager    *Manager
	downloader *downloader.Downloader
}

// NewSyncer is the Syncer constructor
func (m *Manager) NewSyncer() *Syncer {
	mode := downloader.FastSync               // hack
	chaindb, _ := ethdb.NewMemDatabase()      // hack
	eventMux := new(event.TypeMux)            // hack
	blockchain := db.LoadBlockChain(m.dbPool) // not so much

	return &Syncer{
		manager:    m,
		downloader: downloader.New(mode, chaindb, eventMux, blockchain, m.peerstore.remove),
	}
}

// Start the main function of this service.
func (s *Syncer) Start() {
	go func() {
		for {
			// Consume the deliver headers channel
			msg := <-s.manager.deliverHeaderCh

			// Send it through the pipeline
			s.downloader.DeliverHeaders(msg.PeerID, msg.Headers)
		}
	}()

	for {
		// best peer
		bestPeer := s.manager.peerstore.bestPeer()

		// no peers no nothing
		if bestPeer == nil {
			// make sure to wait one step
			time.Sleep(1 * time.Second)

			continue
		}

		head, td := bestPeer.Head()

		// TODO
		// Must be bound to a channel of sorts from devp2p library
		s.downloader.RegisterPeer(bestPeer.String(), 63, bestPeer)

		// Start the synchronization from that peer
		s.downloader.Synchronise(bestPeer.String(), head, td, downloader.FastSync)

		// TODO
		s.downloader.UnregisterPeer(bestPeer.String())

		time.Sleep(1 * time.Second)
	}
}
