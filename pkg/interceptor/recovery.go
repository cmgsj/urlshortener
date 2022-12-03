package interceptor

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrInternal = status.Error(codes.Internal, "internal error")

type Recovery struct {
	logger *zap.Logger
}

func NewRecovery(logger *zap.Logger) *Recovery {
	return &Recovery{logger: logger}
}

func (r *Recovery) Unary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var err error
	panicked := true
	defer func() {
		if rec := recover(); rec != nil || panicked {
			r.logger.Error("panic recovered", zap.Any("err", rec))
			err = ErrInternal
		}
	}()
	res, err := handler(ctx, req)
	panicked = false
	return res, err
}

func (r *Recovery) Stream(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	var err error
	panicked := true
	defer func() {
		if rec := recover(); rec != nil || panicked {
			r.logger.Error("panic recovered", zap.Any("err", rec))
			err = ErrInternal
		}
	}()
	err = handler(srv, stream)
	panicked = false
	return err
}
