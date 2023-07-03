package service

import (
	"go.uber.org/fx"
)

type HttpServiceFactory struct {
	fx.Out
	HttpService `group:"HttpService"`
}

func (f *HttpServiceFactory) Execute(service HttpService) error {
	f.HttpService = service
	return nil
}

type GrpcServiceFactory struct {
	fx.Out
	GrpcService `group:"GrpcService"`
}

func (f *GrpcServiceFactory) Execute(service GrpcService) error {
	f.GrpcService = service
	return nil
}

type SocketIOServiceFactory struct {
	fx.Out
	SocketIOService `group:"SocketIOService"`
}

func (f *SocketIOServiceFactory) Execute(service SocketIOService) error {
	f.SocketIOService = service
	return nil
}

type ZinxServiceFactory struct {
	fx.Out
	ZinxTcpService `group:"ZinxTcpService"`
}

//
//func (f *ZinxServiceFactory) Execute(service ZinxTcpService) error {
//	f.ZinxTcpService = service
//	return nil
//}

type GatewayServiceFactory struct {
	fx.Out
	GatewayService `group:"GatewayService"`
}

func (f *GatewayServiceFactory) Execute(service GatewayService) error {
	f.GatewayService = service
	return nil
}
