package cache

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"urlshortener/pkg/grpc/interceptor"
	"urlshortener/pkg/proto/apipb"
	"urlshortener/pkg/proto/cachepb"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
)

type cacheServer struct {
	cachepb.UnimplementedCacheServer
	rdb             *redis.Client
	cacheExpiryTime time.Duration
}

var (
	port           = flag.Int("port", 8082, "the port to serve on")
	redis_addr     = flag.String("redis_addr", "redis_cache:6379", "the redis address")
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

	cachepb.RegisterCacheServer(grpcServer, server)

	log.Printf("Starting cache_service at %v\n", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func (server *cacheServer) GetUrl(ctx context.Context, req *cachepb.GetUrlRequest) (*cachepb.GetUrlResponse, error) {
	redirectUrl, err := server.rdb.Get(ctx, req.GetUrlId()).Result()
	if err != nil {
		return nil, err
	}
	return &cachepb.GetUrlResponse{RedirectUrl: redirectUrl}, nil
}

func (server *cacheServer) SetUrl(ctx context.Context, req *cachepb.SetUrlRequest) (*apipb.NoContent, error) {
	err := server.rdb.Set(ctx, req.GetUrlId(), req.GetRedirectUrl(), server.cacheExpiryTime).Err()
	if err != nil {
		return nil, err
	}
	return &apipb.NoContent{}, nil
}

func (server *cacheServer) Ping(ctx context.Context, req *apipb.PingRequest) (*apipb.PingResponse, error) {
	return &apipb.PingResponse{Message: "pong"}, nil
}

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
