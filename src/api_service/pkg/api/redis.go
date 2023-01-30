package api

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func (s *Service) initRedisDB(redisAddr string, redisPassword string, redisDb int) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDb,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		s.logger.Fatal("failed to connect to redis:", zap.Error(err))
	}
	s.redisDb = rdb
}

func (s *Service) getFromCache(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	value, err := s.redisDb.Get(ctx, key).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			s.logger.Error("failed to get from cache", zap.Error(err))
		}
		return "", err
	}
	return value, nil
}

func (s *Service) putInCache(key string, value string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := s.redisDb.Set(ctx, key, value, s.cacheExpTime).Err(); err != nil {
		s.logger.Error("failed to set cache", zap.Error(err))
	}
}
