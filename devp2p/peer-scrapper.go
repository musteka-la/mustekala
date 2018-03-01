package devp2p

import (
	"fmt"
	"regexp"
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
	if !m.config.IsPeerScrapperActive {
		return
	}

	var toDBKey, toDBVal string

	// cleanup
	if peerid[:4] == "Peer" {
		peerid = peerid[5:]
	}

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

	// fill the key and value to the DB accordingly
	toDBKey = "peer " + peerid
	toDBVal = fmt.Sprintf("%d, %s, %s", time.Now().UnixNano(), status, statusPlus)

	// DEBUG
	fmt.Printf("peerScrapper:\nkey:\t%v\nval:\t%v\n", toDBKey, toDBVal)
	// DEBUG

	// TODO
	// 0. Send the peer to the set devp2p-peers (will create if it does not exist)
	// 1. Send the status to the set devp2p-peerstatus:<peerid>
	// 2. If it is a byzantium peer, send it to the set devp2p-byzantium-peers
	// 3. If we have a negative for byzantium (wrong genesis, network or fork), send it to the set: devp2p-not-byzantium
}
