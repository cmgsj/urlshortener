package api

import (
	"context"
	"sync"
	"time"

	"urlshortener/pkg/logger"
	"urlshortener/pkg/proto/cachepb"
	"urlshortener/pkg/proto/healthpb"
	"urlshortener/pkg/proto/urlspb"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type apiServer struct {
	urlsClient         urlspb.UrlsClient
	urlsServiceActive  bool
	cacheClient        cachepb.CacheClient
	cacheServiceActive bool
	router             *gin.Engine
	trustedProxies     []string
	mu                 sync.Mutex
}

func (server *apiServer) registerEndpoints() {
	server.router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server.router.GET("/ping", server.Pong)
	server.router.Any("/:urlId", server.RedirectToUrl)
	server.router.GET("/url/:urlId", server.GetUrl)
	server.router.POST("/url", server.PostUrl)
}

func (server *apiServer) registerTrustedProxies() {
	err := server.router.SetTrustedProxies(nil)
	if err != nil {
		logger.Fatal("failed to set trusted proxies:", err)
	}
}

func (server *apiServer) checkServices() {
	go makeCheckCall(server.urlsClient, "urls_service", &server.urlsServiceActive, &server.mu)
	go makeCheckCall(server.cacheClient, "cache_service", &server.cacheServiceActive, &server.mu)
}

func makeCheckCall(client healthServer, name string, active *bool, mu *sync.Mutex) {
	c, cancel := makeCtx()
	defer cancel()
	_, err := client.Check(c, &healthpb.HealthCheckRequest{Service: "api_service"})
	if err != nil {
		*active = false
		logger.Errorf("failed to ping %s: %v", name, err)
	} else {
		*active = true
		logger.Info(name, "is active")
	}
}

func makeCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second)
}
