package fx

import (
	"github.com/donglei1234/platform/services/buddy/pkg/client"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type BuddyClientParams struct {
	fx.In
	BuddySettings
	BuddyClient client.PublicClient `name:"BuddyClient"`
}

type BuddyClientResult struct {
	fx.Out
	BuddyClient client.PublicClient `name:"BuddyClient"`
}

func (g *BuddyClientResult) Execute(
	l *zap.Logger,
	a BuddySettings,
) (err error) {
	l.Info("connect",
		zap.String("service", "condition client"),
		zap.String("url", a.BuddyUrl),
	)
	g.BuddyClient, err = client.NewPublicClient(a.BuddyUrl, false)
	return
}

var BuddyClientModule = fx.Provide(
	func(
		l *zap.Logger,
		a BuddySettings,
	) (out BuddyClientResult, err error) {
		err = out.Execute(l, a)
		return
	},
)
