package memory

import (
	"github.com/donglei1234/platform/services/common/nosql/memory/keys"
	"github.com/donglei1234/platform/services/common/nosql/memory/options"
)

// MemoryStoreProvider knows how to open document stores by name.
type MemoryStoreProvider interface {
	OpenMemoryStore(name string) (MemoryStore, error)
	Shutdown() error
}

type MemoryStore interface {
	// Name returns the name of this MemoryStore.
	Name() string

	Set(key keys.Key, opts ...options.Option) error

	MSet(opts ...options.Option) error

	SetNX(key keys.Key, opts ...options.Option) (bool, error)

	Incur(key keys.Key) (int64, error)

	Get(key keys.Key, opts ...options.Option) error

	Del(key keys.Key) error

	Exists(key keys.Key) (bool, error)

	HSet(path keys.Key, key keys.Key, opts ...options.Option) error

	HGet(path keys.Key, key keys.Key, opts ...options.Option) error

	HDel(path keys.Key, keys ...keys.Key) error

	HLen(path keys.Key) (int64, error)

	HMSet(path keys.Key, opts ...options.Option) error

	HMGet(path keys.Key, keys []keys.Key, opts ...options.Option) error

	HGetAll(path keys.Key, opts ...options.Option) error

	ZAdd(path keys.Key, opts ...options.Option) error

	ZScore(path keys.Key, key keys.Key) (float64, error)

	ZRevRank(path keys.Key, key keys.Key) (int64, error)

	ZRank(path keys.Key, key keys.Key) (int64, error)

	ZRevRange(path keys.Key, destOpt options.Option, opts ...options.ZRangeOption) error

	LPush(key keys.Key, opts ...options.Option) error

	RPush(key keys.Key, opts ...options.Option) error

	LRange(key keys.Key, ro options.ZRangeOption, opts ...options.Option) error

	LSet(key keys.Key, index int64, opts ...options.Option) error

	LPop(key keys.Key, opts ...options.Option) error

	LLen(key keys.Key) (int64, error)

	Scan(path keys.Key, dest options.Option, scanOpts ...options.ScanOption) (int, error)

	Remove(keys ...keys.Key) error

	Sub(key keys.Key) (chan interface{}, error)
}
