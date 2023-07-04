package internal

import (
	"context"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"go.uber.org/zap"
)

type ZinxTcpServer struct {
	logger *zap.Logger
	mux    HttpServeMux
	//listener HasHttpListener
	port Port
	sio  ziface.IServer
}

func (s *ZinxTcpServer) StartServing(ctx context.Context) error {
	go func() {
		s.sio.Serve()
	}()
	return nil
}

func (s *ZinxTcpServer) StopServing(ctx context.Context) error {
	s.sio.Stop()
	return nil
}

func (s *ZinxTcpServer) ZinxTcpServer() ziface.IServer {
	return s.sio
}

func (s *ZinxTcpServer) Port() Port {
	return s.port
}

func NewZinxTcpServer(
	logger *zap.Logger,
	listener HasHttpListener,
	port Port,
) (result *ZinxTcpServer, err error) {
	sio := znet.NewServer()
	result = &ZinxTcpServer{
		logger: logger,
		//listener: listener,
		port: port,
		sio:  sio,
	}
	return
}
