package server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"time"

	_ "urlshortener/docs"
	"urlshortener/pkg/api"
	"urlshortener/pkg/cache"
	"urlshortener/pkg/urls"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	urlClient   urls.UrlsClient
	cacheClient cache.CacheClient
	router      *gin.Engine
}

var (
	port       = flag.Int("port", 8080, "The server port")
	url_addr   = flag.String("url_addr", "localhost:50051", "url service address")
	cache_addr = flag.String("cache_addr", "localhost:50052", "cache service address")
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

	fmt.Println("Starting http server")
	flag.Parse()

	var err error
	urlConn, err := grpc.Dial(*url_addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
		urlClient:   urls.NewUrlsClient(urlConn),
		cacheClient: cache.NewCacheClient(cacheConn),
		router:      gin.Default(),
	}
	server.pingServices()
	server.registerEndpoints()

	exec.Command("open", fmt.Sprintf("http://localhost:%d/docs/index.html", *port)).Start()

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
	_, err := server.urlClient.Ping(c, &api.PingRequest{})
	if err != nil {
		log.Fatal(err)
	}
	c, cancel = makeCtx()
	defer cancel()
	_, err = server.cacheClient.Ping(c, &api.PingRequest{})
	if err != nil {
		log.Fatal(err)
	}
}

func makeCtx() (context.Context, context.CancelFunc) {
	// ctx := metadata.AppendToOutgoingContext(context.Background(), "token", "xxxxx")
	return context.WithTimeout(context.Background(), time.Second)
}
