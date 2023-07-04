package datadog

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	dd_grpc "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
	dd_http "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
	dd_tracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/donglei1234/platform/services/common/tracing"
)

const Provider = "datadog"

type tracer struct {
	logger  *zap.Logger
	address string
	name    string
	tags    map[string]string

	mux *dd_http.ServeMux
}

func NewTracer(
	logger *zap.Logger,
	host string,
	port int,
	name string,
	tags ...string,
) (t *tracer, e error) {
	if ts, err := parseTags(tags...); err != nil {
		e = err
	} else {
		t = &tracer{
			logger:  logger,
			name:    name,
			address: fmt.Sprintf("%s:%d", host, port),
			tags:    ts,
		}
	}

	return
}

func (t *tracer) Start() error {
	opts := []dd_tracer.StartOption{
		dd_tracer.WithServiceName(t.name),
		dd_tracer.WithAgentAddr(t.address),
		dd_tracer.WithLogger(
			newErrorLogger(t.logger, t.name),
		),
		dd_tracer.WithAnalytics(true),
	}
	for k, v := range t.tags {
		opts = append(opts, dd_tracer.WithGlobalTag(k, v))
	}

	dd_tracer.Start(opts...)

	return nil
}

func (t *tracer) Stop() error {
	dd_tracer.Stop()

	return nil
}

func (t *tracer) StreamClientInterceptor() grpc.StreamClientInterceptor {
	return dd_grpc.StreamClientInterceptor(t.grpcOptions()...)
}

func (t *tracer) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return dd_grpc.StreamServerInterceptor(t.grpcOptions()...)
}

func (t *tracer) UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return dd_grpc.UnaryClientInterceptor(t.grpcOptions()...)
}

func (t *tracer) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return dd_grpc.UnaryServerInterceptor(t.grpcOptions()...)
}

func (t *tracer) ServeMux() tracing.HttpServeMux {
	if t.mux == nil {
		t.mux = dd_http.NewServeMux(t.httpOptions()...)
	}

	return t.mux
}

func (t *tracer) SpanContext(ctx context.Context) (sctx tracing.SpanContext, ok bool) {
	var span dd_tracer.Span
	if span, ok = dd_tracer.SpanFromContext(ctx); ok {
		sctx = span.Context()
	}

	return
}

func (t *tracer) grpcOptions() []dd_grpc.Option {
	return []dd_grpc.Option{
		dd_grpc.WithServiceName(t.name),
		dd_grpc.WithAnalytics(true),
	}
}

func (t *tracer) httpOptions() []dd_http.Option {
	return []dd_http.Option{
		dd_http.WithServiceName(t.name),
		dd_http.WithAnalytics(true),
	}
}

type errorLogger struct {
	logger *zap.Logger
	name   string
}

func newErrorLogger(logger *zap.Logger, name string) *errorLogger {
	return &errorLogger{
		logger: logger,
		name:   name,
	}
}

func (l *errorLogger) Log(msg string) {
	l.logger.Error(
		msg,
		zap.String("tracer", l.name),
	)
}

func parseTags(tags ...string) (ts map[string]string, err error) {
	ts = make(map[string]string)

	for _, t := range tags {
		ps := strings.Split(t, "=")
		if len(ps) != 2 {
			err = errors.Errorf("invalid trace tag '%s'", t)
			break
		}

		ts[ps[0]] = ps[1]
	}

	return
}
