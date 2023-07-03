package utils

import (
	"context"
	"crypto/tls"
	"time"

	"google.golang.org/grpc/credentials"

	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

var (
	InsecureDialOption = grpc.WithInsecure()
	SecureDialOption   = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{}))
)

var blockingDialTimeout = 5 * time.Second

func SetBlockingDialTimeout(d time.Duration) {
	blockingDialTimeout = d
}

func BlockingDial(target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	cxt, cancel := context.WithTimeout(context.Background(), blockingDialTimeout)
	defer cancel()

	opts = append(opts, grpc.WithBlock())
	return DialContext(cxt, target, opts...)
}

// Return a connection
func Dial(target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
	return DialContext(context.Background(), target, opts...)
}

// Return a contextual connection
func DialContext(
	cxt context.Context,
	target string,
	opts ...grpc.DialOption,
) (conn *grpc.ClientConn, err error) {
	return grpc.DialContext(cxt, target, opts...)
}

func InjectTracingHeaders(opts ...grpc.DialOption) []grpc.DialOption {
	tracer := grpc_opentracing.WithTracer(opentracing.GlobalTracer())

	opts = append(opts, grpc.WithStreamInterceptor(grpc_opentracing.StreamClientInterceptor(tracer)))
	opts = append(opts, grpc.WithUnaryInterceptor(grpc_opentracing.UnaryClientInterceptor(tracer)))

	return opts
}

func TransportSecurity(secure bool) grpc.DialOption {
	if secure {
		return SecureDialOption
	} else {
		return InsecureDialOption
	}
}
