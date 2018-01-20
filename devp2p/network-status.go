package devp2p

import (
	"fmt"
	"os"
	"sync"
	"time"
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

	t := time.Now()
	fmt.Sprintf("%v", t.Format("20060102150405"))

	// TODO
	// This file code should not be here, output this to the FromDevP2P channel
	f, err := os.Create(fmt.Sprintf("/tmp/network-status-%v.csv", t.Format("20060102150405")))
	if err != nil {
		fmt.Printf("Error creating the file! %v\n", err)
		return
	}
	defer f.Close()

	var line string
	for k, v := range n.peers {
		line = fmt.Sprintf("\"%v\",\"%v\",\"%v\"\n", k, v.status, v.statusPlus)
		_, err := f.WriteString(line)
		if err != nil {
			fmt.Printf("Error writing the file! %v\n", err)
			return
		}
	}
}
