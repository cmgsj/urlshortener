package cache

import (
	"context"
	"time"
	"urlshortener/pkg/proto/cachepb"
	"urlshortener/pkg/proto/healthpb"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	UrlNotFoundError    = status.Error(codes.NotFound, "url not found")
	InternalServerError = status.Error(codes.Internal, "internal server error")
)

type Service struct {
	cachepb.UnimplementedCacheServiceServer
	healthpb.UnimplementedHealthServiceServer
	rdb             *redis.Client
	cacheExpiryTime time.Duration
}

func (service *Service) GetUrl(ctx context.Context, req *cachepb.GetUrlRequest) (*cachepb.GetUrlResponse, error) {
	redirectUrl, err := service.rdb.Get(ctx, req.GetUrlId()).Result()
	if err != nil {
		return nil, UrlNotFoundError
	}
	return &cachepb.GetUrlResponse{RedirectUrl: redirectUrl}, nil
}

func (service *Service) SetUrl(ctx context.Context, req *cachepb.SetUrlRequest) (*cachepb.NoContent, error) {
	err := service.rdb.Set(ctx, req.GetUrl().GetUrlId(), req.GetUrl().GetRedirectUrl(), service.cacheExpiryTime).Err()
	if err != nil {
		return nil, InternalServerError
	}
	return &cachepb.NoContent{}, nil
}

func (service *Service) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (service *Service) Watch(req *healthpb.HealthCheckRequest, stream healthpb.HealthService_WatchServer) error {
	for {
		select {
		case <-stream.Context().Done():
			return nil
		default:
			err := stream.Send(&healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING})
			if err != nil {
				return err
			}
			time.Sleep(time.Minute)
		}
	}
}
