package migrator

type Migrator interface {
	CurrentSchemaVersion() int                        // Returns the current data-structure version
	CurrentFixupVersion() int                         // Returns the current fixup version
	TargetFixupVersion() int                          // Returns the highest fixup version available for the current data-structure
	UnderlyingObject() interface{}                    // Returns a reference to the underlying object which implements the interface
	Migrate() (Migrator, error)                       // Migrates to next data-structure version returning a new underlying object
	ApplyNextFixup() error                            // Applies next data fixup to the current data-structure version of the BuddyQueue
	ApplyChronicFixups() (dataChange bool, err error) // Applies all chronic fixups to the current data-structure version of the BuddyQueue
}
