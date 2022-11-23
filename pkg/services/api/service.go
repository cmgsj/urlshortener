package api

import (
	"context"
	"time"

	"urlshortener/pkg/logger"
	"urlshortener/pkg/proto/cachepb"

	"urlshortener/pkg/proto/urlspb"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Service struct {
	name              string
	router            *gin.Engine
	trustedProxies    []string
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
		logger.Fatal("failed to set trusted proxies:", err)
	}
}

func (s *Service) CheckServices() {
	go checkService(s.name, s.urlsHealthClient, s.urlsServiceName, &s.urlsServiceOk)
	go checkService(s.name, s.cacheHealthClient, s.cacheServiceName, &s.cacheServiceOk)
}

func (s *Service) WatchServices() {
	go watchService(s.name, s.urlsHealthClient, s.urlsServiceName, &s.urlsServiceOk)
	go watchService(s.name, s.cacheHealthClient, s.cacheServiceName, &s.cacheServiceOk)
}

func checkService(name string, client healthpb.HealthClient, serviceName string, active *bool) {
	ctx, cancel := makeUnaryCtx()
	defer cancel()
	_, err := client.Check(ctx, &healthpb.HealthCheckRequest{Service: name})
	if err != nil {
		*active = false
		logger.Error("failed to ping:", serviceName, err)
	} else {
		*active = true
		logger.Info(serviceName, "is active")
	}
}

func watchService(name string, client healthpb.HealthClient, serviceName string, active *bool) {
	ctx, cancel := makeStreamCtx()
	defer cancel()
	stream, err := client.Watch(ctx, &healthpb.HealthCheckRequest{Service: name})
	if err != nil {
		*active = false
		logger.Error("failed to watch:", serviceName, err)
		return
	}
	for {
		res, err := stream.Recv()
		if err != nil {
			*active = false
			logger.Error("failed to receive:", serviceName, err)
			return
		}
		if res.GetStatus() == healthpb.HealthCheckResponse_SERVING {
			*active = true
			logger.Info(serviceName, "is active")
		} else {
			*active = false
			logger.Error(serviceName, "is inactive, status:", res.GetStatus())
		}
	}
}

func makeUnaryCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second)
}

func makeStreamCtx() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}
