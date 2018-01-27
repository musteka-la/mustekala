package message

// Message is the structure to pass information between the peer
// component of the node in a network and the bridge.
// The network components should import this library in order to be
// able to send and receive messages over the channels.
type Message struct {
	Header  string
	Payload interface{}
}
