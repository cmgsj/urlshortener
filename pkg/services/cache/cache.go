package cache

import (
	"flag"
	"fmt"
	"net"
	"time"

	"urlshortener/pkg/interceptor"
	"urlshortener/pkg/logger"
	"urlshortener/pkg/proto/cachepb"

	"google.golang.org/grpc"
)

var (
	port          = flag.Int("port", 8082, "the port to serve on")
	redisAddr     = flag.String("redis_addr", "redis_cache:6379", "the redis address")
	redisPassword = flag.String("redis_password", "", "the redis password")
	redisDb       = flag.Int("redis_db", 0, "the redis db")
	cacheExpTime  = flag.Duration("cache_exp_time", time.Hour, "the cache expiry time")
)

func NewService() *cacheServer {
	flag.Parse()

	server := &cacheServer{
		rdb:             initRedisDB(*redisAddr, *redisPassword, *redisDb),
		cacheExpiryTime: *cacheExpTime,
	}
	return server
}

func (server *cacheServer) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		logger.Fatal("failed to listen:", err)
	}

	loggerInterceptor := interceptor.NewLoggerInterceptor()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggerInterceptor.Unary),
		grpc.StreamInterceptor(loggerInterceptor.Stream))

	cachepb.RegisterCacheServer(grpcServer, server)

	logger.Info("Starting cache_service at:", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("failed to serve:", err)
	}
}
