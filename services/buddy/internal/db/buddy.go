package db

import (
	"github.com/donglei1234/platform/services/common/nosql/document"
	errors2 "github.com/donglei1234/platform/services/common/nosql/errors"
	"regexp"

	"github.com/donglei1234/platform/services/buddy/internal/db/schemas/latest"
)

const (
	validatePathPattern = `[a-zA-Z0-9_.-]`
)

var (
	validatePathExp *regexp.Regexp
)

func init() {
	validatePathExp = regexp.MustCompile(validatePathPattern)
}

type BuddyQueue struct {
	document.Document
	latest.BuddyQueue
}

func (b *BuddyQueue) init(appId, id string, ros document.DocumentStore) error {
	if ros == nil {
		return errors2.ErrDocumentStoreIsNil
	}
	key, e := newBuddyQueueKey(appId, id)
	if e != nil {
		return e
	}
	b.BuddyQueue = *latest.NewBuddyQueue()
	b.Uid = id
	b.Document.Init(&b.BuddyQueue, b.clear, ros, key)
	return nil
}

func (b *BuddyQueue) clear() {
	b.BuddyQueue = latest.BuddyQueue{Buddies: make(map[string]*latest.Buddy, 0)}
}

func (b *BuddyQueue) InitDefault() error {
	return nil
}

func NewRelativeBuddyQueuePath(appId string, id string) (string, error) {
	if !validatePathExp.MatchString(id) {
		return "", document.ErrInvalidKeyFormat
	} else if !validatePathExp.MatchString(appId) {
		return "", document.ErrInvalidKeyFormat
	} else {
		return "/" + appId + "/buddies/" + id, nil
	}
}

func newBuddyQueueKey(appId string, id string) (document.Key, error) {
	if result, err := NewRelativeBuddyQueuePath(appId, id); err != nil {
		return document.Key{}, err
	} else {
		return document.NewKeyFromString(result)
	}
}
