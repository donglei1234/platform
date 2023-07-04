package internal

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/donglei1234/platform/services/common/utils"
)

const (
	timeoutDuration = 10 * time.Second
)

type TcpGrpcServer struct {
	logger   *zap.Logger
	server   *grpc.Server
	listener HasGrpcListener
	port     Port
}

func (s *TcpGrpcServer) StartServing(ctx context.Context) error {
	if listener, err := s.listener.GrpcListener(); err != nil {
		return err
	} else {
		s.logger.Info(
			"serving grpc",
			zap.String("network", listener.Addr().Network()),
			zap.String("address", listener.Addr().String()),
			zap.Int("port", s.port.Value()),
		)

		go func() {
			if err := s.server.Serve(listener); err != nil && err != grpc.ErrServerStopped {
				panic(err)
			}
		}()
	}

	return nil
}

func (s *TcpGrpcServer) StopServing(ctx context.Context) error {
	s.server.GracefulStop()
	return nil
}

func (s *TcpGrpcServer) Dial(target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	return utils.DialContext(ctx, target, opts...)
}

func (s *TcpGrpcServer) GrpcServer() *grpc.Server {
	return s.server
}

func (s *TcpGrpcServer) Port() Port {
	return s.port
}

func NewTcpGrpcServer(
	logger *zap.Logger,
	listener HasGrpcListener,
	port Port,
	opts ...grpc.ServerOption,
) (result *TcpGrpcServer, err error) {
	result = &TcpGrpcServer{
		logger:   logger,
		listener: listener,
		port:     port,
		server:   grpc.NewServer(opts...),
	}
	return
}

type UnixGrpcServer struct {
	logger   *zap.Logger
	server   *grpc.Server
	listener HasGrpcListener
	socket   Socket
}

func (s *UnixGrpcServer) StartServing(ctx context.Context) error {
	if listener, err := s.listener.GrpcListener(); err != nil {
		return err
	} else {
		s.logger.Info(
			"serving grpc",
			zap.String("network", listener.Addr().Network()),
			zap.String("address", listener.Addr().String()),
		)

		go func() {
			if err := s.server.Serve(listener); err != nil && err != grpc.ErrServerStopped {
				panic(err)
			}
		}()
	}

	return nil
}

func (s *UnixGrpcServer) StopServing(ctx context.Context) error {
	s.server.GracefulStop()
	return nil
}

func (s *UnixGrpcServer) GrpcServer() *grpc.Server {
	return s.server
}

func (s *UnixGrpcServer) Dial(target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	anyScheme := func(t string) bool {
		return strings.Index(t, "//") == 0
	}
	missingScheme := func(t string) bool {
		return strings.Index(t, "://") == -1
	}
	unixScheme := func(t string) bool {
		return strings.Index(t, "unix://") == 0
	}

	var unixTarget string
	switch {
	case anyScheme(target):
		unixTarget = "unix" + ":" + target
	case missingScheme(target):
		unixTarget = "unix" + "://" + target
	case unixScheme(target):
		unixTarget = target
	default:
		err = errors.New("target must specify unix scheme or none at all")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	return utils.DialContext(ctx, unixTarget, opts...)
}

func (s *UnixGrpcServer) Socket() Socket {
	return s.socket
}

func NewUnixGrpcServer(
	logger *zap.Logger,
	listener HasGrpcListener,
	socket Socket,
	opts ...grpc.ServerOption,
) (result *UnixGrpcServer, err error) {
	result = &UnixGrpcServer{
		logger:   logger,
		listener: listener,
		socket:   socket,
		server:   grpc.NewServer(opts...),
	}

	return
}

// TestGrpcServer is an in-process grpc server that eschews socket communication in favor of a shared db buffer.
type TestGrpcServer struct {
	logger   *zap.Logger
	server   *grpc.Server
	listener *bufconn.Listener
	port     Port
}

// Start starts serving the grpc services inside a goroutine.  You must call Stop at a later point to stop serving.
func (s *TestGrpcServer) StartServing(ctx context.Context) error {
	s.logger.Info(
		"serving test grpc",
		zap.String("network", s.listener.Addr().Network()),
		zap.String("address", s.listener.Addr().String()),
		zap.Int("port", s.port.Value()),
	)
	go func() {
		if err := s.server.Serve(s.listener); err != nil && err != grpc.ErrServerStopped {
			panic(err)
		}
	}()
	return nil
}

// Stop stops serving grpc services previously started with Start.
func (s *TestGrpcServer) StopServing(ctx context.Context) error {
	s.server.Stop()
	return nil
}

// Dial creates a new connection to the server and returns it.
func (s *TestGrpcServer) Dial(target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	opts = append(opts,
		grpc.WithInsecure(),
		grpc.WithDialer(func(name string, duration time.Duration) (net.Conn, error) {
			return s.listener.Dial()
		}),
	)

	return utils.Dial(target, opts...)
}

func (s *TestGrpcServer) GrpcServer() *grpc.Server {
	return s.server
}

func (s *TestGrpcServer) Port() Port {
	return s.port
}

// NewTestGrpcServer creates a TestGrpcServer.
func NewTestGrpcServer(
	logger *zap.Logger,
	port Port,
	opts ...grpc.ServerOption,
) *TestGrpcServer {
	return &TestGrpcServer{
		logger:   logger,
		listener: bufconn.Listen(256 * 1024),
		port:     port,
		server:   grpc.NewServer(opts...),
	}
}
