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
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	port          = flag.Int("port", 8082, "the port to serve on")
	redisAddr     = flag.String("redis_addr", "redis_cache:6379", "the redis address")
	redisPassword = flag.String("redis_password", "", "the redis password")
	redisDb       = flag.Int("redis_db", 0, "the redis db")
	cacheExpTime  = flag.Duration("cache_exp_time", time.Hour, "the cache expiry time")
)

func NewService() *Service {
	flag.Parse()
	service := &Service{
		rdb:             initRedisDB(*redisAddr, *redisPassword, *redisDb),
		cacheExpiryTime: *cacheExpTime,
	}
	return service
}

func (s *Service) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		logger.Fatal("failed to listen:", err)
	}

	loggerInterceptor := interceptor.NewLoggerInterceptor()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggerInterceptor.Unary),
		grpc.StreamInterceptor(loggerInterceptor.Stream),
	)

	healthServer := health.NewServer()

	healthServer.SetServingStatus("api.service", healthpb.HealthCheckResponse_SERVING)

	healthpb.RegisterHealthServer(grpcServer, healthServer)
	cachepb.RegisterCacheServiceServer(grpcServer, s)

	logger.Info("Starting cache.service at:", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("failed to serve:", err)
	}
}
