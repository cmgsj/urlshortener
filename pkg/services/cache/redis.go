package cache

import (
	"context"
	"urlshortener/pkg/logger"

	"github.com/go-redis/redis/v8"
)

func initRedisDB(redisAddr string, redisPassword string, redisDb int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDb,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		logger.Fatal("failed to connect to redis:", err)
	}
	return rdb
}
