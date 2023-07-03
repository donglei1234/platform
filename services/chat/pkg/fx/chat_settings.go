package fx

import (
	"github.com/donglei1234/platform/services/common/config"
	"go.uber.org/fx"
)

type ChatSettings struct {
	fx.In
	Name         string `name:"Name"`
	ChatInterval int    `name:"ChatInterval"`
}

type ChatSettingsLoader struct {
	fx.Out
	config.EnvironmentBlock
	Name         string `name:"Name" envconfig:"NAME" default:"chat"`
	ChatInterval int    `name:"ChatInterval" envconfig:"WORLD_CHAT_INTERVAL" default:"5"`
}

func (l *ChatSettingsLoader) LoadFromEnv() (err error) {
	err = config.Load(l)
	return
}

var ChatSettingsModule = fx.Provide(
	func() (out ChatSettingsLoader, err error) {
		err = out.LoadFromEnv()
		return
	},
)
