package fx

import (
	"go.uber.org/fx"

	"github.com/donglei1234/platform/services/common/config"
)

type BuddySettings struct {
	fx.In

	BuddyUrl        string `name:"BuddyUrl"`
	InviterMaxCount int32  `name:"InviterMaxCount"`
	BuddyMaxCount   int32  `name:"BuddyMaxCount"`
	BlockedMaxCount int32  `name:"BlockedMaxCount"`
	Name            string `name:"Uid"`
}

type BuddySettingsLoader struct {
	fx.Out
	config.EnvironmentBlock

	BuddyUrl        string `name:"BuddyUrl" envconfig:"BUDDY_URL" default:"localhost:8081"`
	BuddyMaxCount   int32  `name:"BuddyMaxCount" envconfig:"BUDDY_MAX_COUNT" default:"100"`
	BlockedMaxCount int32  `name:"BlockedMaxCount" envconfig:"BLOCKED_MAX_COUNT" default:"100"`
	InviterMaxCount int32  `name:"InviterMaxCount" envconfig:"INVITER_MAX_COUNT" default:"20"`
	Name            string `name:"Uid" envconfig:"NAME" default:"buddy"`
}

func (g *BuddySettingsLoader) LoadFromEnv() (err error) {
	err = config.Load(g)
	return
}

var BuddySettingsModule = fx.Provide(
	func() (out BuddySettingsLoader, err error) {
		err = out.LoadFromEnv()
		return
	},
)
