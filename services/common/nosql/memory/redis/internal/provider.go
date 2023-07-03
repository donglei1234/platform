package internal

import (
	"github.com/go-redis/redis"
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/common/nosql/memory"
)

type MemoryStoreProvider struct {
	client *redis.Client
	logger *zap.Logger
	store  *MemoryStore
}

func NewMemoryStoreProvider(addr string, password string, l *zap.Logger) (memory.MemoryStoreProvider, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	return &MemoryStoreProvider{client: c, logger: l}, nil
}

func (d *MemoryStoreProvider) OpenMemoryStore(name string) (mm memory.MemoryStore, err error) {
	d.store, err = NewMemoryStore(name, d.client, d.logger)
	mm = d.store
	return
}

func (d *MemoryStoreProvider) Shutdown() error {
	d.store.ShutDown()
	d.client.Close()
	return nil
}
