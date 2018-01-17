package devp2p

import (
	"fmt"

	"github.com/ethereum/go-ethereum/p2p"
)

// protocolHandler controls the lifecycle of a connected peer.
// It creates the peer as an entity in this level of abstraction, to be able
// to add it into our peerstore if it succeeds the ethereum handshake.
// Also, it establishes a permanent loop to manage its incoming messages.
// On error, the permanent loop closes, and the peer is removed from the peerstore.
func (m *Manager) protocolHandler(p *p2p.Peer, rw p2p.MsgReadWriter) error {
	// this peer is formatted as an eth peer
	ethPeer := &Peer{
		id: p.String(),
		rw: rw,
	}

	if err := ethPeer.sendStatusMsg(); err != nil {
		log.Debug("failed eth protocol handshake", p, "error", err)
		return err
	}

	// in the lifecycle of a peer, after the ethereum handshake is succesful,
	// we add this peer into our store, which will make them indirectly available
	// to the caller of the manager, to do requests to the devp2p network.
	// once the loop for this connection (implemented below) is broke, the peer
	// will be removed from the store.
	m.peerstore.add(ethPeer)
	defer m.peerstore.remove(ethPeer)

	// we don't want peers that aren't in the byzantium hard fork.
	// we will send a message to the peer asking for its block
	// we have logic inside handleIncomingMsg() to deal with this specific block header.
	p2p.Send(ethPeer.rw,
		GetBlockHeadersMsg,
		&getBlockHeadersData{
			Origin: hashOrNumber{
				Number: uint64(ByzantiumBlockNumber)},
			Amount:  uint64(1),
			Skip:    uint64(0),
			Reverse: true,
		})

	// this is a permanent loop, it waits for the p2p library to ReadMsg()
	// and then switches over the code of the message (New block, Get receipts, etc, etc)
	// for further processing. This loop is broken at the first error,
	// triggering the removal of the peer from our store also.
	for {
		if err := m.handleIncomingMsg(ethPeer); err != nil {
			log.Debug("failed ethereum message handling", "peer", ethPeer.id, "err", err)
			return err
		}
	}
}

// handleIncomingMsg manages the message received from a peer.
// As this library is a simple abstraction for a caller to communicate
// with the DevP2P network, messages in this function get repackaged
// for their adecquate handling at the caller layer.
func (m *Manager) handleIncomingMsg(peer *Peer) error {

	msg, err := peer.rw.ReadMsg()
	if err != nil {
		return err
	}
	defer msg.Discard()

	protocolMaxMsgSize := uint32(10 * 1024 * 1024)
	if msg.Size > protocolMaxMsgSize {
		return fmt.Errorf("message too large: %v > %v", msg.Size, protocolMaxMsgSize)
	}

	switch {
	case msg.Code == StatusMsg:
		// ethereum handshake has been done. We should not be asked about status afterwards.
		return fmt.Errorf("got a status message after handshake")

	case msg.Code == NewBlockHashesMsg:
		log.Debug("NewBlockHashes", "peer", peer.id)
		return fmt.Errorf("not Implemented")

	// this is the Broadcast message of a Transaction
	case msg.Code == TxMsg:
		log.Debug("Tx", "peer", peer.id)
		return fmt.Errorf("not Implemented")

	case msg.Code == GetBlockHeadersMsg:
		log.Debug("GetBlockHeaders", "peer", peer.id)
		return m.handleGetBlockHeaderMsg(peer, &msg)

	case msg.Code == BlockHeadersMsg:
		log.Debug("BlockHeaders", "peer", peer.id)
		return m.handleBlockHeaderMsg(peer, &msg)

	case msg.Code == GetBlockBodiesMsg:
		log.Debug("GetBlockBodies", "peer", peer.id)
		return fmt.Errorf("not Implemented")

	case msg.Code == BlockBodiesMsg:
		log.Debug("BlockBodies", "peer", peer.id)
		return fmt.Errorf("not Implemented")

	// This is the Broadcast message of a Block
	case msg.Code == NewBlockMsg:
		log.Debug("NewBlock", "peer", peer.id)
		return fmt.Errorf("not Implemented")

	case msg.Code == GetNodeDataMsg:
		log.Debug("GetNodeData", "peer", peer.id)
		return fmt.Errorf("not Implemented")

	case msg.Code == NodeDataMsg:
		log.Debug("NodeData", "peer", peer.id)
		return fmt.Errorf("not Implemented")

	case msg.Code == GetReceiptsMsg:
		log.Debug("GetReceipts", "peer", peer.id)
		return fmt.Errorf("not Implemented")

	case msg.Code == ReceiptsMsg:
		log.Debug("Receipts", "peer", peer.id)
		return fmt.Errorf("not Implemented")

	default:
		return fmt.Errorf("message code not supported")
	}
}
