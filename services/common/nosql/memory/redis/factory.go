package redis

import (
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/common/nosql/memory"
	"github.com/donglei1234/platform/services/common/nosql/memory/redis/internal"
)

type MemoryStoreProvider = internal.MemoryStoreProvider

func NewMemoryStoreProvider(addr string, password string, l *zap.Logger) (memory.MemoryStoreProvider, error) {
	return internal.NewMemoryStoreProvider(addr, password, l)
}
