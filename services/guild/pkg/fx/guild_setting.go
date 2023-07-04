package fx

import (
	"go.uber.org/fx"

	"github.com/donglei1234/platform/services/common/config"
)

type GuildSettings struct {
	fx.In
	GuildStoreName string `name:"GuildStoreName"`
	GuildUrl       string `name:"GuildUrl"`
}

type GuildSettingsLoader struct {
	fx.Out
	config.EnvironmentBlock
	GuildStoreName string `name:"GuildStoreName" envconfig:"UNION_STORE_NAME" default:"Guild"`
	GuildUrl       string `name:"GuildUrl" envconfig:"UNION_URL" default:"localhost:8081"`
}

func (g *GuildSettingsLoader) LoadFromEnv() (err error) {
	err = config.Load(g)
	return
}

var GuildSettingsModule = fx.Provide(
	func() (out GuildSettingsLoader, err error) {
		err = out.LoadFromEnv()
		return
	},
)
