package devp2p

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

// networkStatus is the object that will give you the peers we have seen
// and some useful stats
type networkStatus struct {
	lock  sync.RWMutex
	peers map[string]*peerNetworkStatus

	regex *regexMessages
}

// peerNetworkStatus is the n-tuple of the networkStatus object
type peerNetworkStatus struct {
	hash       string // the hash id's first 8 characters
	status     string // We'll use an enum later
	statusPlus string // why it failed, observations, etc
}

type regexMessages struct {
	statusMsg        *regexp.Regexp
	msgTooLarge      *regexp.Regexp
	decoding         *regexp.Regexp
	networkMismatch  *regexp.Regexp
	genesisMismatch  *regexp.Regexp
	protocolMismatch *regexp.Regexp
}

// newNetworkStatus is the networkStatus constructor
func newNetworkStatus() *networkStatus {
	rx := &regexMessages{
		statusMsg:        regexp.MustCompile(`^status message:`),
		msgTooLarge:      regexp.MustCompile(`^message too large:`),
		decoding:         regexp.MustCompile(`^decoding error:`),
		networkMismatch:  regexp.MustCompile(`^network mismatch:`),
		genesisMismatch:  regexp.MustCompile(`^genesis block mismatch:`),
		protocolMismatch: regexp.MustCompile(`^protocol version mismatch:`),
	}

	return &networkStatus{
		peers: make(map[string]*peerNetworkStatus),
		regex: rx,
	}
}

func (n *networkStatus) updateStatus(peerid string, status string, statusPlus string) {
	n.lock.Lock()
	defer n.lock.Unlock()

	// cleanup
	if peerid[:4] == "Peer" {
		peerid = peerid[5:]
	}

	// check that the peer exists in our registries, create it if not
	if _, ok := n.peers[peerid]; !ok {
		n.peers[peerid] = &peerNetworkStatus{}
	}

	n.peers[peerid].status = status

	switch {
	case status == "failed-eth-handshake":
		switch {
		case n.regex.statusMsg.MatchString(statusPlus):
			n.peers[peerid].statusPlus = "status message\",\"" + statusPlus[16:]
		case n.regex.msgTooLarge.MatchString(statusPlus):
			n.peers[peerid].statusPlus = "message too large\",\"" + statusPlus[19:]
		case n.regex.decoding.MatchString(statusPlus):
			n.peers[peerid].statusPlus = "rlp decoding\",\"" + statusPlus[16:]
		case n.regex.networkMismatch.MatchString(statusPlus):
			n.peers[peerid].statusPlus = "network mismatch\",\"" + statusPlus[18:]
		case n.regex.genesisMismatch.MatchString(statusPlus):
			n.peers[peerid].statusPlus = "genesis block mismatch\",\"" + statusPlus[24:]
		case n.regex.protocolMismatch.MatchString(statusPlus):
			n.peers[peerid].statusPlus = "p2p protocol mismatch\",\"" + statusPlus[27:]
		default:
			n.peers[peerid].statusPlus = statusPlus
		}
	default:
		n.peers[peerid].statusPlus = statusPlus
	}
}

func (n *networkStatus) dumpStatus() {
	n.lock.Lock()
	defer n.lock.Unlock()

	t := time.Now()
	fmt.Sprintf("%v", t.Format("20060102150405"))

	// TODO
	// This file code should not be here, output this to the FromDevP2P channel
	filepath := fmt.Sprintf("/tmp/network-status-%v.csv", t.Format("20060102150405"))
	f, err := os.Create(filepath)
	if err != nil {
		log.Errorf("Error creating the file! %v\n", err)
		return
	}
	defer f.Close()

	var line string
	for k, v := range n.peers {
		// split the id into hash and remote address
		id := strings.Split(k, " ")
		if len(id) != 2 {
			log.Errorf("peerid should be like <hash> <remoteAddr>. It is %v\n", k)
			return
		}
		// prepare the line and write it to the file
		line = fmt.Sprintf("\"%v\",\"%v\",\"%v\",\"%v\"\n", id[0], id[1], v.status, v.statusPlus)
		_, err := f.WriteString(line)
		if err != nil {
			log.Errorf("Error writing the file! %v\n", err)
			return
		}
	}
	log.Debug("Network status file created: ", filepath)
}
