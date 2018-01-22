package devp2p

// consumeToDevP2PChan is a wrapper to the for loop that processes the
// messages sent to the ToDevP2P channel.
func (m *Manager) consumeToDevP2PChan() {
	log.Debug("consuming todevp2p channel")

	// Easy to read shorthand
	toDevP2PChan := m.toDevP2PChan

	for {
		// blocks until it gets a message
		msg := <-toDevP2PChan
		// TODO
		// IMPLEMENT
		log.Debugf("incoming Message ToDevP2P: %v", msg)
	}
}
