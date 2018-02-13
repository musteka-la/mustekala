package devp2p

import (
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/p2p"
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
	// the id of the devp2p node, shorted to 8 chars.
	id string

	// the communication pipeline
	rw p2p.MsgReadWriter

	// total difficulty informed by the peer in the eth handshake
	td *big.Int

	// current block informed by the peer in the eth handshake
	currentBlock common.Hash

	////////////////
	// Sent and Received Requests Data
	////////////////
	requestCntLock      sync.RWMutex
	nextPeerTimesPicked int // Times it has been returned by the function nextPeer()
	sentRequestsCnt     int
	receivedRequestsCnt int
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
// PEERSTORE CONSTRUCTOR
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

	p.peers[peer.id] = peer

	p.sortedIndexByTD = append(p.sortedIndexByTD, peer)
	sort.Sort(byTD(p.sortedIndexByTD))

	log.Debug("added peer to store", peer.id)
}

// rRmove excludes a peer from the peerstore, recalculating
// the sorting index by total difficulty
func (p *peerStore) remove(peer *Peer) {
	p.lock.Lock()
	defer p.lock.Unlock()

	delete(p.peers, peer.id)

	newSortedIndexByTD := make([]*Peer, 0)
	for _, sortedPeer := range p.sortedIndexByTD {
		if sortedPeer.id != peer.id {
			newSortedIndexByTD = append(newSortedIndexByTD, sortedPeer)
		}
	}
	p.sortedIndexByTD = newSortedIndexByTD

	log.Debug("removed peer from store", peer.id)
}

// capacity returns the number of peers available to make requests,
// and the timestamp the request was made.
func (p *peerStore) capacity() (int, int64) {
	p.lock.Lock()
	defer p.lock.Unlock()

	return len(p.peers), time.Now().UnixNano()
}

// nextPeer returns the next best peer to send a request.
// Peers are sorted by total difficulty and picked by
// number of times a peer has been returned by this very function.
func (p *peerStore) nextPeer() *Peer {
	p.lock.Lock()
	defer p.lock.Unlock()

	// easier to read shorthand
	s := p.sortedIndexByTD

	for i := 0; i < len(s)-1; i++ {
		if s[i].nextPeerTimesPicked < s[i+1].nextPeerTimesPicked {
			s[i].nextPeerTimesPicked++
			return s[i]
		}
	}

	s[len(s)-1].nextPeerTimesPicked++
	return s[len(s)-1]
}
