package client

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/donglei1234/platform/services/common/nosql/document"
	guild "github.com/donglei1234/platform/services/guild/pkg/client"
	nfx "github.com/donglei1234/platform/services/guild/pkg/fx"
)

type GuildClient struct {
	fx.In
	GuildClient guild.GuildClient `name:"GuildClient"`
}

type GuildClientFactory struct {
	fx.Out
	GuildClient guild.GuildClient `name:"GuildClient"`
}

// using this locally to make the Guild/session store situationally optional
type GuildClientParams struct {
	fx.In
	nfx.GuildSettings
	GuildStore   document.DocumentStore `name:"GuildStore" optional:"true"`
	SessionStore document.DocumentStore `name:"SessionStore" optional:"true"`
}

func (g *GuildClientFactory) Execute(
	l *zap.Logger,
	t fxsvcapp.GlobalTracer,
	s fxsvcapp.SecuritySettings,
	a GuildClientParams,
) (err error) {
	l.Info("connect",
		zap.String("service", "guildClient"),
		//zap.Bool("loopback", a.LoopbackAuth),
		zap.String("url", a.GuildUrl),
	)
	g.GuildClient, err = guild.NewGuildClient(l, a.GuildUrl, s.SecureClients)

	return
}

var GuildClientModule = fx.Provide(
	func(
		l *zap.Logger,
		t fxsvcapp.GlobalTracer,
		s fxsvcapp.SecuritySettings,
		a GuildClientParams,
	) (out GuildClientFactory, err error) {
		err = out.Execute(l, t, s, a)
		return
	},
)
