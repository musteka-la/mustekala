package devp2p

import "sync"

// networkStatus is the object that will give you the peers we have seen
// and some useful stats
type networkStatus struct {
	lock  sync.RWMutex
	peers map[string]*peerNetworkStatus
}

// peerNetworkStatus is the n-tuple of the networkStatus object
type peerNetworkStatus struct {
}

// newNetworkStatus is the networkStatus constructor
func newNetworkStatus() *networkStatus {
	return &networkStatus{
		peers: make(map[string]*peerNetworkStatus),
	}
}

// insertOrUpdate is the main function of the networkStatus object
func (n *networkStatus) insertOrUpdate(p peerNetworkStatus) {
	n.lock.Lock()
	defer n.lock.Unlock()

}
