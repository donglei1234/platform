package internal

import (
	"context"
	"go.uber.org/zap"
	"net/http"

	"github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
)

type TcpSocketIOServer struct {
	logger   *zap.Logger
	server   *http.Server
	mux      HttpServeMux
	listener HasHttpListener
	port     Port
	sio      *socketio.Server
}

func (s *TcpSocketIOServer) StartServing(ctx context.Context) error {
	if listener, err := s.listener.HttpListener(); err != nil {
		return err
	} else {
		s.logger.Info(
			"serving socket.io",
			zap.String("network", listener.Addr().Network()),
			zap.String("address", listener.Addr().String()),
			zap.Int("port", s.port.Value()),
		)
		go func() {
			if err := s.sio.Serve(); err != nil {
				panic(err)
			}
		}()
		go func() {
			if err := s.server.Serve(listener); err != http.ErrServerClosed {
				panic(err)
			}
		}()
	}

	return nil
}

func (s *TcpSocketIOServer) StopServing(ctx context.Context) error {
	if err := s.sio.Close(); err != nil {
		return err
	}
	return s.server.Shutdown(ctx)
}

func (s *TcpSocketIOServer) HttpServeMux() HttpServeMux {
	return s.mux
}

func (s *TcpSocketIOServer) SocketIOServer() *socketio.Server {
	return s.sio
}

func (s *TcpSocketIOServer) Port() Port {
	return s.port
}

func NewTcpSocketIOServer(
	logger *zap.Logger,
	listener HasHttpListener,
	port Port,
	opt *engineio.Options,
) (result *TcpSocketIOServer, err error) {
	mux := http.NewServeMux()
	sio := socketio.NewServer(opt)
	result = &TcpSocketIOServer{
		logger:   logger,
		listener: listener,
		port:     port,
		mux:      mux,
		server: &http.Server{
			Handler: mux,
		},
		sio: sio,
	}
	return
}
