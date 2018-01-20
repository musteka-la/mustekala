package devp2p

import (
	"fmt"
	"sync"
)

// networkStatus is the object that will give you the peers we have seen
// and some useful stats
type networkStatus struct {
	lock  sync.RWMutex
	peers map[string]*peerNetworkStatus
}

// peerNetworkStatus is the n-tuple of the networkStatus object
type peerNetworkStatus struct {
	status     string // We'll use an enum later
	statusPlus string // why it failed, observations, etc
}

// newNetworkStatus is the networkStatus constructor
func newNetworkStatus() *networkStatus {
	return &networkStatus{
		peers: make(map[string]*peerNetworkStatus),
	}
}

func (n *networkStatus) insert(p *Peer) {
	n.lock.Lock()
	defer n.lock.Unlock()

	n.peers[p.id] = &peerNetworkStatus{
		status: "encrypted-handshake",
	}

	// DEBUG
	for k, v := range n.peers {
		fmt.Printf("%v\t%v\n", k, v)
	}
	// DEBUG
}

func (n *networkStatus) updateStatus(p *Peer, status string, statusPlus string) {
	n.lock.Lock()
	defer n.lock.Unlock()

	n.peers[p.id].status = status
	n.peers[p.id].statusPlus = statusPlus

	// DEBUG
	for k, v := range n.peers {
		fmt.Printf("%v\t%v\n", k, v)
	}
	// DEBUG
}

func (n *networkStatus) dumpStatus() []*peerNetworkStatus {
	n.lock.Lock()
	defer n.lock.Unlock()

	response := make([]*peerNetworkStatus, 0)

	for _, v := range n.peers {
		response = append(response, v)
	}

	return response
}
