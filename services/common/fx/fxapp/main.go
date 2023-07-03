package fxapp

import (
	"context"
	"time"

	"go.uber.org/fx"
)

type App struct {
	*fx.App
	controller *Controller
}

func (app *App) Run() error {
	ctx := context.Background()

	if err := app.Start(ctx); err != nil {
		return err
	} else {
		select {
		case <-app.controller.quitChannel:
			if err := app.Stop(ctx); err != nil {
				return err
			}
		case <-app.Done():
		}

		return nil
	}
}

func NewApp(opts ...fx.Option) *App {
	modules := fx.Options(opts...)

	type ControllerOut struct {
		fx.Out
		Controller *Controller
	}

	controller := newController()

	return &App{
		App: fx.New(
			fx.Provide(func() ControllerOut {
				return ControllerOut{Controller: controller}
			}),
			modules,
		),
		controller: controller,
	}
}

func Main(opts ...fx.Option) error {
	app := NewApp(opts...)
	if err := app.Run(); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if err := app.Stop(ctx); err != nil {
		return err
	}

	return nil
}
