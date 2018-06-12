package devp2p

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/p2p"
)

func (m *Manager) handleGetBlockHeaderMsg(peer *Peer, msg *p2p.Msg) error {
	// for now, we only respond with an empty header slice.
	// in the future we want to respond with info from our bridge.
	headers := make([]*types.Header, 0)
	return p2p.Send(peer.rw, BlockHeadersMsg, headers)
}
