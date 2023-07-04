package service

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
)

// WithoutAuth overrides the default auth behavior and allows all methods to be called without an access token.
type WithoutAuth struct{}

// AuthFuncOverride allows all methods to be unauthenticated.
func (w *WithoutAuth) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	AddDisabledTag(ctx)
	return ctx, nil
}

// WithUnauthenticatedMethods gives the embedding service the option to specify methods that are unauthenticated.
type WithUnauthenticatedMethods struct {
	methods map[string]struct{}
}

// AddUnauthenticatedMethod adds a method to the set of unauthenticated methods.  You must provide the full method name
// to this function.  The full method name is based on the service proto file.  As a format string it'd look something
// like this: `fmt.Sprintf("/%s.%s/%s", proto.package, service.name, method.name)`.
func (w *WithUnauthenticatedMethods) AddUnauthenticatedMethod(fullMethodName string) {
	if !strings.HasPrefix(fullMethodName, "/") {
		panic("must provide full method name")
	}
	if w.methods == nil {
		w.methods = make(map[string]struct{})
	}
	w.methods[fullMethodName] = struct{}{}
}

// AuthFuncOverride allows specific named methods to go unauthenticated.
func (w *WithUnauthenticatedMethods) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	unauthenticated := false
	if w.methods != nil {
		_, unauthenticated = w.methods[fullMethodName]
	}
	if unauthenticated {
		AddDisabledTag(ctx)
		return ctx, nil
	} else {
		authFunc, ok := ctx.Value("AuthFunc").(grpc_auth.AuthFunc)
		if ok && authFunc != nil {
			return authFunc(ctx)
		} else {
			panic("WithUnauthenticatedMethods requires AuthFunc to be in the provided context")
		}
	}
}

// AddDisabledTag sets the auth disabled tag in the provided context tags.
func AddDisabledTag(ctx context.Context) {
	ctxzap.AddFields(ctx, zap.Bool("auth.disabled", true))
}
