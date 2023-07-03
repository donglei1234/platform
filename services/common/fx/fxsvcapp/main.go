package fxsvcapp

import (
	"go.uber.org/fx"

	afx "github.com/donglei1234/platform/services/auth/pkg/fx"
	"github.com/donglei1234/platform/services/common/fx/fxapp"
	"github.com/donglei1234/platform/services/common/logging"
	nfx "github.com/donglei1234/platform/services/common/nosql/pkg/fx"
	tfx "github.com/donglei1234/platform/services/common/tracing/pkg/fx"
)

var AppNameModule = fx.Provide(
	func() (out fxapp.AppNameLoader, err error) {
		err = out.LoadFromExecutable()
		return
	},
)

var LoggingModule = fx.Options(
	fxapp.NewAppLifecycleLogger(),
	logging.Module,
)

var HttpDebugModule = fx.Provide(
	NewTcpDebugHttpService,
)

var HttpMetricsModule = fx.Provide(
	NewMetricsHttpService,
)

var ServicesModule = fx.Options(
	AppNameModule,
	SettingsModule,
	afx.SettingsModule,
	tfx.TracerSettingsModule,
	LoggingModule,
	ConnectionMuxModule,
	ServersModule,
	TracerModule,
	SecuritySettingsModule,
	AuthClientModule,
	nfx.NoSQLSettingsModule,
	DocumentStoreProviderModule,
	MemoryStoreProviderModule,
	RedisStoreProviderModule,
)

func StandardMain(opts ...fx.Option) {
	if err := fxapp.Main(
		ServicesModule,
		fx.Options(opts...),
		fx.Invoke(cli),
	); err != nil {
		panic(err)
	}
}
