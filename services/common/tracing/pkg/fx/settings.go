package fx

import (
	"go.uber.org/fx"

	"github.com/donglei1234/platform/services/common/config"
)

type TracerSettings struct {
	fx.In

	TraceProvider    string   `name:"TraceProvider"`
	TraceAgentHost   string   `name:"TraceAgentHost"`
	TraceAgentPort   int      `name:"TraceAgentPort"`
	TraceServiceName string   `name:"TraceServiceName"`
	TraceTags        []string `name:"TraceTags"`
}

type TracerSettingsLoader struct {
	fx.Out
	config.EnvironmentBlock

	TraceProvider    string   `name:"TraceProvider" envconfig:"TRACE_PROVIDER"`
	TraceAgentHost   string   `name:"TraceAgentHost" envconfig:"TRACE_AGENT_HOST"`
	TraceAgentPort   int      `name:"TraceAgentPort" envconfig:"TRACE_AGENT_PORT"`
	TraceServiceName string   `name:"TraceServiceName" envconfig:"TRACE_SERVICE_NAME"`
	TraceTags        []string `name:"TraceTags" envconfig:"TRACE_TAGS"`
}

func (l *TracerSettingsLoader) LoadFromEnv() (err error) {
	err = config.Load(l)
	return
}

var TracerSettingsModule = fx.Provide(
	func() (out TracerSettingsLoader, err error) {
		err = out.LoadFromEnv()
		return
	},
)
