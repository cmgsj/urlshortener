package websvc

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func (s *Service) InitRedisDB(redisAddr string, redisPassword string, redisDb int) {
	db := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDb,
	})
	_, err := db.Ping(context.Background()).Result()
	if err != nil {
		s.Logger.Fatal("failed to connect to redis:", zap.Error(err))
	}
	s.RedisDb = db
}

func (s *Service) getFromCache(key string) (string, error) {
	ctx := context.Background()
	value, err := s.RedisDb.Get(ctx, key).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			s.Logger.Error("failed to get from cache", zap.Error(err))
		}
		return "", err
	}
	return value, nil
}

func (s *Service) putInCache(key string, value string) {
	ctx := context.Background()
	if err := s.RedisDb.Set(ctx, key, value, s.CacheTimeout).Err(); err != nil {
		s.Logger.Error("failed to set cache", zap.Error(err))
	}
}
