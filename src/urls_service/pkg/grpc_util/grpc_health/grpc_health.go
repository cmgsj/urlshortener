package grpc_health

import (
	"context"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type HealthClient struct {
	healthpb.HealthClient
	Active *atomic.Bool
	Name   string
}

func CheckService(svcName string, logger *zap.Logger, client *HealthClient, done chan<- string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := client.Check(ctx, &healthpb.HealthCheckRequest{Service: svcName})
	if err != nil {
		client.Active.Store(false)
		logger.Error("failed to check:", zap.String("service", client.Name), zap.Error(err))
	} else {
		client.Active.Store(true)
		logger.Info("service is active", zap.String("service", client.Name))
	}
	done <- client.Name
}

func CheckServices(svcName string, logger *zap.Logger, clients []*HealthClient) {
	done := make(chan string)
	for _, client := range clients {
		go CheckService(svcName, logger, client, done)
	}
}

func WatchService(svcName string, logger *zap.Logger, client *HealthClient, done chan<- string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := client.Watch(ctx, &healthpb.HealthCheckRequest{Service: svcName})
	if err != nil {
		client.Active.Store(false)
		logger.Error("failed to watch:", zap.String("service", client.Name), zap.Error(err))
		done <- client.Name
		return
	}
	for {
		res, err := stream.Recv()
		if err != nil {
			client.Active.Store(false)
			logger.Error("failed to receive:", zap.String("service", client.Name), zap.Error(err))
			done <- client.Name
			return
		}
		if res.GetStatus() == healthpb.HealthCheckResponse_SERVING {
			client.Active.Store(true)
			logger.Info("service is active", zap.String("service", client.Name))
		} else {
			client.Active.Store(false)
			logger.Info("service is not active", zap.String("service", client.Name))
		}
	}
}

func WatchServices(svcName string, logger *zap.Logger, d time.Duration, clients []*HealthClient) {
	m := make(map[string]*HealthClient)
	done := make(chan string)
	for _, client := range clients {
		m[client.Name] = client
		go WatchService(svcName, logger, client, done)
	}
	for {
		select {
		case s := <-done:
			client, ok := m[s]
			if ok {
				logger.Error("unknown service", zap.String("service", s))
				continue
			}
			go WatchService(svcName, logger, client, done)
		default:
			time.Sleep(d)
		}
	}
}
