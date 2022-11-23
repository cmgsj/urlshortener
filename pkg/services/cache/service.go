package cache

import (
	"context"
	"time"
	"urlshortener/pkg/proto/cachepb"

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
	rdb             *redis.Client
	cacheExpiryTime time.Duration
}

func (s *Service) GetUrl(ctx context.Context, req *cachepb.GetRequest) (*cachepb.GetResponse, error) {
	redirectUrl, err := s.rdb.Get(ctx, req.GetKey()).Result()
	if err != nil {
		return nil, UrlNotFoundError
	}
	return &cachepb.GetResponse{Value: redirectUrl}, nil
}

func (s *Service) SetUrl(ctx context.Context, req *cachepb.SetRequest) (*cachepb.NoContent, error) {
	err := s.rdb.Set(ctx, req.GetKey(), req.GetValue(), s.cacheExpiryTime).Err()
	if err != nil {
		return nil, InternalServerError
	}
	return &cachepb.NoContent{}, nil
}
