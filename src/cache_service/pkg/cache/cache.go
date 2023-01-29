package cache

import (
	"cache_service/pkg/cache/grpc_util/grpc_interceptor"
	"cache_service/pkg/cache/proto/cachepb"
	"fmt"

	"net"
	"os"

	"strconv"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	API_SVC_NAME       = os.Getenv("API_SVC_NAME")
	CACHE_SVC_NAME     = os.Getenv("CACHE_SVC_NAME")
	CACHE_SVC_PORT     = os.Getenv("CACHE_SVC_PORT")
	CACHE_SVC_EXP_TIME = os.Getenv("CACHE_SVC_EXP_TIME")
	REDIS_ADDR         = os.Getenv("REDIS_ADDR")
	REDIS_PASSWORD     = os.Getenv("REDIS_PASSWORD")
	REDIS_DB           = os.Getenv("REDIS_DB")
)

func NewService() *Service {
	service := &Service{
		healthServer: health.NewServer(),
		logger:       zap.Must(zap.NewDevelopment()),
	}
	cacheExpTime, err := strconv.Atoi(CACHE_SVC_EXP_TIME)
	if err != nil {
		service.logger.Fatal("failed to convert CACHE_SVC_EXP_TIME to int:", zap.Error(err))
	}
	service.cacheExpTime = time.Duration(cacheExpTime) * time.Second
	redisdb, err := strconv.Atoi(REDIS_DB)
	if err != nil {
		service.logger.Fatal("failed to convert REDIS_DB to int:", zap.Error(err))
	}
	service.initRedisDB(REDIS_ADDR, REDIS_PASSWORD, redisdb)
	service.healthServer.SetServingStatus(API_SVC_NAME, healthpb.HealthCheckResponse_SERVING)
	return service
}

func (s *Service) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", CACHE_SVC_PORT))
	if err != nil {
		s.logger.Fatal("failed to listen:", zap.Error(err))
	}

	loggerInterceptor := grpc_interceptor.NewLogger(s.logger)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggerInterceptor.Unary),
		grpc.StreamInterceptor(loggerInterceptor.Stream),
	)

	reflection.Register(grpcServer)
	healthpb.RegisterHealthServer(grpcServer, s.healthServer)
	cachepb.RegisterCacheServiceServer(grpcServer, s)

	s.logger.Info("Starting", zap.String("service", CACHE_SVC_NAME), zap.String("address", lis.Addr().String()))
	if err := grpcServer.Serve(lis); err != nil {
		s.logger.Fatal("failed to serve:", zap.Error(err))
	}
}
