package bridge

import "time"

// headerLoop is the bucle that synchronizes the block headers from
// a minimum block (i.e genesis, dao, byzantium) until the present.
func (b *Bridge) headerLoop() {
	log.Info("Launching the Header Loop")

	// easy to read shorthand
	//toDevP2PChan := b.Channels.ToDevP2P

	for {
		// TODO
		// processing arrived items

		// TODO
		// determine what is needed to request

		// TODO
		// do we know the capacity of the devp2p node?
		// too much time since last time we asked. Ask again, come back later.

		// TODO
		// launch block header requests to devp2p based on what we need and capacity reported

		// PLACEHOLDER
		// we don't wanna block no program with no empty loop yo!
		time.Sleep(1 * time.Second)
		// PLACEHOLDER
	}
}
