package fx

import (
	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/go-redis/redis"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type BuddyMemoryParams struct {
	fx.In
	MemoryStore *redis.Client `name:"BuddyRedis"`
}

type BuddyMemoryResult struct {
	fx.Out
	MemoryStore *redis.Client `name:"BuddyRedis"`
}

func (g *BuddyMemoryResult) Execute(
	l *zap.Logger,
	settings BuddySettings,
	d fxsvcapp.GlobalRedisParams,
) (err error) {
	g.MemoryStore, err = openMemoryStore(
		l,
		d.Redis,
		settings.Name,
		settings.Name,
	)
	return
}

var BuddyMemoryModule = fx.Provide(
	func(
		l *zap.Logger,
		settings BuddySettings,
		d fxsvcapp.GlobalRedisParams,
	) (out BuddyMemoryResult, err error) {
		err = out.Execute(l, settings, d)
		return
	},
)

func openMemoryStore(
	logger *zap.Logger,
	provider *redis.Client,
	name string,
	tag string,
) (result *redis.Client, err error) {
	logger.Info("open",
		zap.String("kind", "nosql.MemoryStore"),
		zap.String("name", name),
		zap.String("tag", tag),
	)
	return provider, nil
}
