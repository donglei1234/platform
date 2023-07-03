package fx

import (
	"github.com/donglei1234/platform/services/common/config"
	"go.uber.org/fx"
)

type AuthSettings struct {
	fx.In

	//LoopbackAuth  bool   `name:"LoopbackAuth"`
	//AuthRequired  bool   `name:"AuthRequired"`
	AuthUrl       string `name:"AuthUrl"`
	AuthStoreName string `name:"AuthStoreName"`

	EphemeralToken string `name:"EphemeralToken"`
}

type AuthSettingsLoader struct {
	fx.Out
	config.EnvironmentBlock

	//LoopbackAuth  bool   `name:"LoopbackAuth" ignored:"true"`
	//AuthRequired  bool   `name:"AuthRequired" ignored:"true"`
	AuthStoreName  string `name:"AuthStoreName" envconfig:"AUTH_STORE_NAME" default:"auth"`
	AuthUrl        string `name:"AuthUrl" envconfig:"AUTH_URL" default:"localhost:8081"`
	EphemeralToken string `name:"EphemeralToken" default:"" envconfig:"EPHEMERAL_TOKEN"`
}

func (g *AuthSettingsLoader) LoadFromEnv() (err error) {
	err = config.Load(g)
	return
}

var SettingsModule = fx.Provide(
	func() (out AuthSettingsLoader, err error) {
		err = out.LoadFromEnv()
		return
	},
)
