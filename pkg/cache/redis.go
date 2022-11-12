package cache

import (
	"context"
	"log"

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
		log.Fatalf("failed to connect to redis: %v", err)
	}
	return rdb
}
