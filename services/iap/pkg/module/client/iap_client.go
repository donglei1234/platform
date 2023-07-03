package client

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/donglei1234/platform/services/common/nosql/document"
	iap "github.com/donglei1234/platform/services/iap/pkg/client"
	nfx "github.com/donglei1234/platform/services/iap/pkg/fx"
)

type IAPClient struct {
	fx.In
	IAPClient iap.IAPClient `name:"IAPClient"`
}

type IAPClientFactory struct {
	fx.Out
	IAPClient iap.IAPClient `name:"IAPClient"`
}

// using this locally to make the IAP/session store situationally optional
type IAPClientParams struct {
	fx.In
	nfx.IAPSettings
	IAPStore     document.DocumentStore `name:"IAPStore" optional:"true"`
	SessionStore document.DocumentStore `name:"SessionStore" optional:"true"`
}

func (g *IAPClientFactory) Execute(
	l *zap.Logger,
	t fxsvcapp.GlobalTracer,
	s fxsvcapp.SecuritySettings,
	a IAPClientParams,
) (err error) {
	l.Info("connect",
		zap.String("service", "iapClient"),
		zap.String("url", a.IAPUrl),
	)
	g.IAPClient, err = iap.NewIAPClient(l, a.IAPUrl, s.SecureClients)

	return
}

var IAPClientModule = fx.Provide(
	func(
		l *zap.Logger,
		t fxsvcapp.GlobalTracer,
		s fxsvcapp.SecuritySettings,
		a IAPClientParams,
	) (out IAPClientFactory, err error) {
		err = out.Execute(l, t, s, a)
		return
	},
)
