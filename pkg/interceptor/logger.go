package interceptor

import (
	"context"
	"urlshortener/pkg/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type LoggerInterceptor struct{}

func NewLoggerInterceptor() *LoggerInterceptor {
	return &LoggerInterceptor{}
}

func (i *LoggerInterceptor) Unary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	logger.Info("--> Unary Interceptor:", info.FullMethod)
	logger.Info(metadata.FromIncomingContext(ctx))
	return handler(ctx, req)
}

func (i *LoggerInterceptor) Stream(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	logger.Info("--> Stream Interceptor:", info.FullMethod)
	logger.Info(metadata.FromIncomingContext(stream.Context()))
	return handler(srv, stream)
}
