package bridge

// consumeFromDevP2PChan is a wrapper to the for loop that processes the
// messages sent to the FromDevP2P channel.
func (b *Bridge) consumeFromDevP2PChan() {
	log.Debug("consuming fromdevp2p channel")

	// Easy to read shorthand
	fromDevP2PChan := b.Channels.FromDevP2P

	for {
		// Blocks until it gets a message
		msg := <-fromDevP2PChan
		// TODO
		// IMPLEMENT
		log.Debugf("incoming Message FromDevP2P: %v", msg)
	}
}
