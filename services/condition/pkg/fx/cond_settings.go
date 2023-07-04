package fx

import (
	"go.uber.org/fx"

	"github.com/donglei1234/platform/services/common/config"
)

type ConditionSettings struct {
	fx.In
	ConditionStoreName string `name:"ConditionStoreName"`
	ConditionUrl       string `name:"ConditionUrl"`
}

type ConditionSettingsLoader struct {
	fx.Out
	config.EnvironmentBlock
	ConditionStoreName string `name:"ConditionStoreName" envconfig:"CONDITION_STORE_NAME" default:"condition"`
	ConditionUrl       string `name:"ConditionUrl" envconfig:"CONDITION_URL" default:"localhost:8081"`
}

func (g *ConditionSettingsLoader) LoadFromEnv() (err error) {
	err = config.Load(g)
	return
}

var ConditionSettingsModule = fx.Provide(
	func() (out ConditionSettingsLoader, err error) {
		err = out.LoadFromEnv()
		return
	},
)
