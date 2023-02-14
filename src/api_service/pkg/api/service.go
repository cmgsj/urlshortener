package api

import (
	"proto/pkg/urlpb"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// @title                      URL Shortener API
// @version                    1.0
// @description                This is a URL shortener service.
// @host                       localhost:8080
// @BasePath                   /
// @query.collection.format    multi
// @schemes                    http
// @contact.name               API Support
// @contact.url                http://www.swagger.io/support
// @contact.email              support@swagger.io
// @license.name               MIT
// @license.url                https://opensource.org/licenses/MIT
// @securityDefinitions.apiKey JWT_AUTH
// @in                         header
// @name                       Authorization
// @description:               'Authorization header: "Bearer [token]"'

type Service struct {
	Name            string
	Addr            string
	TrustedProxies  []string
	Router          *gin.Engine
	Logger          *zap.Logger
	UrlClient       urlpb.UrlServiceClient
	UrlHealthClient healthpb.HealthClient
	UrlServiceOk    atomic.Bool
	UrlServiceName  string
	RedisDb         *redis.Client
	CacheTimeout    time.Duration
}

func (s *Service) RegisterEndpoints() {
	s.Router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	s.Router.GET("/ping", s.Pong)
	s.Router.GET("/url/:urlId", s.GetUrl)
	s.Router.POST("/url", s.PostUrl)
	s.Router.Any("/:urlId", s.RedirectToUrl)
}

func (s *Service) RegisterTrustedProxies() {
	err := s.Router.SetTrustedProxies(nil)
	if err != nil {
		s.Logger.Fatal("failed to set trusted proxies", zap.Error(err))
	}
}
