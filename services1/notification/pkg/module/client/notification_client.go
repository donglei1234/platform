package client

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/donglei1234/platform/services/common/nosql/document"
	notification "github.com/donglei1234/platform/services/notification/pkg/client"
	nfx "github.com/donglei1234/platform/services/notification/pkg/fx"
)

type NotificationClient struct {
	fx.In
	NotificationClient notification.NotificationClient `name:"NotificationClient"`
}

type NotificationClientFactory struct {
	fx.Out
	NotificationClient notification.NotificationClient `name:"NotificationClient"`
}

// using this locally to make the Notification/session store situationally optional
type NotificationClientParams struct {
	fx.In
	nfx.NotificationSettings
	NotificationStore document.DocumentStore `name:"NotificationStore" optional:"true"`
	SessionStore      document.DocumentStore `name:"SessionStore" optional:"true"`
}

func (g *NotificationClientFactory) Execute(
	l *zap.Logger,
	t fxsvcapp.GlobalTracer,
	s fxsvcapp.SecuritySettings,
	a NotificationClientParams,
) (err error) {
	l.Info("connect",
		zap.String("service", "notificationClient"),
		//zap.Bool("loopback", a.LoopbackAuth),
		zap.String("url", a.NotificationUrl),
	)
	g.NotificationClient, err = notification.NewNotificationClient(l, a.NotificationUrl, s.SecureClients)

	return
}

var NotificationClientModule = fx.Provide(
	func(
		l *zap.Logger,
		t fxsvcapp.GlobalTracer,
		s fxsvcapp.SecuritySettings,
		a NotificationClientParams,
	) (out NotificationClientFactory, err error) {
		err = out.Execute(l, t, s, a)
		return
	},
)
