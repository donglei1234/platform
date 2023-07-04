package service

import (
	"github.com/aceld/zinx/ziface"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/donglei1234/platform/services/common/access"
	"github.com/donglei1234/platform/services/common/server"
)

type HasHttpServeMux interface {
	HttpServeMux() server.HttpServeMux
}

type HasSocketIOServer interface {
	HttpServeMux() server.HttpServeMux
	SocketIOServer() *socketio.Server
}

type HasZinxTCPServer interface {
	ZinxTcpServer() ziface.IServer
}

type HasGrpcServer interface {
	GrpcServer() *grpc.Server
}

type HasGatewayServer interface {
	GatewayServer() *http.Server
	GatewayRuntimeMux() *runtime.ServeMux
	GatewayOption() []grpc.DialOption
}

type Service interface {
	access.ProtectedService
	ServiceTransport() Transport
}

type GrpcService interface {
	Service
	RegisterWithGrpcServer(server HasGrpcServer) error
}

type HttpService interface {
	Service
	RegisterWithHttpServer(server HasHttpServeMux)
}

type SocketIOService interface {
	Service
	RegisterWithSocketIOServer(server HasSocketIOServer)
}

type GatewayService interface {
	Service
	RegisterWithGatewayServer(server HasGatewayServer) error
}

type ZinxTcpService interface {
	Service
	RegisterWithTCPServer(server HasZinxTCPServer)
}
