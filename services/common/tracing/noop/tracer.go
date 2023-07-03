package noop

import (
	"context"
	"net/http"

	"google.golang.org/grpc"

	"github.com/donglei1234/platform/services/common/tracing"
)

type tracer struct {
	mux *http.ServeMux
}

func NewTracer() *tracer {
	return &tracer{}
}

func (t *tracer) Start() error {
	return nil
}

func (t *tracer) Stop() error {
	return nil
}

func (t *tracer) StreamClientInterceptor() grpc.StreamClientInterceptor {
	return streamClientInterceptor
}

func (t *tracer) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return streamServerInterceptor
}

func (t *tracer) UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return unaryClientInterceptor
}

func (t *tracer) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return unaryServerInterceptor
}

func (t *tracer) ServeMux() tracing.HttpServeMux {
	if t.mux == nil {
		t.mux = http.NewServeMux()
	}

	return t.mux
}

func (t *tracer) SpanContext(ctx context.Context) (_ tracing.SpanContext, _ bool) {
	return
}

var streamClientInterceptor = func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return streamer(ctx, desc, cc, method, opts...)
}

var streamServerInterceptor = func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return handler(srv, ss)
}

var unaryClientInterceptor = func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return invoker(ctx, method, req, reply, cc, opts...)
}

var unaryServerInterceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	return handler(ctx, req)
}
