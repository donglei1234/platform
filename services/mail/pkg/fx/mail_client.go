package fx

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	mail "github.com/donglei1234/platform/services/mail/pkg/client"
)

type MailClientParams struct {
	fx.In
	MailClient mail.MailClient `name:"MailClient"`
}

type MailClientResult struct {
	fx.Out
	MailClient mail.MailClient `name:"MailClient"`
}

func (g *MailClientResult) Execute(
	l *zap.Logger,
	t fxsvcapp.GlobalTracer,
	s fxsvcapp.SecuritySettings,
	a MailSettings,
) (err error) {
	l.Info("connect",
		zap.String("service", "mailclient"),
		zap.String("url", a.MailUrl),
	)
	g.MailClient, err = mail.NewMailClient(a.MailUrl, s.SecureClients)

	return
}

var MailClientModule = fx.Provide(
	func(
		l *zap.Logger,
		t fxsvcapp.GlobalTracer,
		s fxsvcapp.SecuritySettings,
		a MailSettings,
	) (out MailClientResult, err error) {
		err = out.Execute(l, t, s, a)
		return
	},
)
