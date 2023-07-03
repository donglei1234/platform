package main

import (
	"github.com/donglei1234/platform/services/auth/internal/app/service/external"
	"github.com/donglei1234/platform/services/auth/internal/app/service/public"
	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
)

func main() {
	fxsvcapp.StandardMain(
		public.ServiceModule,
		external.ServicesModules,
	)
}
