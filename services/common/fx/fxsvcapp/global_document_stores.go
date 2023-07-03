package fxsvcapp

import (
	"context"
	"fmt"
	"github.com/donglei1234/platform/services/common/nosql/document"
	"github.com/donglei1234/platform/services/common/nosql/document/mongodb"
	"net/url"
	"runtime"

	"go.uber.org/fx"
	"go.uber.org/zap"

	afx "github.com/donglei1234/platform/services/auth/pkg/fx"
	"github.com/donglei1234/platform/services/common/nosql/document/badger"
	"github.com/donglei1234/platform/services/common/nosql/document/couchbase"
	"github.com/donglei1234/platform/services/common/nosql/document/mock"
	nfx "github.com/donglei1234/platform/services/common/nosql/pkg/fx"
)

// GlobalDocumentStoreProvider loads the global DocumentStoreProvider from the fx dependency graph.
type GlobalDocumentStoreProvider struct {
	fx.In
	DocumentStoreProvider document.DocumentStoreProvider `name:"DocumentStoreProvider"`
}

// GlobalDocumentStoreFactory injects server instances into the fx dependency graph based on values in GlobalSettings.
type GlobalDocumentStoreFactory struct {
	fx.Out
	DocumentStoreProvider document.DocumentStoreProvider `name:"DocumentStoreProvider"`
}

func (g *GlobalDocumentStoreFactory) Execute(
	lc fx.Lifecycle,
	l *zap.Logger,
	s GlobalSettings,
	n nfx.NoSQLSettings,
) (err error) {
	document.SetNamespace(s.Deployment)

	if n.DocumentStoreUrl != "" {
		if u, e := url.Parse(n.DocumentStoreUrl); e != nil {
			err = e
		} else {
			switch u.Scheme {
			case "couchbase":
				username := u.User.Username()
				if username == "" {
					username = n.NoSqlUser
				}

				password, set := u.User.Password()
				if !set {
					password = n.NoSqlPassword
				}

				cfg := couchbase.ClusterConfig{
					ConnUrl:  fmt.Sprintf("couchbase://%s", u.Host),
					Username: username,
					Password: password,
				}
				LogConnectService(l, "couchbase", cfg.ConnUrl)
				g.DocumentStoreProvider, err = couchbase.NewDocumentStoreProvider(cfg, l)
			case "badger":
				LogCreateResource(l, "DocumentStoreProvider", "badger")
				path := u.Path
				if runtime.GOOS == "windows" {
					// u.Path always has leading /, which on Windows does not play well with drive letters.
					if len(path) > 2 && path[0] == '/' && path[2] == ':' {
						path = path[1:]
					}
				}
				g.DocumentStoreProvider, err = badger.NewDocumentStoreProvider(path, n.GCInterval, l)
			case "mongodb":
				username := u.User.Username()
				if username == "" {
					username = n.NoSqlUser
				}

				password, set := u.User.Password()
				if !set {
					password = n.NoSqlPassword
				}

				cfg := mongodb.ClusterConfig{
					ConnUrl:  fmt.Sprintf("mongodb://%s", u.Host),
					Username: username,
					Password: password,
				}
				LogConnectService(l, "mongodb", cfg.ConnUrl)
				g.DocumentStoreProvider, err = mongodb.NewDocumentStoreProvider(cfg, l)

			case "test":
				LogCreateResource(l, "DocumentStoreProvider", "test")
				g.DocumentStoreProvider, err = mock.NewDocumentStoreProvider()

			default:
				return ErrInvalidDocumentStoreURL
			}
		}
	} else {
		return ErrMissingDocumentStoreURL
	}

	if g.DocumentStoreProvider != nil {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return g.DocumentStoreProvider.Shutdown()
			},
		})
	}

	return
}

// GlobalSessionStore loads the global "sessions" DocumentStore from the fx dependency graph.
type GlobalSessionStore struct {
	fx.In
	SessionStore document.DocumentStore `name:"SessionStore"`
}

// GlobalSessionStoreFactory injects couchbase "sessions" store instances into the fx dependency graph based on
// values in GlobalSettings.
type GlobalSessionStoreFactory struct {
	fx.Out
	SessionStore document.DocumentStore `name:"SessionStore"`
}

func (g *GlobalSessionStoreFactory) Execute(
	l *zap.Logger,
	n nfx.NoSQLSettings,
	d GlobalDocumentStoreProvider,
) (err error) {
	g.SessionStore, err = openDocumentStore(
		l,
		d.DocumentStoreProvider,
		n.SessionStoreName,
		"sessions",
	)
	return
}

