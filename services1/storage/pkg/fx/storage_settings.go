package fx

import (
	"go.uber.org/fx"

	"github.com/donglei1234/platform/services/common/config"
)

type StorageSettings struct {
	fx.In
	StorageStoreName string `name:"StorageStoreName"`
	StorageUrl       string `name:"StorageUrl"`
}

type StorageSettingsLoader struct {
	fx.Out
	config.EnvironmentBlock
	StorageStoreName string `name:"StorageStoreName" envconfig:"STORAGE_STORE_NAME" default:"storage"`
	StorageUrl       string `name:"StorageUrl" envconfig:"STORAGE_URL" default:"localhost:8081"`
}

func (g *StorageSettingsLoader) LoadFromEnv() (err error) {
	err = config.Load(g)
	return
}

var StorageSettingsModule = fx.Provide(
	func() (out StorageSettingsLoader, err error) {
		err = out.LoadFromEnv()
		return
	},
)
