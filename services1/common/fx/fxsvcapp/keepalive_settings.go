package fxsvcapp

import (
	"time"

	"go.uber.org/fx"
	"google.golang.org/grpc/keepalive"

	"github.com/donglei1234/platform/services/common/config"
)

type KeepaliveSettings struct {
	fx.In

	KeepaliveClientParameters  KeepaliveClientParameters  `name:"KeepaliveClientParameters"`
	KeepaliveEnforcementPolicy KeepaliveEnforcementPolicy `name:"KeepaliveEnforcementPolicy"`
	KeepaliveServerParameters  KeepaliveServerParameters  `name:"KeepaliveServerParameters"`
}

func (s KeepaliveSettings) GetKeepaliveClientParameters() keepalive.ClientParameters {
	return keepalive.ClientParameters{
		Time:                s.KeepaliveClientParameters.Time,
		Timeout:             s.KeepaliveClientParameters.Timeout,
		PermitWithoutStream: s.KeepaliveClientParameters.PermitWithoutStream,
	}
}

func (s KeepaliveSettings) GetKeepaliveEnforcementPolicy() keepalive.EnforcementPolicy {
	return keepalive.EnforcementPolicy{
		MinTime:             s.KeepaliveEnforcementPolicy.MinTime,
		PermitWithoutStream: s.KeepaliveEnforcementPolicy.PermitWithoutStream,
	}
}

func (s KeepaliveSettings) GetKeepaliveServerParameters() keepalive.ServerParameters {
	return keepalive.ServerParameters{
		MaxConnectionIdle:     s.KeepaliveServerParameters.MaxConnectionIdle,
		MaxConnectionAge:      s.KeepaliveServerParameters.MaxConnectionAge,
		MaxConnectionAgeGrace: s.KeepaliveServerParameters.MaxConnectionAgeGrace,
		Time:                  s.KeepaliveServerParameters.Time,
		Timeout:               s.KeepaliveServerParameters.Timeout,
	}
}

type KeepaliveSettingsLoader struct {
	fx.Out
	config.EnvironmentBlock

	KeepaliveClientParameters  KeepaliveClientParameters  `name:"KeepaliveClientParameters"`
	KeepaliveEnforcementPolicy KeepaliveEnforcementPolicy `name:"KeepaliveEnforcementPolicy"`
	KeepaliveServerParameters  KeepaliveServerParameters  `name:"KeepaliveServerParameters"`
}

func (g *KeepaliveSettingsLoader) LoadFromEnv() (err error) {
	err = config.Load(g)

	return
}

type KeepaliveClientParameters struct {
	Time                time.Duration `envconfig:"KEEPALIVE_TIME" default:"30s"`
	Timeout             time.Duration `envconfig:"KEEPALIVE_TIMEOUT" default:"20s"`
	PermitWithoutStream bool          `envconfig:"KEEPALIVE_PERMIT_WITHOUT_STREAM" default:"true"`
}

type KeepaliveEnforcementPolicy struct {
	MinTime             time.Duration `envconfig:"KEEPALIVE_MIN_TIME" default:"30s"`
	PermitWithoutStream bool          `envconfig:"KEEPALIVE_PERMIT_WITHOUT_STREAM" default:"true"`
}

type KeepaliveServerParameters struct {
	MaxConnectionIdle     time.Duration `envconfig:"KEEPALIVE_MAX_CONNECTION_IDLE"`
	MaxConnectionAge      time.Duration `envconfig:"KEEPALIVE_MAX_CONNECTION_AGE"`
	MaxConnectionAgeGrace time.Duration `envconfig:"KEEPALIVE_MAX_CONNECTION_AGE_GRACE"`
	Time                  time.Duration `envconfig:"KEEPALIVE_TIME" default:"2h"`
	Timeout               time.Duration `envconfig:"KEEPALIVE_TIMEOUT" default:"20s"`
}

var KeepaliveSettingsModule = fx.Provide(
	func() (out KeepaliveSettingsLoader, err error) {
		err = out.LoadFromEnv()
		return
	},
)
