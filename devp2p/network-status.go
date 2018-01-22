package devp2p

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
)

/*
	Network Status File

	Status Codes:

	00		TCP Dialing
		01	TCP Dialing failed
		29	Connection Setup failed
		39	Ethereum Handshake failed
		49 	Byzantium Block Handshake failed
	50 		Byzantium Block Handshake succeed

*/

// networkStatus is the object that will give you the peers we have seen
// and some useful stats
type networkStatus struct {
	lock  sync.RWMutex
	peers map[string]*peerNetworkStatus

	regex *regexMessages

	filepath string
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
func newNetworkStatus(filepath string) *networkStatus {
	rx := &regexMessages{
		statusMsg:        regexp.MustCompile(`^status message:`),
		msgTooLarge:      regexp.MustCompile(`^message too large:`),
		decoding:         regexp.MustCompile(`^decoding error:`),
		networkMismatch:  regexp.MustCompile(`^network mismatch:`),
		genesisMismatch:  regexp.MustCompile(`^genesis block mismatch:`),
		protocolMismatch: regexp.MustCompile(`^protocol version mismatch:`),
	}

	return &networkStatus{
		peers:    make(map[string]*peerNetworkStatus),
		regex:    rx,
		filepath: filepath,
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

	// no file defined? go home
	if n.filepath == "" {
		return
	}

	f, err := os.OpenFile(n.filepath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Errorf("Error opening/creating the file! %v\n", err)
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
		line = fmt.Sprintf("\"%v\",\"%v\",\"%v\",\"%v\",\"%v\"\n", id[0], id[1], v.status[0:2], v.status[3:], v.statusPlus)
		_, err := f.WriteString(line)
		if err != nil {
			log.Errorf("Error writing the file! %v\n", err)
			return
		}
	}
	log.Info("Writing to the network status file ", n.filepath)
}
