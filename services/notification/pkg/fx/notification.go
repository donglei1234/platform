package fx

import (
	"go.uber.org/fx"

	"github.com/donglei1234/platform/services/common/config"
)

type NotificationSettings struct {
	fx.In
	NotificationStoreName string `name:"NotificationStoreName"`
	NotificationUrl       string `name:"NotificationUrl"`
}

type NotificationSettingsLoader struct {
	fx.Out
	config.EnvironmentBlock
	NotificationStoreName string `name:"NotificationStoreName" envconfig:"NOTIFICATION_STORE_NAME" default:"Notification"`
	NotificationUrl       string `name:"NotificationUrl" envconfig:"NOTIFICATION_URL" default:"localhost:8081"`
}

func (g *NotificationSettingsLoader) LoadFromEnv() (err error) {
	err = config.Load(g)
	return
}

var NotificationSettingsModule = fx.Provide(
	func() (out NotificationSettingsLoader, err error) {
		err = out.LoadFromEnv()
		return
	},
)
