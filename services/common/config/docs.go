package config

import "go.uber.org/fx"

// Help contains usage information for an environment configuration block.
type Help struct {
	Name  string
	Usage string
}

// AllHelp collects all config.Help instances that are in the current app container.
type AllHelp struct {
	fx.In
	Values []Help `group:"config.Help"`
}

// HelpOut provides a Help instance using Fx dependency injection.
type HelpOut struct {
	fx.Out
	Help Help `group:"config.Help" ignored:"true"`
}
