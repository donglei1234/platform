package external

import (
	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/donglei1234/platform/services/common/nosql/memory"
	nfx "github.com/donglei1234/platform/services/common/nosql/pkg/fx"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type GlobalGMServerStoreParams struct {
	fx.In
	GMServerStore memory.MemoryStore `name:"GMServerStore"`
}

type GlobalGMServerStoreResult struct {
	fx.Out
	GMServerStore memory.MemoryStore `name:"GMServerStore"`
}

func (g *GlobalGMServerStoreResult) Execute(
	l *zap.Logger,
	d fxsvcapp.GlobalMemoryStoreProvider,
) (err error) {
	g.GMServerStore, err = openMemoryStore(
		l,
		d.MemoryStoreProvider,
		"gm",
		"gm",
	)
	return
}

func openMemoryStore(
	logger *zap.Logger,
	provider memory.MemoryStoreProvider,
	name string,
	tag string,
) (result memory.MemoryStore, err error) {
	return provider.OpenMemoryStore(name)
}

var GMStoreModule = fx.Provide(
	func(
		l *zap.Logger,
		n nfx.NoSQLSettings,
		d fxsvcapp.GlobalMemoryStoreProvider,
	) (out GlobalGMServerStoreResult, err error) {
		err = out.Execute(l, d)
		return
	},
)
