package devp2p

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type regexMessages struct {
	statusMsg        *regexp.Regexp
	msgTooLarge      *regexp.Regexp
	decoding         *regexp.Regexp
	networkMismatch  *regexp.Regexp
	genesisMismatch  *regexp.Regexp
	protocolMismatch *regexp.Regexp
}

var peerScrapperRegex regexMessages

func init() {
	peerScrapperRegex = regexMessages{
		statusMsg:        regexp.MustCompile(`^status message:`),
		msgTooLarge:      regexp.MustCompile(`^message too large:`),
		decoding:         regexp.MustCompile(`^decoding error:`),
		networkMismatch:  regexp.MustCompile(`^network mismatch:`),
		genesisMismatch:  regexp.MustCompile(`^genesis block mismatch:`),
		protocolMismatch: regexp.MustCompile(`^protocol version mismatch:`),
	}
}

// peerScrapper is a hook to add an update in the database
// about the situation of a peer.
// this is the main feature of the service devp2p-node-scrapper.
// however it can be used as well in any other service by activating
// the flag IsPeerScrapperActive
func (m *Manager) peerScrapper(peerid, status, statusPlus string) {
	var err error

	if !m.config.IsPeerScrapperActive {
		return
	}

	// cleanup
	if peerid[:4] == "Peer" {
		peerid = peerid[5:]
	}
	peerid = strings.Replace(peerid, " ", ":", -1) // more redis key friendy

	// make additional processing to the status if needed
	// just csv friendly commas
	switch {
	case status == "39-ethereum handshake failed":
		switch {
		case peerScrapperRegex.statusMsg.MatchString(statusPlus):
			statusPlus = "status message, " + statusPlus[16:]
		case peerScrapperRegex.msgTooLarge.MatchString(statusPlus):
			statusPlus = "message too large, " + statusPlus[19:]
		case peerScrapperRegex.decoding.MatchString(statusPlus):
			statusPlus = "rlp decoding, " + statusPlus[16:]
		case peerScrapperRegex.networkMismatch.MatchString(statusPlus):
			statusPlus = "network mismatch, " + statusPlus[18:]
		case peerScrapperRegex.genesisMismatch.MatchString(statusPlus):
			statusPlus = "genesis block mismatch, " + statusPlus[24:]
		case peerScrapperRegex.protocolMismatch.MatchString(statusPlus):
			statusPlus = "p2p protocol mismatch, " + statusPlus[27:]
		default:
			statusPlus = statusPlus
		}
	}

	// storing into the database
	// we go without using abstractions because we want to use a single conn
	conn := m.dbPool.Get()
	defer conn.Close()

	// 0. send the peer to the set devp2p-peers (will create if it does not exist)
	_, err = conn.Do("SADD", "devp2p-scrapped-peers:all", peerid)
	if err != nil {
		fmt.Printf("Error setting value in redisDB: %v\n", err)
	}

	// 1. send the concatenated status to the set devp2p-peerstatuses:<peerid>
	concatStatus := fmt.Sprintf("%d, %s, %s", time.Now().UnixNano(), status, statusPlus)
	_, err = conn.Do("SADD", "devp2p-peerstatus:"+peerid, concatStatus)
	if err != nil {
		fmt.Printf("Error setting value in redisDB: %v\n", err)
	}

	// 2. if it is a byzantium peer, send it to the set devp2p-byzantium-peers
	if status == "50-byzantium block check passed" {
		_, err = conn.Do("SADD", "devp2p-byzantium-peers:all", peerid)
		if err != nil {
			fmt.Printf("Error setting value in redisDB: %v\n", err)
		}
	}

	// 3. if we have a negative for byzantium (wrong genesis, network or fork), send it to the set: devp2p-not-byzantium
	if status == "49-byzantium block check failed" {
		_, err = conn.Do("SADD", "devp2p-non-byzantium-peers:all", peerid)
		if err != nil {
			fmt.Printf("Error setting value in redisDB: %v\n", err)
		}
	}
}
