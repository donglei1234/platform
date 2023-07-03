package db

import (
	"github.com/donglei1234/platform/services/common/nosql/document"
	errors2 "github.com/donglei1234/platform/services/common/nosql/errors"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Database struct {
	logger *zap.Logger
	ros    document.DocumentStore
}

func OpenDatabase(l *zap.Logger, ros document.DocumentStore) Database {
	return Database{
		logger: l,
		ros:    ros,
	}
}

func (d *Database) ValidateSettings() error {
	if d.ros == nil {
		return errors2.ErrDocumentStoreIsNil
	} else {
		return nil
	}
}

func (d *Database) NewBuddyQueue(appId string, id string) (*BuddyQueue, error) {
	bq := new(BuddyQueue)
	err := bq.init(appId, id, d.ros)
	if err != nil {
		return nil, err
	}
	return bq, nil
}

func (d *Database) ContainsBuddyQueue(appId string, id string) (bool, error) {
	if key, e := newBuddyQueueKey(appId, id); e != nil {
		return false, e
	} else {
		return d.ros.Contains(key)
	}
}

func (d *Database) IncrBuddyReward(appId string, id, path string, value int32) error {
	if key, e := newBuddyQueueKey(appId, id); e != nil {
		return e
	} else {
		return d.ros.Incr(key, path, value)
	}
}

func (d *Database) PushBack(appId string, id, path string, profileId string) error {
	if key, e := newBuddyQueueKey(appId, id); e != nil {
		return e
	} else {
		return d.ros.PushBack(key, path, profileId)
	}
}

func (d *Database) CreateBuddyQueue(appId, id string) error {
	if bq, err := d.NewBuddyQueue(appId, id); err != nil {
		return err
	} else if err = bq.Create(); err != nil {
		return err
	}
	return nil
}

func (d *Database) LoadOrCreateBuddyQueue(appId string, id string) (bq *BuddyQueue, err error) {
	if bq, err = d.NewBuddyQueue(appId, id); err != nil {
		return
	} else if err = bq.Load(); errors.Cause(err) == errors2.ErrKeyNotFound {
		if bq, err = d.NewBuddyQueue(appId, id); err != nil {
			return
		} else if err := bq.InitDefault(); err != nil {
			return nil, err
		} else if err = bq.Create(); err != nil {
			err = bq.Load()
		}
	}
	// Even if we know the BuddyQueue is on the latest version, we should still run it through fixups
	if err == nil {
		err = d.FixupBuddyQueue(bq)
	}
	return
}
func (d *Database) LoadBuddyQueue(appId string, id string) (bq *BuddyQueue, err error) {
	if bq, err = d.NewBuddyQueue(appId, id); err != nil {
		return
	} else if err = bq.Load(); err != nil {
		return
	}
	// Even if we know the BuddyQueue is on the latest version, we should still run it through fixups
	if err == nil {
		err = d.FixupBuddyQueue(bq)
	}
	return
}
