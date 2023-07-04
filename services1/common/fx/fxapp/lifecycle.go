package fxapp

import (
	"log"
	"os"

	"go.uber.org/fx"
)

// NewAppLifecycleLogger manages enabling / disabling Fx lifecycle logging based on the value of the FX_DEBUG
// environment variable.  Set FX_DEBUG=true to enable logging of Fx events.
func NewAppLifecycleLogger() fx.Option {
	if debug := os.Getenv("FX_DEBUG"); debug != "true" {
		return fx.NopLogger
	} else {
		return fx.WithLogger(log.New(os.Stderr, "", log.LstdFlags))
	}
}
