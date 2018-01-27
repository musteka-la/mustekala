package bridge

import "time"

// headerInfuraLoop is the bucle that synchronizes the block headers from
// a minimum block (i.e genesis, dao, byzantium) until the present.
func (b *Bridge) headerInfuraLoop() {
	log.Info("Launching the Infura Block Header Syncing Loop")

	// easy to read shorthand
	//toDevP2PChan := b.Channels.ToDevP2P

	for {

		// PLACEHOLDER
		// we don't wanna block no program with no empty loop yo!
		time.Sleep(1 * time.Second)
		// PLACEHOLDER
	}
}
