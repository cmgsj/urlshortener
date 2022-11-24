package interceptor

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Logger struct {
	logger *zap.Logger
}

func NewLogger() *Logger {
	return &Logger{
		logger: zap.Must(zap.NewDevelopment()),
	}
}

func (l *Logger) Unary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	l.logger.Info("--> Unary Interceptor:", zap.String("method", info.FullMethod))
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		l.logger.Info("metadata:", zap.Any("data", md))
	}
	start := time.Now()
	res, err := handler(ctx, req)
	l.logger.Info("<-- Unary Interceptor:", zap.String("method", info.FullMethod), zap.Duration("elapsed", time.Since(start)), zap.Error(err))
	return res, err
}

func (l *Logger) Stream(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	l.logger.Info("--> Stream Interceptor:", zap.String("method", info.FullMethod))
	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		l.logger.Info("metadata:", zap.Any("data", md))
	}
	start := time.Now()
	err := handler(srv, stream)
	l.logger.Info("<-- Stream Interceptor:", zap.String("method", info.FullMethod), zap.Duration("elapsed", time.Since(start)), zap.Error(err))
	return err
}
