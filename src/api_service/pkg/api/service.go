package api

import (
	"api_service/pkg/proto/urlspb"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Service struct {
	name             string
	addr             string
	trustedProxies   []string
	router           *gin.Engine
	logger           *zap.Logger
	urlsClient       urlspb.UrlsServiceClient
	urlsHealthClient healthpb.HealthClient
	urlsServiceOk    atomic.Bool
	urlsServiceName  string
	redisDb          *redis.Client
	cacheExpTime     time.Duration
}

func (s *Service) RegisterEndpoints() {
	s.router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	s.router.GET("/ping", s.Pong)
	s.router.GET("/url/:urlId", s.GetUrl)
	s.router.POST("/url", s.PostUrl)
	s.router.Any("/:urlId", s.RedirectToUrl)
}

func (s *Service) RegisterTrustedProxies() {
	err := s.router.SetTrustedProxies(nil)
	if err != nil {
		s.logger.Fatal("failed to set trusted proxies", zap.Error(err))
	}
}
