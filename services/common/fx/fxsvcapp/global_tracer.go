package fxsvcapp

import (
	"context"

	"github.com/donglei1234/platform/services/common/tracing/noop"

	"github.com/donglei1234/platform/services/common/tracing/datadog"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/common/tracing"
	tfx "github.com/donglei1234/platform/services/common/tracing/pkg/fx"
)

type GlobalTracer struct {
	fx.In

	Tracer tracing.Tracer `name:"Tracer"`
}

type GlobalTracerFactory struct {
	fx.Out

	Tracer tracing.Tracer `name:"Tracer"`
}

func (f *GlobalTracerFactory) Execute(
	lc fx.Lifecycle,
	l *zap.Logger,
	s GlobalSettings,
	t tfx.TracerSettings,
) (err error) {
	if t.TraceProvider != "" && t.TraceAgentHost != "" && t.TraceAgentPort > 0 {
		if t.TraceServiceName == "" {
			err = ErrMissingTracerServiceName
		} else {
			switch t.TraceProvider {
			case datadog.Provider:
				f.Tracer, err = datadog.NewTracer(
					l,
					t.TraceAgentHost,
					t.TraceAgentPort,
					t.TraceServiceName,
					t.TraceTags...,
				)
			default:
				err = ErrUnsupportedTracer
			}
		}
	} else {
		f.Tracer = noop.NewTracer()
	}

	if f.Tracer != nil {
		lc.Append(fx.Hook{
			OnStart: func(_ context.Context) error {
				return f.Tracer.Start()
			},
			OnStop: func(_ context.Context) error {
				return f.Tracer.Stop()
			},
		})
	}

	return
}

var TracerModule = fx.Provide(
	func(
		lc fx.Lifecycle,
		l *zap.Logger,
		s GlobalSettings,
		t tfx.TracerSettings,
	) (out GlobalTracerFactory, err error) {
		err = out.Execute(lc, l, s, t)
		return
	},
)
