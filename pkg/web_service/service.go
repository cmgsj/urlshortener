package web_service

import (
	"time"

	"github.com/cmgsj/url-shortener/pkg/proto/urlpb"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sony/gobreaker"
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

type Service struct {
	Addr           string
	TrustedProxies []string
	Router         *gin.Engine
	Logger         *zap.Logger
	UrlClient      urlpb.UrlServiceClient
	UrlServiceCb   *gobreaker.CircuitBreaker
	RedisDb        *redis.Client
	CacheTimeout   time.Duration
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

func MakeUrlServiceCb(logger *zap.Logger) *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(
		gobreaker.Settings{
			Name:        "url-svc-circuit-breaker",
			MaxRequests: 3,
			Timeout:     3 * time.Second,
			Interval:    1 * time.Second,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				return counts.ConsecutiveFailures > 3
			},
			OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
				logger.Info("circuit breaker state change", zap.String("cb-name", name), zap.String("from", from.String()), zap.String("to", to.String()))
			},
		},
	)
}
