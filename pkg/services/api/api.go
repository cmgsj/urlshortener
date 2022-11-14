package api

import (
	"flag"
	"fmt"

	"urlshortener/pkg/logger"
	"urlshortener/pkg/proto/cachepb"
	"urlshortener/pkg/proto/healthpb"
	"urlshortener/pkg/proto/urlspb"
	_ "urlshortener/pkg/services/api/api_docs"

	"github.com/gin-gonic/gin"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port      = flag.Int("port", 8080, "the port to serve on")
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

func NewService() *Service {
	flag.Parse()
	urlsConn := dialGrpc(*urlsAddr)
	cacheConn := dialGrpc(*cacheAddr)
	service := &Service{
		urlsClient:        urlspb.NewUrlsServiceClient(urlsConn),
		urlsHealthClient:  healthpb.NewHealthServiceClient(urlsConn),
		cacheClient:       cachepb.NewCacheServiceClient(cacheConn),
		cacheHealthClient: healthpb.NewHealthServiceClient(cacheConn),
		router:            gin.Default(),
		trustedProxies:    []string{"127.0.0.1"},
	}
	service.registerEndpoints()
	service.registerTrustedProxies()
	return service
}

func (service *Service) Run() {
	logger.SetLogLevel(logger.LevelFromString(*logLevel))
	service.checkServices()
	service.watchServices()
	// scheduler.SchedulePeriodicTask(service.checkServices, time.Minute)
	logger.Info("Starting server on port:", *port)
	logger.Info(fmt.Sprintf("Swagger docs available at http://localhost:%d/docs/index.html", *port))
	if err := service.router.Run(fmt.Sprintf(":%d", *port)); err != nil {
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
