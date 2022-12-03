package grpcutil

import (
	"context"
	"log"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func DialGrpc(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to connect:", err)
	}
	return conn
}

func CheckService(name string, logger *zap.Logger, client healthpb.HealthClient, serviceName string, active *bool) {
	ctx, cancel := MakeUnaryCtx()
	defer cancel()
	_, err := client.Check(ctx, &healthpb.HealthCheckRequest{Service: name})
	if err != nil {
		*active = false
		logger.Error("failed to check:", zap.String("service", serviceName), zap.Error(err))
	} else {
		*active = true
		logger.Info("service is active", zap.String("service", serviceName))
	}
}

func WatchService(name string, logger *zap.Logger, client healthpb.HealthClient, serviceName string, active *bool) {
	ctx, cancel := MakeStreamCtx()
	defer cancel()
	stream, err := client.Watch(ctx, &healthpb.HealthCheckRequest{Service: name})
	if err != nil {
		*active = false
		logger.Error("failed to watch:", zap.String("service", serviceName), zap.Error(err))
		return
	}
	for {
		res, err := stream.Recv()
		if err != nil {
			*active = false
			logger.Error("failed to receive:", zap.String("service", serviceName), zap.Error(err))
			return
		}
		if res.GetStatus() == healthpb.HealthCheckResponse_SERVING {
			*active = true
			logger.Info("service is active", zap.String("service", serviceName))
		} else {
			*active = false
			logger.Info("service is not active", zap.String("service", serviceName))
		}
	}
}

func MakeUnaryCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second)

}

func MakeStreamCtx() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}
