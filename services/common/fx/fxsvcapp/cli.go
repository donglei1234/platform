package fxsvcapp

import (
	"github.com/donglei1234/platform/services/common/config"
	"github.com/donglei1234/platform/services/common/fx/fxapp"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type cliParams struct {
	fx.In
	config.AllHelp
	Runners       GlobalRunners
	ServiceBinder ServiceBinder

	Lifecycle  fx.Lifecycle
	Controller *fxapp.Controller
	Logger     *zap.Logger

	AppName string `name:"AppName"`
}

func cli(in cliParams) (err error) {
	exitApp := false

	if err := in.ServiceBinder.Execute(in.Logger, in.Lifecycle); err != nil {
		exitApp = true
	}
	if err := in.Runners.Execute(in.Lifecycle); err != nil {
		exitApp = true
	}

	if exitApp {
		in.Controller.Quit()
	}

	return err
}
