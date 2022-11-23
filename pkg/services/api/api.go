package api

import (
	"flag"
	"fmt"

	"urlshortener/pkg/logger"
	"urlshortener/pkg/proto/cachepb"

	"urlshortener/pkg/proto/urlspb"
	_ "urlshortener/pkg/services/api/api_docs"

	"github.com/gin-gonic/gin"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	port      = flag.Int("port", 8080, "the port to serve on")
	urlsAddr  = flag.String("urls_addr", "urls.service:8081", "url service address")
	cacheAddr = flag.String("cache_addr", "cache.service:8082", "cache service address")
	logLevel  = flag.String("log_level", "info", "the log level")
)

// @title                      URL Shortener API
// @version                    1.0
// @description                This is a URL shortener service.
// @host                       localhost:8080
// @BasePath                   /
// @query.collection.format    multi
// @schemes                    http
// @contact.name               API Support
// @contact.url                http://www.swagger.io/support
// @contact.email              support@swagger.io
// @license.name               MIT
// @license.url                https://opensource.org/licenses/MIT
// @securityDefinitions.apiKey JWT_AUTH
// @in                         header
// @name                       Authorization
// @description:               'Authorization header: "Bearer [token]"'

func NewService() *Service {
	flag.Parse()
	urlsConn := dialGrpc(*urlsAddr)
	cacheConn := dialGrpc(*cacheAddr)
	service := &Service{
		name:              "api.service",
		router:            gin.Default(),
		trustedProxies:    []string{"127.0.0.1"},
		urlsClient:        urlspb.NewUrlsServiceClient(urlsConn),
		urlsHealthClient:  healthpb.NewHealthClient(urlsConn),
		urlsServiceName:   "urls.service",
		cacheClient:       cachepb.NewCacheServiceClient(cacheConn),
		cacheHealthClient: healthpb.NewHealthClient(cacheConn),
		cacheServiceName:  "cache.service",
	}
	service.RegisterEndpoints()
	service.RegisterTrustedProxies()
	return service
}

func (s *Service) Run() {
	logger.SetLogLevel(logger.LevelFromString(*logLevel))
	// s.checkServices()
	s.WatchServices()
	logger.Info("Starting server on port:", *port)
	logger.Info(fmt.Sprintf("Swagger docs available at http://localhost:%d/docs/index.html", *port))
	if err := s.router.Run(fmt.Sprintf(":%d", *port)); err != nil {
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
