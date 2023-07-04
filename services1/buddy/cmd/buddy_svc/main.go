package main

import (
	"github.com/donglei1234/platform/services/buddy/internal/app/service/external"
	"github.com/donglei1234/platform/services/buddy/internal/app/service/private"
	"github.com/donglei1234/platform/services/buddy/internal/app/service/public"
	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
)

func main() {
	fxsvcapp.StandardMain(
		external.ServiceModules,
		public.ServiceModule,
		private.ServiceModule,
	)
}
