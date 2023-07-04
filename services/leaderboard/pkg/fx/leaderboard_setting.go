package fx

import (
	"go.uber.org/fx"

	"github.com/donglei1234/platform/services/common/config"
)

type LeaderboardSettings struct {
	fx.In
	LeaderboardStoreName string `name:"LeaderboardStoreName"`
	LeaderboardUrl       string `name:"LeaderboardUrl"`
}

type LeaderboardSettingsLoader struct {
	fx.Out
	config.EnvironmentBlock
	LeaderboardStoreName string `name:"LeaderboardStoreName" envconfig:"LEADERBOARD_STORE_NAME" default:"Leaderboard"`
	LeaderboardUrl       string `name:"LeaderboardUrl" envconfig:"LEADERBOARD_URL" default:"localhost:8081"`
}

func (g *LeaderboardSettingsLoader) LoadFromEnv() (err error) {
	err = config.Load(g)
	return
}

var LeaderboardSettingsModule = fx.Provide(
	func() (out LeaderboardSettingsLoader, err error) {
		err = out.LoadFromEnv()
		return
	},
)
