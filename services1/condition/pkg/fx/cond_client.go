package fx

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/donglei1234/platform/services/common/nosql/document"
	condition "github.com/donglei1234/platform/services/condition/pkg/client"
)

type ConditionClient struct {
	fx.In
	ConditionClient condition.ConditionClient `name:"ConditionClient"`
}

type ConditionClientFactory struct {
	fx.Out
	ConditionClient condition.ConditionClient `name:"ConditionClient"`
}

// using this locally to make the condition/session store situationally optional
type ConditionClientParams struct {
	fx.In
	ConditionSettings
	ConditionStore document.DocumentStore `name:"ConditionStore" optional:"true"`
	SessionStore   document.DocumentStore `name:"SessionStore" optional:"true"`
}

func (g *ConditionClientFactory) Execute(
	l *zap.Logger,
	t fxsvcapp.GlobalTracer,
	s fxsvcapp.SecuritySettings,
	a ConditionClientParams,
) (err error) {
	l.Info("connect",
		zap.String("service", "condition client"),
		zap.String("url", a.ConditionUrl),
	)
	g.ConditionClient, err = condition.NewConditionClient(a.ConditionUrl, s.SecureClients)

	return
}

var ConditionClientModule = fx.Provide(
	func(
		l *zap.Logger,
		t fxsvcapp.GlobalTracer,
		s fxsvcapp.SecuritySettings,
		a ConditionClientParams,
	) (out ConditionClientFactory, err error) {
		err = out.Execute(l, t, s, a)
		return
	},
)
