package client

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/donglei1234/platform/services/common/nosql/document"
	leaderboard "github.com/donglei1234/platform/services/leaderboard/pkg/client"
	nfx "github.com/donglei1234/platform/services/leaderboard/pkg/fx"
)

type Leaderboard struct {
	fx.In
	LeaderboardClient leaderboard.LeaderboardClient `name:"LeaderboardClient"`
}

type LeaderboardClientFactory struct {
	fx.Out
	LeaderboardClient leaderboard.LeaderboardClient `name:"LeaderboardClient"`
}

// using this locally to make the Leaderboard/session store situationally optional
type LeaderboardClientParams struct {
	fx.In
	nfx.LeaderboardSettings
	LeaderboardStore document.DocumentStore `name:"LeaderboardStore" optional:"true"`
	SessionStore     document.DocumentStore `name:"SessionStore" optional:"true"`
}

func (g *LeaderboardClientFactory) Execute(
	l *zap.Logger,
	t fxsvcapp.GlobalTracer,
	s fxsvcapp.SecuritySettings,
	a LeaderboardClientParams,
) (err error) {
	l.Info("connect",
		zap.String("service", "leaderboard"),
		//zap.Bool("loopback", a.LoopbackAuth),
		zap.String("url", a.LeaderboardUrl),
	)
	g.LeaderboardClient, err = leaderboard.NewLeaderboardClient(l, a.LeaderboardUrl, s.SecureClients)

	return
}

var LeaderboardClientModule = fx.Provide(
	func(
		l *zap.Logger,
		t fxsvcapp.GlobalTracer,
		s fxsvcapp.SecuritySettings,
		a LeaderboardClientParams,
	) (out LeaderboardClientFactory, err error) {
		err = out.Execute(l, t, s, a)
		return
	},
)
