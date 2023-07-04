package client

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/donglei1234/platform/services/common/nosql/document"
	storage "github.com/donglei1234/platform/services/storage/pkg/client"
	nfx "github.com/donglei1234/platform/services/storage/pkg/fx"
)

type StorageClient struct {
	fx.In
	StorageClient storage.StorageClient `name:"StorageClient"`
}

type StorageClientFactory struct {
	fx.Out
	StorageClient storage.StorageClient `name:"StorageClient"`
}

// using this locally to make the storage/session store situational optional
type StorageClientParams struct {
	fx.In
	nfx.StorageSettings
	StorageStore document.DocumentStore `name:"StorageStore" optional:"true"`
	SessionStore document.DocumentStore `name:"SessionStore" optional:"true"`
}

func (g *StorageClientFactory) Execute(
	l *zap.Logger,
	t fxsvcapp.GlobalTracer,
	s fxsvcapp.SecuritySettings,
	a StorageClientParams,
) (err error) {
	l.Info("connect",
		zap.String("service", "StorageClient"),
		zap.String("url", a.StorageUrl),
	)
	g.StorageClient, err = storage.NewStorageClient(l, a.StorageUrl, s.SecureClients)

	return
}

var StorageClientModule = fx.Provide(
	func(
		l *zap.Logger,
		t fxsvcapp.GlobalTracer,
		s fxsvcapp.SecuritySettings,
		a StorageClientParams,
	) (out StorageClientFactory, err error) {
		err = out.Execute(l, t, s, a)
		return
	},
)
