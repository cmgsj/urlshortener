package cache

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"urlshortener/pkg/api"
	"urlshortener/pkg/interceptor"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
)

type cacheServer struct {
	UnimplementedCacheServer
	rdb             *redis.Client
	cacheExpiryTime time.Duration
}

var (
	port           = flag.Int("cache_port", 50052, "the port to serve on")
	redis_addr     = flag.String("redis_addr", "localhost:6379", "the redis address")
	redis_password = flag.String("redis_password", "", "the redis password")
	redis_db       = flag.Int("redis_db", 0, "the redis db")
	cache_exp_time = flag.Duration("cache_exp_time", time.Hour, "the cache expiry time")
)

func RunService() {

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcInterceptor := interceptor.NewGrpcInterceptor()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcInterceptor.UnaryLogger),
		grpc.StreamInterceptor(grpcInterceptor.StreamLogger))

	server := &cacheServer{
		rdb:             initRedisDB(*redis_addr, *redis_password, *redis_db),
		cacheExpiryTime: *cache_exp_time,
	}

	RegisterCacheServer(grpcServer, server)

	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func (server *cacheServer) GetUrl(ctx context.Context, req *GetUrlRequest) (*GetUrlResponse, error) {
	redirectUrl, err := server.rdb.Get(ctx, req.GetUrlId()).Result()
	if err != nil {
		return nil, err
	}
	return &GetUrlResponse{RedirectUrl: redirectUrl}, nil
}

func (server *cacheServer) SetUrl(ctx context.Context, req *SetUrlRequest) (*api.NoContent, error) {
	err := server.rdb.Set(ctx, req.GetUrlId(), req.GetRedirectUrl(), server.cacheExpiryTime).Err()
	if err != nil {
		return nil, err
	}
	return &api.NoContent{}, nil
}

func (server *cacheServer) Ping(ctx context.Context, req *api.PingRequest) (*api.PingResponse, error) {
	return &api.PingResponse{Message: "pong"}, nil
}

func initRedisDB(redisAddr string, redisPassword string, redisDb int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDb,
	})
	ctx := context.Background()
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	log.Printf("redis ping response: %v", pong)
	return rdb
}
