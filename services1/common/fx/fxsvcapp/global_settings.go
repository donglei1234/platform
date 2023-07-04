package fxsvcapp

import (
	"go.uber.org/fx"

	"github.com/donglei1234/platform/services/common/config"
)

// GlobalSettings loads from the fx dependency graph based on values provided by GlobalSettingsLoader.
type GlobalSettings struct {
	fx.In

	NetworkSettings
	KeepaliveSettings
	SocketIOSettings

	AppTestMode bool   `name:"AppTestMode"`
	Version     string `name:"Version"`
	Deployment  string `name:"Deployment"`
	AppId       string `name:"AppId"`
	PublicUrl   string `name:"PublicUrl"`

	NatsUrl  string `name:"NatsUrl"`
	RedisUrl string `name:"RedisUrl"`
	RedisPwd string `name:"RedisPwd"`

	AwsAccessKey         string `name:"AwsAccessKey"`
	AwsAccessSecret      string `name:"AwsAccessSecret"`
	AwsRegion            string `name:"AwsRegion"`
	AwsSNSApplicationARN string `name:"AwsSNSApplicationARN"`
	TopicArn             string `name:"TopicArn"`
	AwsBucketName        string `name:"AwsBucketName"`

	FacebookAppId     string `name:"FacebookAppId"`
	FacebookAppSecret string `name:"FacebookAppSecret"`
	GoogleApiKeys     string `name:"GoogleApiKeys"`

	MixpanelToken     string `name:"MixpanelToken"`
	ThinkingDataToken string `name:"ThinkingDataToken"`

	JwtVerificationKey string `name:"JwtVerificationKey"`

	IAPCredential string `name:"IAPCredential"`

	StatHost string `name:"StatHost"`
}

// GlobalSettingsLoader loads from the environment and its members are injected into the fx dependency graph.
type GlobalSettingsLoader struct {
	fx.Out
	config.EnvironmentBlock

	NetworkSettingsLoader
	KeepaliveSettingsLoader
	SocketIOSettingsLoader

	AppTestMode bool   `name:"AppTestMode" ignored:"true"`
	Version     string `name:"Version" default:"unknown"`
	Deployment  string `name:"Deployment" default:"local" envconfig:"DEPLOYMENT"`
	AppId       string `name:"AppId" envconfig:"APP_ID" default:"app_id"`
	PublicUrl   string `name:"PublicUrl" default:"localhost:8081" envconfig:"PUBLIC_URL"`

	NatsUrl  string `name:"NatsUrl" default:"nats://nats:4222" envconfig:"NATS_URL" vault:"nats_url"`
	RedisUrl string `name:"RedisUrl" default:"localhost:6379" envconfig:"REDIS_URL" vault:"redis_url"`
	RedisPwd string `name:"RedisPwd" default:"" envconfig:"REDIS_PWD" vault:"redis_pwd"`

	VaultAddr  string `name:"VaultAddr" default:"http://10.0.1.3:8200" envconfig:"VAULT_ADDR"`
	VaultPath  string `name:"VaultPath" default:"/kv/data/test" envconfig:"VAULT_PATH"`
	VaultToken string `name:"ValutToken" default:"s.ZgcvSsbtRNudg8qz8UmOrXy9" envconfig:"VAULT_TOKEN"`

	AwsAccessKey         string `name:"AwsAccessKey" default:"" envconfig:"AWS_ACCESS_KEY" vault:"aws_key"`
	AwsAccessSecret      string `name:"AwsAccessSecret" default:"" envconfig:"AWS_ACCESS_SECRET" vault:"aws_secret"`
	AwsRegion            string `name:"AwsRegion" default:"" envconfig:"AWS_REGION" vault:"aws_region"`
	AwsBucketName        string `name:"AwsBucketName" default:"" envconfig:"AWS_BUCKET_NAME" vault:"aws_bucket_name"`
	AwsSNSApplicationARN string `name:"AwsSNSApplicationARN" default:"{}" envconfig:"AWS_SNS_APPLICATION_ARN" vault:"aws_sns_application_arn"`
	TopicArn             string `name:"TopicArn" default:"{}" envconfig:"TOPIC_ARN" vault:"topic_arn"`

	FacebookAppId     string `name:"FacebookAppId" default:"" envconfig:"FACEBOOK_APP_ID" vault:"fb_id"`
	FacebookAppSecret string `name:"FacebookAppSecret" default:"" envconfig:"FACEBOOK_APP_SECRET" vault:"fb_secret"`
	GoogleApiKeys     string `name:"GoogleApiKeys" default:"" envconfig:"GOOGLE_API_KEYS" vault:"google"`

	MixpanelToken     string `name:"MixpanelToken" default:"" envconfig:"MIXPANEl_TOKEN" vault:"mixpanel"`
	ThinkingDataToken string `name:"ThinkingDataToken" default:"" envconfig:"THINKING_DATA_TOKEN" vault:"thinkingdata"`

	JwtVerificationKey string `name:"JwtVerificationKey" default:"test" envconfig:"JWT_VERIFICATION_KEY" vault:"jwt"`

	IAPCredential string `name:"IAPCredential" default:"" envconfig:"IAP_CREDENTIAL" vault:"iap"`
	StatHost      string `name:"StatHost" default:"" envconfig:"STAT_HOST" vault:"stat_host"`
}

func (g *GlobalSettingsLoader) LoadFromEnv() (err error) {
	err = config.Load(g)
	return
}

var SettingsModule = fx.Provide(
	func() (out GlobalSettingsLoader, err error) {
		err = out.LoadFromEnv()
		return
	},
)
