package cache

import (
	"context"

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
	s.rdb = rdb
}
