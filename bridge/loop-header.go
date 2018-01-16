package bridge

import "time"

// headerLoop is the bucle that synchronizes the block headers from
// a minimum block (i.e genesis, dao, byzantium) until the present.
//
// Its main activities are:
// * ?
func (b *Bridge) headerLoop() {
	log.Info("Launching the Header Loop")

	// Easy to read shorthand
	//toDevP2PChan := b.Channels.ToDevP2P

	for {
		// TODO
		// Processing arrived items

		// TODO
		// Determine what is needed to request

		// TODO
		// Do we know the capacity of the devp2p node?
		// Too much time since last time we asked. Ask again, come back later.

		// TODO
		// Launch block header requests to devp2p based on what we need and capacity reported

		// PLACEHOLDER
		// We don't wanna block no program with no empty loop yo!
		time.Sleep(1 * time.Second)
		// PLACEHOLDER
	}
}
