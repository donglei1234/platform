package fxsvcapp

import (
	"go.uber.org/fx"

	"github.com/donglei1234/platform/services/common/config"
)

type NetworkSettings struct {
	fx.In

	Port   int    `name:"Port"`
	Socket string `name:"Socket"`
}

type NetworkSettingsLoader struct {
	fx.Out
	config.EnvironmentBlock

	Port   int    `name:"Port" default:"8081" envconfig:"PORT"`
	Socket string `name:"Socket" envconfig:"SOCKET"`
}

func (g *NetworkSettingsLoader) LoadFromEnv() (err error) {
	err = config.Load(g)
	return
}

var NetworkSettingsModule = fx.Provide(
	func() (out NetworkSettingsLoader, err error) {
		err = out.LoadFromEnv()
		return
	},
)