// GlobalAuthStore loads all of the global couchbase "auth" store instances from the fx dependency graph.
type GlobalAuthStore struct {
	fx.In
	AuthStore document.DocumentStore `name:"AuthStore"`
}

// GlobalSessionStoreFactory injects couchbase "auth" store instances into the fx dependency graph based on
// values in GlobalSettings.
type GlobalAuthStoreFactory struct {
	fx.Out
	AuthStore document.DocumentStore `name:"AuthStore"`
}

func (g *GlobalAuthStoreFactory) Execute(
	l *zap.Logger,
	a afx.AuthSettings,
	d GlobalDocumentStoreProvider,
) (err error) {
	g.AuthStore, err = openDocumentStore(
		l,
		d.DocumentStoreProvider,
		a.AuthStoreName,
		"auth",
	)
	return
}

// GlobalRosDataStore loads all of the global couchbase "ros data" store instances from the fx dependency graph.
type GlobalRosDataStore struct {
	fx.In
	RosDataStore document.DocumentStore `name:"RosDataStore"`
}

// GlobalRosDataStoreFactory injects couchbase "ros data" store instances into the fx dependency graph based on
// values in GlobalSettings.
type GlobalRosDataStoreFactory struct {
	fx.Out
	RosDataStore document.DocumentStore `name:"RosDataStore"`
}

//func (g *GlobalRosDataStoreFactory) Execute(
//	l *zap.Logger,
//	r rfx.RosSettings,
//	d GlobalDocumentStoreProvider,
//) (err error) {
//	g.RosDataStore, err = openDocumentStore(
//		l,
//		d.DocumentStoreProvider,
//		r.RosStoreName,
//		"ros",
//	)
//	return
//}

// GlobalGameServerStore loads all of the global couchbase "gameservers" store instances from the fx dependency graph.
type GlobalGameServerStore struct {
	fx.In
	GameServerStore document.DocumentStore `name:"GameServerStore"`
}

// GlobalGameServerStoreFactory injects couchbase "gameservers" store instances into the fx dependency graph based on
// values in GlobalSettings.
type GlobalGameServerStoreFactory struct {
	fx.Out
	GameServerStore document.DocumentStore `name:"GameServerStore"`
}

func (g *GlobalGameServerStoreFactory) Execute(
	l *zap.Logger,
	s GlobalSettings,
	//tata fx2.TaTaSettingsLoader,
	d GlobalDocumentStoreProvider,
) (err error) {
	g.GameServerStore, err = openDocumentStore(
		l,
		d.DocumentStoreProvider,
		"game",
		"game",
	)
	return
}

func openDocumentStore(
	logger *zap.Logger,
	provider document.DocumentStoreProvider,
	name string,
	tag string,
) (result document.DocumentStore, err error) {
	LogOpenDocumentStore(logger, name, tag)
	return provider.OpenDocumentStore(name)
}

var DocumentStoreProviderModule = fx.Provide(
	func(
		s GlobalSettings,
		lc fx.Lifecycle,
		l *zap.Logger,
		n nfx.NoSQLSettings,
	) (out GlobalDocumentStoreFactory, err error) {
		err = out.Execute(lc, l, s, n)
		return
	},
)

var SessionsStoreModule = fx.Provide(
	func(
		l *zap.Logger,
		n nfx.NoSQLSettings,
		d GlobalDocumentStoreProvider,
	) (out GlobalSessionStoreFactory, err error) {
		err = out.Execute(l, n, d)
		return
	},
)

var AuthStoreModule = fx.Provide(
	func(
		l *zap.Logger,
		a afx.AuthSettings,
		d GlobalDocumentStoreProvider,
	) (out GlobalAuthStoreFactory, err error) {
		err = out.Execute(l, a, d)
		return
	},
)

//var RosDataStoreModule = fx.Provide(
//	func(
//		l *zap.Logger,
//		r rfx.RosSettings,
//		d GlobalDocumentStoreProvider,
//	) (out GlobalRosDataStoreFactory, err error) {
//		err = out.Execute(l, r, d)
//		return
//	},
//)

var GameServerStoreModule = fx.Provide(
	func(
		l *zap.Logger,
		s GlobalSettings,
		//tata fx2.TaTaSettingsLoader,
		d GlobalDocumentStoreProvider,
	) (out GlobalGameServerStoreFactory, err error) {
		err = out.Execute(l, s, d)
		return
	},
)
