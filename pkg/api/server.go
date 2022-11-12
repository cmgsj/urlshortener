package api

import (
	"log"

	"github.com/mike9107/urlshortener/pkg/protobuf/cachepb"
	"github.com/mike9107/urlshortener/pkg/protobuf/urlspb"

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
}

func (server *apiServer) registerTrustedProxies() {
	err := server.router.SetTrustedProxies(server.trustedProxies)
	if err != nil {
		log.Fatalf("failed to set trusted proxies: %v", err)
	}
}

func (server *apiServer) registerEndpoints() {
	server.router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server.router.GET("/ping", server.Pong)
	server.router.Any("/:urlId", server.RedirectToUrl)
	server.router.GET("/url/:urlId", server.GetUrl)
	server.router.POST("/url", server.PostUrl)
}

func (server *apiServer) pingServices() {
	clients := []client{
		{name: "urls service", service: server.urlsClient, active: &server.urlsServiceActive},
		{name: "cache service", service: server.cacheClient, active: &server.cacheServiceActive},
	}
	for _, client := range clients {
		go makePingCall(client.service, client.name, client.active)
	}
}
