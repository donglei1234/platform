package v0

import (
	"github.com/donglei1234/platform/services/buddy/internal/db/schemas/latest"
	"github.com/donglei1234/platform/services/buddy/internal/db/schemas/migrator"
)

// "version zero" of BuddyQueue data has never been used. This exists only for example purposes and this
// package can be deleted after we conduct the first actual BuddyQueue version migration. (when latest becomes v1)

func (bq *BuddyQueue) CurrentSchemaVersion() int {
	return bq.SchemaVersion
}

func (bq *BuddyQueue) CurrentFixupVersion() int {
	return bq.FixupVersion
}

func (bq *BuddyQueue) TargetFixupVersion() int {
	return len(v0Fixups)
}

func (bq *BuddyQueue) UnderlyingObject() interface{} {
	return bq
}

func (bq *BuddyQueue) Migrate() (migrator.Migrator, error) {
	versionTwoBuddyQueue := latest.NewBuddyQueue()
	return versionTwoBuddyQueue, nil
}

func (bq *BuddyQueue) ApplyNextFixup() error {
	return migrator.ErrNoFixupToApply
}

func (bq *BuddyQueue) ApplyChronicFixups() (dataChange bool, err error) {
	return false, nil
}
