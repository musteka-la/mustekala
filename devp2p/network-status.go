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

}

func (n *networkStatus) updateStatus(p *Peer, status string, statusPlus string) {
	n.lock.Lock()
	defer n.lock.Unlock()

	n.peers[p.id].status = status
	n.peers[p.id].statusPlus = statusPlus
}

func (n *networkStatus) dumpStatus() {
	n.lock.Lock()
	defer n.lock.Unlock()

	for k, v := range n.peers {
		// Super sofisticated CSV lib
		fmt.Printf("\"%v\",\"%v\",\"%v\"\n", k, v.status, v.statusPlus)
	}
}
