package api

import (
	"api_service/pkg/proto/cachepb"
	"api_service/pkg/proto/urlspb"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Service struct {
	name              string
	addr              string
	trustedProxies    []string
	router            *gin.Engine
	logger            *zap.Logger
	urlsClient        urlspb.UrlsServiceClient
	urlsHealthClient  healthpb.HealthClient
	urlsServiceOk     atomic.Bool
	urlsServiceName   string
	cacheClient       cachepb.CacheServiceClient
	cacheHealthClient healthpb.HealthClient
	cacheServiceOk    atomic.Bool
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
