package fxsvcapp

import (
	"context"
	nfx "github.com/donglei1234/platform/services/common/nosql/pkg/fx"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/url"

	"github.com/donglei1234/platform/services/common/nosql/memory"
	"github.com/donglei1234/platform/services/common/nosql/memory/keys"
	"github.com/donglei1234/platform/services/common/nosql/memory/redis"
	goredis "github.com/go-redis/redis"
)

type MemoryStoreName string

const (
	MemoryStoreRoom = "room"
)

// GlobalMemoryStoreProvider loads the global MemoryStoreProvider from the fx dependency graph.
type GlobalMemoryStoreProvider struct {
	fx.In
	MemoryStoreProvider memory.MemoryStoreProvider `name:"MemoryStoreProvider"`
}

// GlobalMemoryStoreFactory injects server instances into the fx dependency graph based on values in GlobalSettings.
type GlobalMemoryStoreFactory struct {
	fx.Out
	MemoryStoreProvider memory.MemoryStoreProvider `name:"MemoryStoreProvider"`
}

func (g *GlobalMemoryStoreFactory) Execute(
	lc fx.Lifecycle,
	l *zap.Logger,
	s GlobalSettings,
	n nfx.NoSQLSettings,
) (err error) {
	keys.SetNamespace(s.Deployment)

	if n.MemoryStoreUrl != "" {
		if u, e := url.Parse(n.MemoryStoreUrl); e != nil {
			err = e
		} else {
			switch u.Scheme {
			case "redis":
				password, set := u.User.Password()
				if !set {
					password = n.NoSqlPassword
				}
				LogConnectService(l, "redis", u.Host)
				g.MemoryStoreProvider, err = redis.NewMemoryStoreProvider(u.Host, password, l)
			default:
				return ErrInvalidMemoryStoreURL
			}
		}
	} else {
		return ErrMissingMemoryStoreURL
	}

	if g.MemoryStoreProvider != nil {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return g.MemoryStoreProvider.Shutdown()
			},
		})
	}

	return
}

// GlobalRedisParams provides the GlobalRedisParams to the fx dependency graph.
type GlobalRedisParams struct {
	fx.In
	Redis *goredis.Client `name:"GlobalRedis"`
	Cache *goredis.Client `name:"GlobalCache"`
}

// GlobalRedisResult provides the GlobalRedisResult to the fx dependency graph.
type GlobalRedisResult struct {
	fx.Out
	Redis *goredis.Client `name:"GlobalRedis"`
	Cache *goredis.Client `name:"GlobalCache"`
}

func (g *GlobalRedisResult) Execute(
	lc fx.Lifecycle,
	l *zap.Logger,
	s GlobalSettings,
	n nfx.NoSQLSettings,
) (err error) {
	keys.SetNamespace(s.Deployment)
	if n.MemoryStoreUrl != "" {
		if u, e := url.Parse(n.MemoryStoreUrl); e != nil {
			err = e
		} else {
			switch u.Scheme {
			case "redis":
				password, set := u.User.Password()
				if !set {
					password = n.NoSqlPassword
				}
				LogConnectService(l, "redis", u.Host)
				g.Redis = goredis.NewClient(&goredis.Options{
					Addr:     u.Host,
					Password: password,
					DB:       0,
				})
				g.Cache = goredis.NewClient(&goredis.Options{
					Addr:     u.Host,
					Password: password,
					DB:       1,
				})
			default:
				return ErrInvalidMemoryStoreURL
			}
		}
	} else {
		return ErrMissingMemoryStoreURL
	}

	if g.Redis != nil {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				g.Redis.Close()
				return g.Cache.Close()
			},
		})
	}

	return
}

type GlobalRoomServerStore struct {
	fx.In
	RoomServerStore memory.MemoryStore `name:"RoomServerStore"`
}

type GlobalRoomServerStoreFactory struct {
	fx.Out
	RoomServerStore memory.MemoryStore `name:"RoomServerStore"`
}

func (g *GlobalRoomServerStoreFactory) Execute(
	l *zap.Logger,
	d GlobalMemoryStoreProvider,
) (err error) {
	g.RoomServerStore, err = openMemoryStore(
		l,
		d.MemoryStoreProvider,
		MemoryStoreRoom,
		MemoryStoreRoom,
	)
	return
}

var RoomStoreModule = fx.Provide(
	func(
		l *zap.Logger,
		n nfx.NoSQLSettings,
		d GlobalMemoryStoreProvider,
	) (out GlobalRoomServerStoreFactory, err error) {
		err = out.Execute(l, d)
		return
	},
)

type GlobalConditionServerStore struct {
	fx.In
	ConditionServerStore memory.MemoryStore `name:"ConditionServerStore"`
}

type GlobalConditionServerStoreFactory struct {
	fx.Out
	ConditionServerStore memory.MemoryStore `name:"ConditionServerStore"`
}

func (g *GlobalConditionServerStoreFactory) Execute(
	l *zap.Logger,
	d GlobalMemoryStoreProvider,
) (err error) {
	g.ConditionServerStore, err = openMemoryStore(
		l,
		d.MemoryStoreProvider,
		"condition",
		"condition",
	)
	return
}

var ConditionStoreModule = fx.Provide(
	func(
		l *zap.Logger,
		n nfx.NoSQLSettings,
		d GlobalMemoryStoreProvider,
	) (out GlobalConditionServerStoreFactory, err error) {
		err = out.Execute(l, d)
		return
	},
)

func openMemoryStore(
	logger *zap.Logger,
	provider memory.MemoryStoreProvider,
	name string,
	tag string,
) (result memory.MemoryStore, err error) {
	LogOpenMemoryStore(logger, name, tag)
	return provider.OpenMemoryStore(name)
}

var MemoryStoreProviderModule = fx.Provide(
	func(
		s GlobalSettings,
		lc fx.Lifecycle,
		l *zap.Logger,
		n nfx.NoSQLSettings,
	) (out GlobalMemoryStoreFactory, err error) {
		err = out.Execute(lc, l, s, n)
		return
	},
)

var RedisStoreProviderModule = fx.Provide(
	func(
		s GlobalSettings,
		lc fx.Lifecycle,
		l *zap.Logger,
		n nfx.NoSQLSettings,
	) (out GlobalRedisResult, err error) {
		err = out.Execute(lc, l, s, n)
		return
	},
)
