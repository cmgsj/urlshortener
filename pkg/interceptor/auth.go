package interceptor

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	InvalidAuthHeader = status.Error(codes.Unauthenticated, "invalid auth header")
)

type AuthFunc func(ctx context.Context) error

type Auth struct {
	logger   *zap.Logger
	authFunc AuthFunc
}

func NewAuth(logger *zap.Logger, authFunc AuthFunc) *Auth {
	return &Auth{logger: logger, authFunc: authFunc}
}

func Authenticate(ctx context.Context) error {
	data, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(data["auth"]) != 1 {
		return InvalidAuthHeader
	}
	_ = data["auth"][0]
	// TODO: Authenticate token
	return nil
}

func (a *Auth) Unary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if err := a.authFunc(ctx); err != nil {
		a.logger.Error("authentication failed", zap.Error(err))
		return nil, err
	}
	return handler(ctx, req)
}

func (a *Auth) Stream(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if err := a.authFunc(stream.Context()); err != nil {
		a.logger.Error("authentication failed", zap.Error(err))
		return err
	}
	return handler(srv, stream)
}
