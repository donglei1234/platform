package fx

import (
	"go.uber.org/fx"

	"github.com/donglei1234/platform/services/common/config"
)

type IAPSettings struct {
	fx.In
	IAPStoreName string `name:"IAPStoreName"`
	IAPUrl       string `name:"IAPUrl"`
}

type IAPSettingsLoader struct {
	fx.Out
	config.EnvironmentBlock
	IAPStoreName string `name:"IAPStoreName" envconfig:"IAP_STORE_NAME" default:"IAP"`
	IAPUrl       string `name:"IAPUrl" envconfig:"IAP_URL" default:"localhost:8081"`
}

func (g *IAPSettingsLoader) LoadFromEnv() (err error) {
	err = config.Load(g)
	return
}

var IAPSettingsModule = fx.Provide(
	func() (out IAPSettingsLoader, err error) {
		err = out.LoadFromEnv()
		return
	},
)
