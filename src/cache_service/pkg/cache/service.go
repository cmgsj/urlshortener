package cache

import (
	"cache_service/pkg/proto/cachepb"
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

func (s *Service) GetItem(ctx context.Context, req *cachepb.GetItemRequest) (*cachepb.GetItemResponse, error) {
	value, err := s.rdb.Get(ctx, req.GetKey()).Result()
	if err != nil {
		return nil, ErrKeyNotFound
	}
	return &cachepb.GetItemResponse{Item: &cachepb.Item{Key: req.GetKey(), Value: value}}, nil
}

func (s *Service) PutItem(ctx context.Context, req *cachepb.PutItemRequest) (*cachepb.PutItemResponse, error) {
	err := s.rdb.Set(ctx, req.GetItem().GetKey(), req.GetItem().GetValue(), s.cacheExpTime).Err()
	if err != nil {
		return nil, ErrInternal
	}
	return &cachepb.PutItemResponse{}, nil
}
