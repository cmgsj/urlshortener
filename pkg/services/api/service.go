package api

import (
	"context"
	"time"

	"urlshortener/pkg/proto/cachepb"

	"urlshortener/pkg/proto/urlspb"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Service struct {
	name              string
	trustedProxies    []string
	router            *gin.Engine
	logger            *zap.Logger
	urlsClient        urlspb.UrlsServiceClient
	urlsHealthClient  healthpb.HealthClient
	urlsServiceOk     bool
	urlsServiceName   string
	cacheClient       cachepb.CacheServiceClient
	cacheHealthClient healthpb.HealthClient
	cacheServiceOk    bool
	cacheServiceName  string
}

func (s *Service) RegisterEndpoints() {
	s.router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	s.router.GET("/ping", s.Pong)
	s.router.Any("/:urlId", s.RedirectToUrl)
	s.router.GET("/url/:urlId", s.GetUrl)
	s.router.POST("/url", s.PostUrl)
}

func (s *Service) RegisterTrustedProxies() {
	err := s.router.SetTrustedProxies(nil)
	if err != nil {
		s.logger.Fatal("failed to set trusted proxies", zap.Error(err))
	}
}

func (s *Service) CheckServices() {
	go s.checkService(s.name, s.urlsHealthClient, s.urlsServiceName, &s.urlsServiceOk)
	go s.checkService(s.name, s.cacheHealthClient, s.cacheServiceName, &s.cacheServiceOk)
}

func (s *Service) WatchServices() {
	go s.watchService(s.name, s.urlsHealthClient, s.urlsServiceName, &s.urlsServiceOk)
	go s.watchService(s.name, s.cacheHealthClient, s.cacheServiceName, &s.cacheServiceOk)
}

func (s *Service) checkService(name string, client healthpb.HealthClient, serviceName string, active *bool) {
	ctx, cancel := makeUnaryCtx()
	defer cancel()
	_, err := client.Check(ctx, &healthpb.HealthCheckRequest{Service: name})
	if err != nil {
		*active = false
		s.logger.Error("failed to check:", zap.String("service", serviceName), zap.Error(err))
	} else {
		*active = true
		s.logger.Info("service is active", zap.String("service", serviceName))
	}
}

func (s *Service) watchService(name string, client healthpb.HealthClient, serviceName string, active *bool) {
	ctx, cancel := makeStreamCtx()
	defer cancel()
	stream, err := client.Watch(ctx, &healthpb.HealthCheckRequest{Service: name})
	if err != nil {
		*active = false
		s.logger.Error("failed to watch:", zap.String("service", serviceName), zap.Error(err))
		return
	}
	for {
		res, err := stream.Recv()
		if err != nil {
			*active = false
			s.logger.Error("failed to receive:", zap.String("service", serviceName), zap.Error(err))
			return
		}
		if res.GetStatus() == healthpb.HealthCheckResponse_SERVING {
			*active = true
			s.logger.Info("service is active", zap.String("service", serviceName))
		} else {
			*active = false
			s.logger.Info("service is not active", zap.String("service", serviceName))
		}
	}
}

func makeUnaryCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second)

}

func makeStreamCtx() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}
