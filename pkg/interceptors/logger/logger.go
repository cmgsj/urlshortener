package logger

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

type loggerInterceptor struct{}

func NewLoggerInterceptor() *loggerInterceptor {
	return &loggerInterceptor{}
}

func (i *loggerInterceptor) UnaryLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("--> Unary Interceptor:", info.FullMethod)
	// log.Println(metadata.FromIncomingContext(ctx))
	return handler(ctx, req)
}

func (i *loggerInterceptor) StreamLogger(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Println("--> Stream Interceptor:", info.FullMethod)
	// log.Println(metadata.FromIncomingContext(stream.Context()))
	return handler(srv, stream)
}
