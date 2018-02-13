package devp2p

import "regexp"

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

	var toDBKey, toDBValue string

	// cleanup
	if peerid[:4] == "Peer" {
		peerid = peerid[5:]
	}

	// make additional processing to the status, if needed
	switch {
	case status == "failed-eth-handshake":
		switch {
		case peerScrapperRegex.statusMsg.MatchString(statusPlus):
			statusPlus = "status message\",\"" + statusPlus[16:]
		case peerScrapperRegex.msgTooLarge.MatchString(statusPlus):
			statusPlus = "message too large\",\"" + statusPlus[19:]
		case peerScrapperRegex.decoding.MatchString(statusPlus):
			statusPlus = "rlp decoding\",\"" + statusPlus[16:]
		case peerScrapperRegex.networkMismatch.MatchString(statusPlus):
			statusPlus = "network mismatch\",\"" + statusPlus[18:]
		case peerScrapperRegex.genesisMismatch.MatchString(statusPlus):
			statusPlus = "genesis block mismatch\",\"" + statusPlus[24:]
		case peerScrapperRegex.protocolMismatch.MatchString(statusPlus):
			statusPlus = "p2p protocol mismatch\",\"" + statusPlus[27:]
		default:
			statusPlus = statusPlus
		}
	}

	// TODO
	// fill the key and value to the DB accordingly

	// TODO
	// (just in memory case)
	// Just fmt.Printf your key and value!

	// TODO
	// sending to the db
	_ = toDBKey
	_ = toDBValue
}
