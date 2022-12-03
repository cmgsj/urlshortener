package api

import (
	"urlshortener/pkg/grpcutil"
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
	go grpcutil.CheckService(s.name, s.logger, s.urlsHealthClient, s.urlsServiceName, &s.urlsServiceOk)
	go grpcutil.CheckService(s.name, s.logger, s.cacheHealthClient, s.cacheServiceName, &s.cacheServiceOk)
}

func (s *Service) WatchServices() {
	go grpcutil.WatchService(s.name, s.logger, s.urlsHealthClient, s.urlsServiceName, &s.urlsServiceOk)
	go grpcutil.WatchService(s.name, s.logger, s.cacheHealthClient, s.cacheServiceName, &s.cacheServiceOk)
}
