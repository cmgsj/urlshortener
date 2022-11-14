package api

import (
	"context"
	"time"

	"urlshortener/pkg/logger"
	"urlshortener/pkg/proto/cachepb"
	"urlshortener/pkg/proto/healthpb"
	"urlshortener/pkg/proto/urlspb"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Service struct {
	urlsClient        urlspb.UrlsServiceClient
	urlsHealthClient  healthpb.HealthServiceClient
	urlsServiceOk     bool
	cacheClient       cachepb.CacheServiceClient
	cacheHealthClient healthpb.HealthServiceClient
	cacheServiceOk    bool
	router            *gin.Engine
	trustedProxies    []string
}

func (service *Service) registerEndpoints() {
	service.router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	service.router.GET("/ping", service.Pong)
	service.router.Any("/:urlId", service.RedirectToUrl)
	service.router.GET("/url/:urlId", service.GetUrl)
	service.router.POST("/url", service.PostUrl)
}

func (service *Service) registerTrustedProxies() {
	err := service.router.SetTrustedProxies(nil)
	if err != nil {
		logger.Fatal("failed to set trusted proxies:", err)
	}
}

func (service *Service) checkServices() {
	go checkService(service.urlsHealthClient, "urls_service", &service.urlsServiceOk)
	go checkService(service.urlsHealthClient, "cache_service", &service.cacheServiceOk)
}

func (service *Service) watchServices() {
	go watchService(service.urlsHealthClient, "urls_service", &service.urlsServiceOk)
	go watchService(service.cacheHealthClient, "cache_service", &service.cacheServiceOk)
}

func checkService(client healthpb.HealthServiceClient, name string, active *bool) {
	c, cancel := makeCtx()
	defer cancel()
	_, err := client.Check(c, &healthpb.HealthCheckRequest{Service: "api_service"})
	if err != nil {
		*active = false
		logger.Error("failed to ping:", name, err)
	} else {
		*active = true
		logger.Info(name, "is active")
	}
}

func watchService(client healthpb.HealthServiceClient, name string, active *bool) {
	c, cancel := makeCtx()
	defer cancel()
	stream, err := client.Watch(c, &healthpb.HealthCheckRequest{Service: "api_service"})
	if err != nil {
		logger.Error("failed to watch:", name, err)
		return
	}
	for {
		res, err := stream.Recv()
		if err != nil {
			logger.Error("failed to receive:", name, err)
			return
		}
		if res.Status == healthpb.HealthCheckResponse_SERVING {
			*active = true
			logger.Info(name, "is active")
		} else {
			*active = false
			logger.Error(name, "is inactive, status:", res.Status)
		}
	}
}

func makeCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second)
}
