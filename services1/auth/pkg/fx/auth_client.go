package fx

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	auth "github.com/donglei1234/platform/services/auth/pkg/client"
)

type GlobalAuthClient struct {
	fx.In
	AuthClient auth.AuthClient `name:"AuthClient"`
}

type GlobalAuthClientFactory struct {
	fx.Out
	AuthClient auth.AuthClient `name:"AuthClient"`
}

func (g *GlobalAuthClientFactory) Execute(
	l *zap.Logger,
	a AuthSettings,
) (err error) {
	l.Info("connect",
		zap.String("service", "authclient"),
		zap.String("url", a.AuthUrl),
	)
	g.AuthClient, err = auth.NewAuthClient(l, a.AuthUrl, false)
	return
}

var ClientModule = fx.Provide(
	func(
		l *zap.Logger,
		a AuthSettings,
	) (out GlobalAuthClientFactory, err error) {
		err = out.Execute(l, a)
		return
	},
)
