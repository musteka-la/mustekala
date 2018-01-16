package devp2p

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/p2p"
)

// handshakeTimeout is how much we want to wait for the ethereum handshake
// to be completed.
var handshakeTimeout = 5 * time.Second

// statusData is the information to be sent and received at the
// ethereum handshake.
type statusData struct {
	ProtocolVersion uint32
	NetworkId       uint32
	TD              *big.Int
	CurrentBlock    common.Hash
	GenesisBlock    common.Hash
}

// getOurStatusData prepares the information we send to the clients
// we are performing our ethereum handshake.
//
// In future versions of this package, we may want to send a better
// status than the one of a node at the genesis block.
func getOurStatusData() *statusData {
	_td := new(big.Int)
	td, _ := _td.SetString("17179869184", 10)

	ourStatus := &statusData{
		ProtocolVersion: uint32(63),
		NetworkId:       uint32(1),
		TD:              td,
		CurrentBlock:    common.HexToHash("d4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3"),
		GenesisBlock:    common.HexToHash("d4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3"),
	}

	return ourStatus
}

// sendStatusMsg initiates the ethereum handshake, sending the status message.
func (p *Peer) sendStatusMsg() error {
	errc := make(chan error, 2)
	ourStatus := getOurStatusData()

	go func() {
		errc <- p2p.Send(p.rw, StatusMsg, ourStatus)
	}()
	go func() {
		errc <- p.readStatusMsg()
	}()
	timeout := time.NewTimer(handshakeTimeout)
	defer timeout.Stop()
	for i := 0; i < 2; i++ {
		select {
		case err := <-errc:
			if err != nil {
				return err
			}
		case <-timeout.C:
			return p2p.DiscReadTimeout
		}
	}
	return nil
}

// readStatusMsg deals with the received status from the contacted peer.
func (p *Peer) readStatusMsg() error {
	// we start by reading the message
	msg, err := p.rw.ReadMsg()
	if err != nil {
		return err
	}

	// ... and find where anything can have gone wrong
	if msg.Code != StatusMsg {
		return fmt.Errorf("status message: first msg has code %x (!= %x)", msg.Code, StatusMsg)
	}
	protocolMaxMsgSize := uint32(10 * 1024 * 1024)
	if msg.Size > protocolMaxMsgSize {
		return fmt.Errorf("message too large: %v > %v", msg.Size, protocolMaxMsgSize)
	}

	// decode the handshake and make sure everything matches
	var theirStatus statusData
	ourStatus := getOurStatusData()
	if err := msg.Decode(&theirStatus); err != nil {
		return fmt.Errorf("decoding error: %v %v", msg, err)
	}

	// ethereum handshake checks: genesis block, network and version
	if theirStatus.GenesisBlock != ourStatus.GenesisBlock {
		return fmt.Errorf("genesis block mismatch: %x (!= %x)",
			theirStatus.GenesisBlock[:8],
			ourStatus.GenesisBlock[:8])
	}
	if theirStatus.NetworkId != ourStatus.NetworkId {
		return fmt.Errorf("network mismatch: %d (!= %d)",
			theirStatus.NetworkId,
			ourStatus.NetworkId)
	}
	if int(theirStatus.ProtocolVersion) != int(ourStatus.ProtocolVersion) {
		return fmt.Errorf("protocol version mismatch: %d (!= %d)",
			theirStatus.ProtocolVersion,
			ourStatus.ProtocolVersion)
	}

	// no match means no errors. We are good to go. Get the values we need.
	p.td, p.currentBlock = theirStatus.TD, theirStatus.CurrentBlock
	return nil
}
