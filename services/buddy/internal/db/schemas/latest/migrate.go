package latest

import (
	"github.com/donglei1234/platform/services/buddy/internal/db/schemas/migrator"
)

func (bq *BuddyQueue) CurrentSchemaVersion() int {
	return bq.SchemaVersion
}

func (bq *BuddyQueue) CurrentFixupVersion() int {
	return bq.FixupVersion
}

func (bq *BuddyQueue) TargetFixupVersion() int {
	return len(v1Fixups)
}

func (bq *BuddyQueue) UnderlyingObject() interface{} {
	return bq
}

func (bq *BuddyQueue) Migrate() (migrator.Migrator, error) {
	return bq, migrator.ErrUnableToMigrateFromLatest
}

func (bq *BuddyQueue) ApplyNextFixup() error {
	if bq.FixupVersion < len(v1Fixups) {
		if err := v1Fixups[bq.FixupVersion](bq); err != nil {
			return err
		} else {
			bq.FixupVersion = bq.FixupVersion + 1
			return nil
		}
	} else {
		return migrator.ErrNoFixupToApply
	}
}

func (bq *BuddyQueue) ApplyChronicFixups() (dataChange bool, err error) {
	return false, nil
}
