package api

import (
	"flag"
	"fmt"
	"time"

	_ "urlshortener/docs"
	"urlshortener/pkg/logger"
	"urlshortener/pkg/proto/cachepb"
	"urlshortener/pkg/proto/urlspb"
	"urlshortener/pkg/scheduler"

	"github.com/gin-gonic/gin"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port      = flag.Int("port", 8080, "The server port")
	urlsAddr  = flag.String("urls_addr", "urls_service:8081", "url service address")
	cacheAddr = flag.String("cache_addr", "cache_service:8082", "cache service address")
	logLevel  = flag.String("log_level", "info", "the log level")
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
	server := &apiServer{
		urlsClient:     urlspb.NewUrlsClient(dialGrpc(*urlsAddr)),
		cacheClient:    cachepb.NewCacheClient(dialGrpc(*cacheAddr)),
		router:         gin.Default(),
		trustedProxies: []string{"127.0.0.1"},
	}
	server.registerEndpoints()
	server.registerTrustedProxies()
	return server
}

func (server *apiServer) Run() {
	logger.SetLogLevel(logger.LevelFromString(*logLevel))

	go server.checkServices()
	scheduler.SchedulePeriodicTask(server.checkServices, time.Minute)

	logger.Info("Starting server on port:", *port)
	logger.Infof("Swagger docs available at http://localhost:%d/docs/index.html", *port)
	if err := server.router.Run(fmt.Sprintf(":%d", *port)); err != nil {
		logger.Fatal("failed to serve:", err)
	}
}

func dialGrpc(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("failed to connect:", err)
	}
	return conn
}
