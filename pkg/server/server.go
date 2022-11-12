package server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	_ "urlshortener/docs"
	"urlshortener/pkg/proto/apipb"
	"urlshortener/pkg/proto/cachepb"
	"urlshortener/pkg/proto/urlspb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	urlClient   urlspb.UrlsClient
	cacheClient cachepb.CacheClient
	router      *gin.Engine
}

var (
	port       = flag.Int("port", 8080, "The server port")
	urls_addr  = flag.String("urls_addr", "urls_service:8081", "url service address")
	cache_addr = flag.String("cache_addr", "cache_service:8082", "cache service address")
)

// @title       Go + Gin API
// @version     1.0
// @description This is a sample server.

// @contact.name  API Support
// @contact.url   http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url  https://opensource.org/licenses/MIT

// @host                    localhost:8080
// @BasePath                /
// @query.collection.format multi
// @schemes                 http

// @securityDefinitions.apiKey JWT_AUTH
// @in                         header
// @name                       Authorization
// @description:               'Authorization header: "Bearer [token]"'
func RunService() {

	flag.Parse()

	var err error
	urlConn, err := grpc.Dial(*urls_addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer urlConn.Close()

	cacheConn, err := grpc.Dial(*cache_addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer cacheConn.Close()

	server := &Server{
		urlClient:   urlspb.NewUrlsClient(urlConn),
		cacheClient: cachepb.NewCacheClient(cacheConn),
		router:      gin.Default(),
	}
	server.pingServices()
	server.registerEndpoints()

	log.Println("Starting server on port", *port)
	log.Printf("Swagger docs available at http://localhost:%d/docs/index.html\n", *port)

	if err := server.router.Run(fmt.Sprintf(":%d", *port)); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func (server *Server) registerEndpoints() {
	server.router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server.router.GET("/ping", server.Ping)
	server.router.Any("/:urlId", server.RedirectToUrl)
	server.router.GET("/url/:urlId", server.GetUrl)
	server.router.POST("/url", server.PostUrl)
}

func (server *Server) pingServices() {
	c, cancel := makeCtx()
	defer cancel()
	_, err := server.urlClient.Ping(c, &apipb.PingRequest{})
	if err != nil {
		log.Fatal("failed to ping url service: ", err)
	}
	c, cancel = makeCtx()
	defer cancel()
	_, err = server.cacheClient.Ping(c, &apipb.PingRequest{})
	if err != nil {
		log.Fatal("failed to ping cache service: ", err)
	}
}

func makeCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second)
}
