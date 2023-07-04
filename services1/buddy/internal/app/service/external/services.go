package external

import (
	"go.uber.org/fx"

	ffx "github.com/donglei1234/platform/services/buddy/pkg/fx"
)

var ServiceModules = fx.Options(
	ffx.BuddySettingsModule,
)
