package fx

import (
	"go.uber.org/fx"

	"github.com/donglei1234/platform/services/common/config"
)

type MailSettings struct {
	fx.In
	MailStoreName string `name:"MailStoreName"`
	MailUrl       string `name:"MailUrl"`
}

type MailSettingsLoader struct {
	fx.Out
	config.EnvironmentBlock
	MailStoreName string `name:"MailStoreName" envconfig:"ROOM_STORE_NAME" default:"mail"`
	MailUrl       string `name:"MailUrl" envconfig:"MAIL_URL" default:"localhost:8081"`
}

func (g *MailSettingsLoader) LoadFromEnv() (err error) {
	err = config.Load(g)
	return
}

var MailSettingsModule = fx.Provide(
	func() (out MailSettingsLoader, err error) {
		err = out.LoadFromEnv()
		return
	},
)
