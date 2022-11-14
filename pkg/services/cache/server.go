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

type cacheServer struct {
	cachepb.UnimplementedCacheServer
	rdb             *redis.Client
	cacheExpiryTime time.Duration
}

func (server *cacheServer) GetUrl(ctx context.Context, req *cachepb.GetUrlRequest) (*cachepb.GetUrlResponse, error) {
	redirectUrl, err := server.rdb.Get(ctx, req.GetUrlId()).Result()
	if err != nil {
		return nil, UrlNotFoundError
	}
	return &cachepb.GetUrlResponse{RedirectUrl: redirectUrl}, nil
}

func (server *cacheServer) SetUrl(ctx context.Context, req *cachepb.SetUrlRequest) (*cachepb.NoContent, error) {
	err := server.rdb.Set(ctx, req.GetUrlId(), req.GetRedirectUrl(), server.cacheExpiryTime).Err()
	if err != nil {
		return nil, InternalServerError
	}
	return &cachepb.NoContent{}, nil
}

func (server *cacheServer) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}
