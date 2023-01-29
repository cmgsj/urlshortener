package cache

import (
	"cache_service/pkg/cache/proto/cachepb"
	"context"

	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/status"
)

var (
	ErrKeyNotFound = status.Error(codes.NotFound, "key not found")
	ErrInternal    = status.Error(codes.Internal, "internal error")
)

type Service struct {
	cachepb.UnimplementedCacheServiceServer
	healthServer *health.Server
	logger       *zap.Logger
	rdb          *redis.Client
	cacheExpTime time.Duration
}

func (s *Service) Get(ctx context.Context, req *cachepb.GetRequest) (*cachepb.GetResponse, error) {
	value, err := s.rdb.Get(ctx, req.GetKey()).Result()
	if err != nil {
		return nil, ErrKeyNotFound
	}
	return &cachepb.GetResponse{Value: value}, nil
}

func (s *Service) Set(ctx context.Context, req *cachepb.SetRequest) (*cachepb.NoContent, error) {
	err := s.rdb.Set(ctx, req.GetKey(), req.GetValue(), s.cacheExpTime).Err()
	if err != nil {
		return nil, ErrInternal
	}
	return &cachepb.NoContent{}, nil
}
