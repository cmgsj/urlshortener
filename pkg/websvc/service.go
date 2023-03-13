package websvc

import (
	"time"

	"github.com/cmgsj/urlshortener/pkg/proto/urlpb"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// @title                   URL Shortener API
// @version                 1.0
// @description             This is a URL shortener service.
// @host                    localhost:8080
// @BasePath                /
// @query.collection.format multi
// @schemes                 http

type (
	Service struct {
		Addr           string
		TrustedProxies []string
		Router         *gin.Engine
		Logger         *zap.Logger
		UrlClient      urlpb.UrlServiceClient
		RedisDb        *redis.Client
		CacheTimeout   time.Duration
	}

	Options struct {
		Addr           string
		TrustedProxies []string
		Logger         *zap.Logger
		UrlClient      urlpb.UrlServiceClient
		CacheTimeout   time.Duration
	}
)

func New(opt Options) *Service {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())
	return &Service{
		Addr:           opt.Addr,
		TrustedProxies: opt.TrustedProxies,
		Router:         router,
		Logger:         opt.Logger,
		UrlClient:      opt.UrlClient,
		CacheTimeout:   opt.CacheTimeout,
	}
}

func (s *Service) RegisterEndpoints() {
	s.Router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	s.Router.GET("/ping", s.Pong)
	s.Router.GET("/url/:urlId", s.GetUrl)
	s.Router.POST("/url", s.PostUrl)
	s.Router.Any("/:urlId", s.RedirectToUrl)
}

func (s *Service) RegisterTrustedProxies() {
	err := s.Router.SetTrustedProxies(s.TrustedProxies)
	if err != nil {
		s.Logger.Fatal("failed to set trusted proxies", zap.Error(err))
	}
}
