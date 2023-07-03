package fx

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	gm "github.com/donglei1234/platform/services/gm/pkg/client"
)

type GmClientParams struct {
	fx.In
	GmClient gm.GmClient `name:"GmClient"`
}

type GmClientResult struct {
	fx.Out
	GmClient gm.GmClient `name:"GmClient"`
}

func (g *GmClientResult) Execute(
	l *zap.Logger,
	t fxsvcapp.GlobalTracer,
	s fxsvcapp.SecuritySettings,
	a GMSettings,
) (err error) {
	l.Info("connect",
		zap.String("service", "gmclient"),
		zap.String("url", a.GmUrl),
	)
	g.GmClient, err = gm.NewGMClient(a.GmUrl, s.SecureClients)

	return
}

var GmClientModule = fx.Provide(
	func(
		l *zap.Logger,
		t fxsvcapp.GlobalTracer,
		s fxsvcapp.SecuritySettings,
		a GMSettings,
	) (out GmClientResult, err error) {
		err = out.Execute(l, t, s, a)
		return
	},
)
