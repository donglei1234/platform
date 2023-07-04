package fxsvcapp

import (
	"go.uber.org/fx"

	"github.com/donglei1234/platform/services/common/config"
)

type SecuritySettings struct {
	fx.In

	TlsCert string `name:"TlsCert"`
	TlsKey  string `name:"TlsKey"`

	SecureClients bool `name:"SecureClients"`
}

type SecuritySettingsLoader struct {
	fx.Out
	config.EnvironmentBlock

	TlsCert string `name:"TlsCert" envconfig:"TLSCERT"`
	TlsKey  string `name:"TlsKey" envconfig:"TLSKEY"`

	SecureClients bool `name:"SecureClients" envonfig:"SECURE_CLIENTS"`
}

func (g *SecuritySettingsLoader) LoadFromEnv() (err error) {
	err = config.Load(g)
	return
}

var SecuritySettingsModule = fx.Provide(
	func() (out SecuritySettingsLoader, err error) {
		err = out.LoadFromEnv()
		return
	},
)
