package api

import (
	"flag"
	"fmt"
	"log"
	"time"

	_ "urlshortener/docs"
	"urlshortener/pkg/protobuf/cachepb"
	"urlshortener/pkg/protobuf/urlspb"
	"urlshortener/pkg/scheduler"

	"github.com/gin-gonic/gin"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port      = flag.Int("port", 8080, "The server port")
	urlsAddr  = flag.String("urls_addr", "urls_service:8081", "url service address")
	cacheAddr = flag.String("cache_addr", "cache_service:8082", "cache service address")
)

// @title                   URL Shortener API
// @version                 1.0
// @description             This is a URL shortener service.
// @host                    localhost:8080
// @BasePath                /
// @query.collection.format multi
// @schemes                 http
// @contact.name            API Support
// @contact.url             http://www.swagger.io/support
// @contact.email           support@swagger.io
// @license.name            MIT
// @license.url             https://opensource.org/licenses/MIT

func NewService() *apiServer {
	flag.Parse()

	urlConn, err := grpc.Dial(*urlsAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	cacheConn, err := grpc.Dial(*cacheAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	server := &apiServer{
		urlsClient:     urlspb.NewUrlsClient(urlConn),
		cacheClient:    cachepb.NewCacheClient(cacheConn),
		router:         gin.Default(),
		trustedProxies: []string{"127.0.0.1"},
	}
	server.registerEndpoints()
	server.registerTrustedProxies()
	return server
}

func (server *apiServer) Run() {
	go server.pingServices()
	scheduler.SchedulePeriodicTask(server.pingServices, time.Minute)

	log.Println("Starting server on port", *port)
	log.Printf("Swagger docs available at http://localhost:%d/docs/index.html\n", *port)

	if err := server.router.Run(fmt.Sprintf(":%d", *port)); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
