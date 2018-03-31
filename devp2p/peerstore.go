package devp2p

import (
	"fmt"
	"math/big"
	"net"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/discover"
)

// peerStore keeps track of the devp2p peers after a succesful
// handshake (i.e. a match in protocols, version and fork).
// It is also able to give us the available nodes to make requests to.
type peerStore struct {
	lock sync.RWMutex

	peers map[string]*Peer

	sortedIndexByTD []*Peer
}

// Peer is an arrangement we use to organize the life cycle of a
// devp2p peer after it is dialed, and the encryption and protocol
// handshakes are performed.
// After a succesful ethereum handshake, this Peer can be added
// in the peerstore, to be managed for further requesting.
type Peer struct {
	// the id of the devp2p node, remote address and name
	id         discover.NodeID
	remoteAddr net.Addr
	name       string

	// the communication pipeline
	rw p2p.MsgReadWriter

	// head and total difficulty informed by the peer in the eth handshake
	currentBlock common.Hash
	td           *big.Int
	lock         sync.RWMutex

	// is this peer useful for us?
	byzantiumChecked bool
}

////////////////////////////////////////////////////////////////////////////////
// SORTING HELPERS
////////////////////////////////////////////////////////////////////////////////

type byTD []*Peer

func (b byTD) Len() int {
	return len(b)
}

func (b byTD) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byTD) Less(i, j int) bool {
	// reverse the sorting. We want to have at the head
	// of this index the peers with the most total difficulty
	if b[i].td.Cmp(b[j].td) == 1 {
		return true
	}
	return false
}

////////////////////////////////////////////////////////////////////////////////
// PEER
////////////////////////////////////////////////////////////////////////////////

// String implements fmt.Stringer.
func (p *Peer) String() string {
	return fmt.Sprintf("Peer %x %v", p.id[:8], p.remoteAddr)
}

// Head returns reported current block and total difficulty of the peer
func (p *Peer) Head() (currentBlock common.Hash, td *big.Int) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	copy(currentBlock[:], p.currentBlock[:])
	return currentBlock, new(big.Int).Set(p.td)
}

// RequestHeadersByHash fetches a batch of blocks' headers corresponding to the
// specified header query, based on the hash of an origin block.
func (p *Peer) RequestHeadersByHash(origin common.Hash, amount int, skip int, reverse bool) error {
	log.Debugf("Fetching batch of headers count %v from_hash 0x%x skip %v reverse %v", amount, origin[:8], skip, reverse)
	return p2p.Send(p.rw, GetBlockHeadersMsg, &getBlockHeadersData{Origin: hashOrNumber{Hash: origin}, Amount: uint64(amount), Skip: uint64(skip), Reverse: reverse})
}

// RequestHeadersByNumber fetches a batch of blocks' headers corresponding to the
// specified header query, based on the number of an origin block.
func (p *Peer) RequestHeadersByNumber(origin uint64, amount int, skip int, reverse bool) error {
	log.Debugf("Fetching batch of headers count %v from_number %v skip %v reverse %v", amount, origin, skip, reverse)
	return p2p.Send(p.rw, GetBlockHeadersMsg, &getBlockHeadersData{Origin: hashOrNumber{Number: origin}, Amount: uint64(amount), Skip: uint64(skip), Reverse: reverse})
}

////////////////////////////////////////////////////////////////////////////////
// PEERSTORE
////////////////////////////////////////////////////////////////////////////////

// newPeerStore prepares the peerstore of this library
func newPeerStore() *peerStore {
	return &peerStore{
		peers:           make(map[string]*Peer),
		sortedIndexByTD: make([]*Peer, 0),
	}
}

// add includes a peer into the peerstore and recalculate the
// latter's sorted index by total difficulty
func (p *peerStore) add(peer *Peer) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.peers[peer.String()] = peer

	p.sortedIndexByTD = append(p.sortedIndexByTD, peer)
	sort.Sort(byTD(p.sortedIndexByTD))

	log.Debugf("added peer to store %v %v %x", peer.String(), peer.name, peer.currentBlock[:8])
}

// remove excludes a peer from the peerstore, recalculating
// the sorting index by total difficulty
func (p *peerStore) remove(id string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	delete(p.peers, id)

	newSortedIndexByTD := make([]*Peer, 0)
	for _, sortedPeer := range p.sortedIndexByTD {
		if sortedPeer.String() != id {
			newSortedIndexByTD = append(newSortedIndexByTD, sortedPeer)
		}
	}
	p.sortedIndexByTD = newSortedIndexByTD

	log.Debug("removed peer from store", id)
}

// bestPeer returns the next best peer to send a request.
// Peers are sorted by total difficulty and picked by
// number of times a peer has been returned by this very function.
func (p *peerStore) bestPeer() *Peer {
	p.lock.Lock()
	defer p.lock.Unlock()

	// no peers? No Problem
	if len(p.peers) == 0 {
		return nil
	}

	// easier to read shorthand
	s := p.sortedIndexByTD

	// TODO
	// Keep rotating the best peer

	for i := 0; i < len(s)-1; i++ {
		if s[i].byzantiumChecked {
			return s[i]
		}
	}

	return nil
}
