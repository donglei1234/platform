package fxsvcapp

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	auth "github.com/donglei1234/platform/services/auth/pkg/client"
	afx "github.com/donglei1234/platform/services/auth/pkg/fx"
	"github.com/donglei1234/platform/services/common/nosql/document"
)

type GlobalAuthClient struct {
	fx.In
	AuthClient auth.AuthClient `name:"AuthClient"`
}

type GlobalAuthClientFactory struct {
	fx.Out
	AuthClient auth.AuthClient `name:"AuthClient"`
}

// using this locally to make the auth/session store situationally optional
type GlobalAuthClientParams struct {
	fx.In
	afx.AuthSettings
	AuthStore    document.DocumentStore `name:"AuthStore" optional:"true"`
	SessionStore document.DocumentStore `name:"SessionStore" optional:"true"`
}

func (g *GlobalAuthClientFactory) Execute(
	l *zap.Logger,
	t GlobalTracer,
	s SecuritySettings,
	a GlobalAuthClientParams,
) (err error) {
	//if a.LoopbackAuth {
	//	l.Info("connect",
	//		zap.String("service", "authclient"),
	//		zap.Bool("loopback", a.LoopbackAuth),
	//		zap.Bool("authRequired", a.AuthRequired),
	//		zap.String("authStore.Name", a.AuthStore.Name()),
	//		zap.String("sessionStore.Name", a.SessionStore.Name()),
	//	)
	//	// NB: auth/session store may rightfully be omitted in cases where auth is not required
	//	if a.AuthRequired && (a.AuthStore.Name() == "" || a.SessionStore.Name() == "") {
	//		return nosql.ErrDocumentStoreIsNil
	//	}
	//	g.AuthClient, err = auth.NewLoopbackClient(l, a.AuthRequired, a.AuthStore, a.SessionStore)
	//} else {
	l.Info("connect",
		zap.String("service", "authclient"),
		//zap.Bool("loopback", a.LoopbackAuth),
		zap.String("url", a.AuthUrl),
	)
	g.AuthClient, err = auth.NewAuthClient(l, a.AuthUrl, s.SecureClients)
	//}

	return
}

var AuthClientModule = fx.Provide(
	func(
		l *zap.Logger,
		t GlobalTracer,
		s SecuritySettings,
		a GlobalAuthClientParams,
	) (out GlobalAuthClientFactory, err error) {
		err = out.Execute(l, t, s, a)
		return
	},
)
