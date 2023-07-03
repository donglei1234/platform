package fx

import (
	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/donglei1234/platform/services/common/nosql/document"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// GlobalGameServerStore loads all of the global couchbase "gameservers" store instances from the fx dependency graph.
type StoreParams struct {
	fx.In
	ServerStore document.DocumentStore `name:"ChatServerStore"`
}

// GlobalGameServerStoreFactory injects couchbase "gameservers" store instances into the fx dependency graph based on
// values in GlobalSettings.
type StoreResult struct {
	fx.Out
	ServerStore document.DocumentStore `name:"ChatServerStore"`
}

func (g *StoreResult) Execute(
	l *zap.Logger,
	s ChatSettings,
	d fxsvcapp.GlobalDocumentStoreProvider,
) (err error) {
	g.ServerStore, err = openDocumentStore(
		l,
		d.DocumentStoreProvider,
		s.Name,
		s.Name,
	)
	return
}

var StoreProviderModule = fx.Provide(
	func(
		l *zap.Logger,
		s ChatSettings,
		d fxsvcapp.GlobalDocumentStoreProvider,
	) (out StoreResult, err error) {
		err = out.Execute(l, s, d)
		return
	},
)

func openDocumentStore(
	logger *zap.Logger,
	provider document.DocumentStoreProvider,
	name string,
	tag string,
) (result document.DocumentStore, err error) {
	logger.Info("open",
		zap.String("kind", "nosql.DocumentStore"),
		zap.String("name", name),
		zap.String("tag", tag),
	)
	return provider.OpenDocumentStore(name)
}
