package cache

import (
	"flag"
	"fmt"
	"net"
	"time"
	"urlshortener/pkg/interceptor"
	"urlshortener/pkg/proto/cachepb"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
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
		healthServer: health.NewServer(),
		logger:       zap.Must(zap.NewDevelopment()),
		cacheExpTime: *cacheExpTime,
	}
	service.initRedisDB(*redisAddr, *redisPassword, *redisDb)
	service.healthServer.SetServingStatus("api.service", healthpb.HealthCheckResponse_SERVING)
	return service
}

func (s *Service) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		s.logger.Fatal("failed to listen:", zap.Error(err))
	}

	loggerInterceptor := interceptor.NewLogger()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggerInterceptor.Unary),
		grpc.StreamInterceptor(loggerInterceptor.Stream),
	)
	reflection.Register(grpcServer)

	healthpb.RegisterHealthServer(grpcServer, s.healthServer)
	cachepb.RegisterCacheServiceServer(grpcServer, s)

	s.logger.Info("Starting cache.service at:", zap.String("address", lis.Addr().String()))
	if err := grpcServer.Serve(lis); err != nil {
		s.logger.Fatal("failed to serve:", zap.Error(err))
	}
}
