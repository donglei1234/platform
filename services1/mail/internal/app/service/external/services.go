package external

import (
	"go.uber.org/fx"

	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
)

var ServicesModules = fx.Options(
	fxsvcapp.RoomStoreModule,
)
