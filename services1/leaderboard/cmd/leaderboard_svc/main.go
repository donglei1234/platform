package main

import (
	"math/rand"
	"os"
	"path"
	"time"

	auth "github.com/donglei1234/platform/services/auth/pkg/module"
	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/donglei1234/platform/services/leaderboard/internal/app/service/external"
	"github.com/donglei1234/platform/services/leaderboard/internal/app/service/public"
	ffx "github.com/donglei1234/platform/services/leaderboard/pkg/fx"
)

func main() {
	rand.Seed(time.Now().Unix())
	home := os.Getenv("HOME")
	setEnvVariable("DOCUMENT_STORE_URL", "badger://badger/"+path.Join(home, "leaderboard"))

	fxsvcapp.StandardMain(
		ffx.LeaderboardSettingsModule,
		public.ServiceModule,
		external.ServicesModules,
		auth.PublicModule,
		fxsvcapp.AuthStoreModule,
	)
}

func setEnvVariable(key string, value string) {
	if os.Getenv(key) == "" {
		_ = os.Setenv(key, value)
	}
}
