package grpc_ctx

import (
	"context"
	"time"
)

func MakeUnaryCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second)
}

func MakeStreamCtx() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}
