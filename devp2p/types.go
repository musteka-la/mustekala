package devp2p

import (
	"github.com/ethereum/go-ethereum/core/types"
)

//
type deliverHeaderMsg struct {
	PeerID  string
	Headers []*types.Header
}
