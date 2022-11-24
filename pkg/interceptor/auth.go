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
	InvalidAuthHeader    = status.Error(codes.Unauthenticated, "invalid auth header")
	UnauthenticatedError = status.Error(codes.Unauthenticated, "unauthenticated")
)

type Auth struct {
	logger *zap.Logger
}

func NewAuth() *Auth {
	return &Auth{
		logger: zap.Must(zap.NewDevelopment()),
	}
}

func (a *Auth) Unary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	data, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(data["auth"]) != 1 {
		a.logger.Error("No metadata found")
		return nil, InvalidAuthHeader
	}
	token := data["auth"][0]
	a.logger.Info("metadata:", zap.Any("data", data))
	a.logger.Info("auth:", zap.String("token", token))
	return handler(ctx, req)
}

func (a *Auth) Stream(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	data, ok := metadata.FromIncomingContext(stream.Context())
	if !ok || len(data["auth"]) != 1 {
		a.logger.Error("No metadata found")
		return InvalidAuthHeader
	}
	token := data["auth"][0]
	a.logger.Info("metadata:", zap.Any("data", data))
	a.logger.Info("auth:", zap.String("token", token))
	return handler(srv, stream)
}
