package interceptor

import (
	"context"
	"urlshortener/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	InvalidAuthHeader    = status.Error(codes.Unauthenticated, "invalid auth header")
	UnauthenticatedError = status.Error(codes.Unauthenticated, "unauthenticated")
)

type AuthInterceptor struct{}

func NewAuthInterceptor() *AuthInterceptor {
	return &AuthInterceptor{}
}

func (i *AuthInterceptor) Unary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	data, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(data["auth"]) != 1 {
		logger.Error("No metadata found")
		return nil, InvalidAuthHeader
	}
	token := data["auth"][0]
	logger.Info("metadata:", data)
	logger.Info("auth:", token)
	return handler(ctx, req)
}

func (i *AuthInterceptor) Stream(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	data, ok := metadata.FromIncomingContext(stream.Context())
	if !ok || len(data["auth"]) != 1 {
		logger.Error("No metadata found")
		return InvalidAuthHeader
	}
	token := data["auth"][0]
	logger.Info("metadata:", data)
	logger.Info("auth:", token)
	return handler(srv, stream)
}
