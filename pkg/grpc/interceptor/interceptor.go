package interceptor

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

type GrpcInterceptor struct{}

func NewGrpcInterceptor() *GrpcInterceptor {
	return &GrpcInterceptor{}
}

func (i *GrpcInterceptor) UnaryLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("--> unary interceptor:", info.FullMethod)
	// log.Println(metadata.FromIncomingContext(ctx))
	return handler(ctx, req)
}

func (i *GrpcInterceptor) StreamLogger(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Println("--> stream interceptor:", info.FullMethod)
	// log.Println(metadata.FromIncomingContext(stream.Context()))
	return handler(srv, stream)
}
