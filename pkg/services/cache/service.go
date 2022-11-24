package cache

import (
	"context"
	"time"
	"urlshortener/pkg/proto/cachepb"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/status"
)

var (
	ErrUrlNotFound         = status.Error(codes.NotFound, "url not found")
	ErrInternalServerError = status.Error(codes.Internal, "internal server error")
)

type Service struct {
	cachepb.UnimplementedCacheServiceServer
	healthServer *health.Server
	logger       *zap.Logger
	rdb          *redis.Client
	cacheExpTime time.Duration
}

func (s *Service) Get(ctx context.Context, req *cachepb.GetRequest) (*cachepb.GetResponse, error) {
	redirectUrl, err := s.rdb.Get(ctx, req.GetKey()).Result()
	if err != nil {
		return nil, ErrUrlNotFound
	}
	return &cachepb.GetResponse{Value: redirectUrl}, nil
}

func (s *Service) Set(ctx context.Context, req *cachepb.SetRequest) (*cachepb.NoContent, error) {
	err := s.rdb.Set(ctx, req.GetKey(), req.GetValue(), s.cacheExpTime).Err()
	if err != nil {
		return nil, ErrInternalServerError
	}
	return &cachepb.NoContent{}, nil
}
