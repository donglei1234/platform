package server

import (
	"context"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	cli "github.com/donglei1234/platform/services/auth/pkg/client"
	"github.com/donglei1234/platform/services/auth/tools"
	"github.com/donglei1234/platform/services/common/access"
	"github.com/donglei1234/platform/services/common/tracing"
)

// addMetadataToLog extracts relevant information from the incoming request metadata and adds it to the context logger.
func addMetadataToLog(ctx context.Context, tracer tracing.Tracer) {
	in := metautils.ExtractIncoming(ctx)

	// tag simple metadata
	ctxzap.AddFields(ctx, zap.String("md.user-agent", in.Get("user-agent")))

	// parse any authorization header
	{
		raw := in.Get("authorization")

		if raw != "" {
			e := strings.SplitN(raw, " ", 2)
			if len(e) == 2 {
				//ctxzap.AddFields(ctx, zap.String("md.auth.token", auth.AnonymizeToken(e[1])))
			} else {
				ctxzap.AddFields(ctx, zap.String("md.auth.raw", raw))
			}
		}
	}

	// inject tracing metadata
	if sctx, ok := tracer.SpanContext(ctx); ok {
		ctxzap.AddFields(
			ctx,
			zap.Uint64("dd.trace_id", sctx.TraceID()),
			zap.Uint64("dd.span_id", sctx.SpanID()),
		)
	}

}

func addInterceptorOptions(
	logger *zap.Logger,
	tracer tracing.Tracer,
	version string,
	authClient cli.AuthClient,
	opts ...grpc.ServerOption,
) []grpc.ServerOption {
	var authFunc grpc_auth.AuthFunc

	si := []grpc.StreamServerInterceptor{
		tracer.StreamServerInterceptor(),
		grpc_prometheus.StreamServerInterceptor,
		grpc_zap.StreamServerInterceptor(logger),
		//grpc_auth.UnaryServerInterceptor(authFunc),
		func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			header := metadata.Pairs("version", version)
			if err := ss.SendHeader(header); err != nil {
				return err
			}

			addMetadataToLog(ss.Context(), tracer)
			ctx := context.WithValue(ss.Context(), access.GrpcServerContextKey, srv)
			wrapped := grpc_middleware.WrapServerStream(ss)
			wrapped.WrappedContext = context.WithValue(ctx, "AuthFunc", authFunc)
			return handler(srv, wrapped)
		},
	}

	ui := []grpc.UnaryServerInterceptor{
		tracer.UnaryServerInterceptor(),
		grpc_prometheus.UnaryServerInterceptor,
		grpc_zap.UnaryServerInterceptor(logger),
		//grpc_auth.UnaryServerInterceptor(authFunc),
		grpc_zap.PayloadUnaryServerInterceptor(logger, func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
			return true
		}),
		func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			header := metadata.Pairs("version", version)
			if err = grpc.SendHeader(ctx, header); err != nil {
				return
			}

			addMetadataToLog(ctx, tracer)
			ctx = context.WithValue(ctx, access.GrpcServerContextKey, info.Server)
			ctx = context.WithValue(ctx, "AuthFunc", authFunc)
			return handler(ctx, req)
		},
	}
	//	 configure auth not used ,instead use authOverrideFunc
	{
		authFunc = func(ctx context.Context) (context.Context, error) {
			logger.Debug("platform authfun run !")
			token, err := tools.ParseJwtToken(ctx)
			if err != nil {
				return nil, ErrMiddlewareParseJwtFailure
			}

			_, err = authClient.ValidateToken(ctx, token)
			// ValidateToken已对err封装
			if err != nil {
				return nil, err
			}
			logger.Debug("platform middleware authfun run over , Check passed !")
			return ctx, nil
		}

		//si = append(si, grpc_auth.StreamServerInterceptor(authFunc))
		//ui = append(ui, grpc_auth.UnaryServerInterceptor(authFunc))
	}

	interceptorOpts := []grpc.ServerOption{
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(si...)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(ui...)),
	}

	if opts == nil {
		opts = interceptorOpts
	} else {
		opts = append(opts, interceptorOpts...)
	}

	return opts
}
