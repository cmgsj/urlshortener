package api

import (
	"flag"
	"fmt"

	"urlshortener/pkg/grpcutil"
	"urlshortener/pkg/proto/cachepb"

	"urlshortener/pkg/proto/urlspb"
	_ "urlshortener/pkg/services/api/api_docs"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	port      = flag.Int("port", 8080, "the port to serve on")
	urlsAddr  = flag.String("urls_addr", "urls.service:8081", "url service address")
	cacheAddr = flag.String("cache_addr", "cache.service:8082", "cache service address")
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
	urlsConn := grpcutil.DialGrpc(*urlsAddr)
	cacheConn := grpcutil.DialGrpc(*cacheAddr)
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())
	service := &Service{
		name:              "api.service",
		trustedProxies:    []string{"127.0.0.1"},
		router:            router,
		logger:            zap.Must(zap.NewDevelopment()),
		urlsClient:        urlspb.NewUrlsServiceClient(urlsConn),
		urlsHealthClient:  healthpb.NewHealthClient(urlsConn),
		urlsServiceName:   "urls.service",
		urlsServiceOk:     false,
		cacheClient:       cachepb.NewCacheServiceClient(cacheConn),
		cacheHealthClient: healthpb.NewHealthClient(cacheConn),
		cacheServiceName:  "cache.service",
		cacheServiceOk:    false,
	}
	service.RegisterEndpoints()
	service.RegisterTrustedProxies()
	return service
}

func (s *Service) Run() {
	s.CheckServices()
	s.WatchServices()
	s.logger.Info("Starting server on:", zap.Int("port", *port))
	s.logger.Info(fmt.Sprintf("Swagger docs available at http://localhost:%d/docs/index.html", *port))
	if err := s.router.Run(fmt.Sprintf(":%d", *port)); err != nil {
		s.logger.Fatal("Failed to start server:", zap.Error(err))
	}
}
