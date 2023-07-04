package fxsvcapp

import (
	"go.uber.org/fx"

	"github.com/donglei1234/platform/services/common/mq"
)

type GlobalMQ struct {
	fx.In
	MessageQueue mq.MessageQueue `name:"MessageQueue"`
}

type GlobalMQFactory struct {
	fx.Out
	MessageQueue mq.MessageQueue `name:"MessageQueue"`
}

type MQImplementations struct {
	fx.In
	KafkaMQ mq.MessageQueue `name:"KafkaMQ" optional:"true"`
	NatsMQ  mq.MessageQueue `name:"NatsMQ" optional:"true"`
	NsqMQ   mq.MessageQueue `name:"NsqMQ" optional:"true"`
	LocalMQ mq.MessageQueue `name:"LocalMQ" optional:"true"`
}

func (g *GlobalMQFactory) Execute(s GlobalSettings, i MQImplementations) (err error) {
	mq.SetNamespace(s.Deployment)

	// If run in TestMode, all Subscribe() and Publish() requests will be run through
	// the local:// mq implementation regardless of their chosen mq protocol
	if s.AppTestMode {
		g.MessageQueue = mq.NewMessageQueue(i.LocalMQ, i.LocalMQ, i.LocalMQ, i.LocalMQ)
	} else {
		g.MessageQueue = mq.NewMessageQueue(i.KafkaMQ, i.NatsMQ, i.NsqMQ, i.LocalMQ)
	}

	return nil
}

var MessageQueueModule = fx.Provide(
	func(s GlobalSettings, i MQImplementations) (out GlobalMQFactory, err error) {
		err = out.Execute(s, i)
		return
	},
)
