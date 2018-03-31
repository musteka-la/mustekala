package devp2p

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/p2p"
)

func (m *Manager) handleEmptyResponseMsg(peer *Peer, msg *p2p.Msg) error {
	// TODO
	// Just empty, not even headers
	//
	headers := make([]*types.Header, 0)
	return p2p.Send(peer.rw, BlockHeadersMsg, headers)
}
