package devp2p

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/p2p"
)

func (m *Manager) handleBlockHeaderMsg(peer *Peer, msg *p2p.Msg) error {
	var headers []*types.Header
	if err := msg.Decode(&headers); err != nil {
		return fmt.Errorf("Error Decoding message: %v %v", msg, err)
	}

	// in protocolHandler() we did shot a request for the byzantium block.
	// if we have that response here, let's check its hash, and drop
	// the peer if it does not comply.
	if len(headers) == 1 && headers[0].Number.Cmp(ByzantiumBlockNumberBigInt) == 0 {
		// check hash
		if headers[0].Hash().String() == ByzantiumBlockHashStr {
			log.Debug("Peer byzantium block is OK")

			m.peerScrapper(peer.id, "50-byzantium block check passed", "OK") // hook

			// no need to ship this to the outgoing channel, we synchronize from here
			return nil
		} else {
			// if you are curious, most of the errors are due to the hash being
			// 0x6ff3e725355c909b52c5aa0e637e7c1d5e5b58bc25bc5fed320bf27278d81bd5
			// this is the one corresponding to Ethereum Classic (ETC)
			hashStr := fmt.Sprintf("%x", headers[0].Hash())

			m.peerScrapper(peer.id, "49-byzantium block check failed", hashStr)

			return fmt.Errorf("Peer byzantium block check failed, got %s", hashStr)
		}
	}

	// TODO
	// validate the sealing (pow) of each received header

	// TODO
	// ship the headers to the DB

	return nil
}
