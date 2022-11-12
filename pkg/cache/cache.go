package cache

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"urlshortener/pkg/grpc/interceptor"
	"urlshortener/pkg/protobuf/cachepb"

	"google.golang.org/grpc"
)

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
