package latest

var v1Fixups []fixup

type fixup = func(p *BuddyQueue) error

func init() {
	v1Fixups = []fixup{
		// THESE MUST STAY IN ORDER FOREVER, PLEASE
		// ONLY ADD TO THE BOTTOM OF THIS LIST
		fixDefaultBuddySettings, // 0 to 1
	}
}

// schema version 1, fixup version 0 to 1
func fixDefaultBuddySettings(bq *BuddyQueue) error {

	// If bq is an existing BuddyQueue, the AllowToBeAdded of Settings property would initially be false.
	// If bq is a new BuddyQueue, the AllowToBeAdded of Settings property would initially be true.
	//
	// This fixup makes sure that the current Character can be added as buddy by others by default.

	bq.Settings = &BuddySettings{
		AllowToBeAdded: true,
	}

	return nil
}
