package fx

import (
	"go.uber.org/fx"

	"github.com/donglei1234/platform/services/common/config"
)

type GMSettings struct {
	fx.In
	GmStoreName string `name:"GmStoreName"`
	GmUrl       string `name:"GmUrl"`
	GmMaxCount  int    `name:"GmMaxCount"`
}

type GmSettingsLoader struct {
	fx.Out
	config.EnvironmentBlock
	GmStoreName string `name:"GmStoreName" envconfig:"ROOM_STORE_NAME" default:"gm"`
	GmUrl       string `name:"GmUrl" envconfig:"ROOM_URL" default:"localhost:8081"`
	GmMaxCount  int    `name:"GmMaxCount" envconfig:"ROOM_MAX_COUNT" default:"200"`
}

func (g *GmSettingsLoader) LoadFromEnv() (err error) {
	err = config.Load(g)
	return
}

var GmSettingsModule = fx.Provide(
	func() (out GmSettingsLoader, err error) {
		err = out.LoadFromEnv()
		return
	},
)
