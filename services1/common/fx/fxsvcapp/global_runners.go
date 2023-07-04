package fxsvcapp

import (
	"context"

	"github.com/donglei1234/platform/services/common/runner"
	"go.uber.org/fx"
)

type GlobalRunners struct {
	fx.In

	Runners []runner.Runner `group:"Runner"`
}

type RunnerFactory struct {
	fx.Out

	Runner runner.Runner `group:"Runner"`
}

func (g *GlobalRunners) Execute(lc fx.Lifecycle) error {
	runnerHook := func(r runner.Runner) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				return r.Start()
			},
			OnStop: func(ctx context.Context) error {
				return r.Stop()
			},
		})
	}

	for _, r := range g.Runners {
		runnerHook(r)
	}

	return nil
}

func (f *RunnerFactory) Execute(r runner.Runner) error {
	f.Runner = r
	return nil
}
