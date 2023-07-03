package db

import (
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/buddy/internal/db/schemas/latest"
	"github.com/donglei1234/platform/services/buddy/internal/db/schemas/migrator"
	"github.com/donglei1234/platform/services/buddy/internal/db/schemas/v0"
	"github.com/donglei1234/platform/services/common/jsonx"
	"github.com/donglei1234/platform/services/common/nosql/document"
)

// FixupBuddyQueue runs BuddyQueues which are already on the "latest" schema through the appropriate fixups. (both
// one-time fixups and chronic fixups)
func (d *Database) FixupBuddyQueue(buddyQueue *BuddyQueue) error {
	var internalErr error

	if updateErr := buddyQueue.Update(func() bool {
		if fixupApplied, e := d.applyFixupsForCurrentVersion(buddyQueue); e != nil {
			d.logger.Error("error encountered during BuddyQueue fixup",
				zap.Int("impacted_buddy_queue_version", latest.Version),
				zap.Error(e),
			)
			internalErr = migrator.ErrFixupFailure
			return false
		} else if fixupApplied {
			// Data was changed and no fixup errors were encountered. Let's save!
			return true
		} else {
			// No data was changed, so no need to save
			internalErr = ErrNoDataChange
			return false
		}
	}); updateErr != nil {
		if internalErr == ErrNoDataChange {
			return nil
		} else {
			d.logger.Error(
				"error encountered during BuddyQueue document update",
				zap.NamedError("updateErr", updateErr),
				zap.NamedError("internalErr", internalErr),
			)
			return ErrFixupFailed
		}
	} else {
		return nil
	}
}

// migrateToLatest runs the BuddyQueue through both schema migrations as well as all appropriate fix-ups
func (d *Database) migrateToLatest(m migrator.Migrator) (migratedMigrator migrator.Migrator, migrated bool, err error) {
	currentVer := m.CurrentSchemaVersion()

	// If the BuddyQueue is a legacy version, we walk it through the version migration process until it is caught up
	for currentVer < latest.Version {
		migrated = true

		if _, err := d.applyFixupsForCurrentVersion(m); err != nil {
			d.logger.Error("fixup error encountered",
				zap.Int("impacted_buddy_queue_version", m.CurrentSchemaVersion()),
				zap.Int("impacted_fixup_version", m.CurrentFixupVersion()),
				zap.Error(err),
			)
			return nil, false, migrator.ErrFixupFailure
		} else if m, err = m.Migrate(); err != nil {
			d.logger.Error(
				"version migration error encountered",
				zap.Int("impacted_buddy_queue_version", currentVer),
				zap.Error(err),
			)
			return nil, false, migrator.ErrVersionMigrationFailure
		} else if m.CurrentSchemaVersion() <= currentVer {
			// If no forward progress was made, we bail out of migration attempt
			d.logger.Error(
				"version migration failed to progress",
				zap.Int("impacted_buddy_queue_version", currentVer+1),
				zap.Int("resulting_buddy_queue_version", m.CurrentSchemaVersion()),
			)
			return nil, false, migrator.ErrUnexpectedVersionValue
		} else {
			currentVer = m.CurrentSchemaVersion()
		}
	}

	migratedMigrator = m
	return migratedMigrator, migrated, nil
}

func (d *Database) migratorAsBuddyQueueDoc(
	m migrator.Migrator,
	docKey document.Key,
	docVersion document.Version,
) (*BuddyQueue, error) {
	var bq *latest.BuddyQueue

	switch v := m.UnderlyingObject().(type) {
	case *latest.BuddyQueue:
		bq = v
	default:
		return nil, migrator.ErrUnexpectedNonLatestVersion
	}

	if bq == nil {
		return nil, migrator.ErrUnexpectedNilUnderlyingValue
	} else {
		buddyQueue := &BuddyQueue{}
		buddyQueue.Document.InitWithVersion(&buddyQueue.BuddyQueue, buddyQueue.Clear, d.ros, docKey, docVersion)
		buddyQueue.BuddyQueue = *bq
		return buddyQueue, nil
	}
}

// applyFixupsForCurrentVersion applies the appropriate set of fixups for the current schema version because
// Each version of the BuddyQueue may have its own set of fixups
func (d *Database) applyFixupsForCurrentVersion(m migrator.Migrator) (fixupApplied bool, err error) {
	currentFixupVer := m.CurrentFixupVersion()
	targetFixupVer := m.TargetFixupVersion()

	// We run this version of the BuddyQueue through all one-off fixups
	for currentFixupVer < targetFixupVer {
		fixupApplied = true

		if err := m.ApplyNextFixup(); err != nil {
			return false, err
		} else if m.CurrentFixupVersion() <= currentFixupVer {
			// If no forward progress was made, we bail out of apply-fixups attempt
			return false, migrator.ErrUnexpectedFixupVersionValue
		} else {
			currentFixupVer = m.CurrentFixupVersion()
		}
	}

	// Chronic fixups are run once for legacy BuddyQueues (during version migration) and always for latest.BuddyQueue
	if dataChanged, err := m.ApplyChronicFixups(); err != nil {
		return false, err
	} else if dataChanged == true {
		fixupApplied = true
	}

	return fixupApplied, nil
}

func (d *Database) loadBuddyQueueAsMigrator(docKey document.Key) (migrator.Migrator, document.Version, error) {
	var unknownVerBuddyQueue []byte

	if version, err := d.ros.Get(
		docKey,
		document.WithDestination(&unknownVerBuddyQueue),
		document.WithAnyVersion(),
	); err != nil {
		return nil, 0, err
	} else if len(unknownVerBuddyQueue) == 0 {
		return nil, 0, ErrUnexpectedEmptyBuddyQueue
	} else {
		// All versions of the BuddyQueue (past and future) should contain a Version field at root
		BuddyQueueVerNum := jsoniter.Get(unknownVerBuddyQueue, "Version").ToInt()

		if m, err := d.unmarshalByBuddyQueueVersion(BuddyQueueVerNum, unknownVerBuddyQueue); err != nil {
			d.logger.Error(
				"BuddyQueue unmarshal error",
				zap.Int("version_field_value", BuddyQueueVerNum),
				zap.Error(err),
			)
			return nil, 0, ErrBuddyQueueUnmarshalFailure
		} else {
			return m, version, nil
		}
	}
}

func (d *Database) unmarshalByBuddyQueueVersion(BuddyQueueVersion int, data []byte) (migrator.Migrator, error) {
	if BuddyQueueVersion > latest.Version {
		return nil, ErrBuddyQueueVersionExceedsLatestVersion
	} else if BuddyQueueVersion == latest.Version {
		return d.unmarshalAsLatest(data)
	} else if BuddyQueueVersion == v0.Version {
		return d.unmarshalAsV0(data)
	} else {
		return nil, ErrUnhandledBuddyQueueVersion
	}
}

func (d *Database) unmarshalAsLatest(data []byte) (*latest.BuddyQueue, error) {
	latestBuddyQueue := &latest.BuddyQueue{}

	if err := jsonx.Unmarshal(data, latestBuddyQueue); err != nil {
		return nil, err
	} else {
		return latestBuddyQueue, nil
	}
}

func (d *Database) unmarshalAsV0(data []byte) (*v0.BuddyQueue, error) {
	v0BuddyQueue := &v0.BuddyQueue{}

	if err := jsonx.Unmarshal(data, v0BuddyQueue); err != nil {
		return nil, err
	} else {
		return v0BuddyQueue, nil
	}
}
