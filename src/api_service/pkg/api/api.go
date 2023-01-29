package api

import (
	_ "api_service/pkg/docs"
	"grpc_util/pkg/grpc_dial"
	"grpc_util/pkg/grpc_health"
	"proto/pkg/cachepb"
	"proto/pkg/urlspb"

	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	API_SVC_PORT   = os.Getenv("API_SVC_PORT")
	API_SVC_ADDR   = os.Getenv("API_SVC_ADDR")
	API_SVC_NAME   = os.Getenv("API_SVC_NAME")
	CACHE_SVC_ADDR = os.Getenv("CACHE_SVC_ADDR")
	CACHE_SVC_NAME = os.Getenv("CACHE_SVC_NAME")
	URLS_SVC_ADDR  = os.Getenv("URLS_SVC_ADDR")
	URLS_SVC_NAME  = os.Getenv("URLS_SVC_NAME")
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
	urlsConn := grpc_dial.DialGrpc(URLS_SVC_ADDR)
	cacheConn := grpc_dial.DialGrpc(CACHE_SVC_ADDR)
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())
	service := &Service{
		name:              API_SVC_NAME,
		addr:              API_SVC_ADDR,
		trustedProxies:    []string{"127.0.0.1"},
		router:            router,
		logger:            zap.Must(zap.NewDevelopment()),
		urlsClient:        urlspb.NewUrlsServiceClient(urlsConn),
		urlsHealthClient:  healthpb.NewHealthClient(urlsConn),
		urlsServiceName:   URLS_SVC_NAME,
		cacheClient:       cachepb.NewCacheServiceClient(cacheConn),
		cacheHealthClient: healthpb.NewHealthClient(cacheConn),
		cacheServiceName:  CACHE_SVC_NAME,
	}
	service.RegisterEndpoints()
	service.RegisterTrustedProxies()
	return service
}

func (s *Service) Run() {
	go grpc_health.WatchServices(s.name, s.logger, []*grpc_health.HealthClient{
		{HealthClient: s.urlsHealthClient, Active: &s.urlsServiceOk, Name: s.urlsServiceName},
		{HealthClient: s.cacheHealthClient, Active: &s.cacheServiceOk, Name: s.cacheServiceName},
	})
	s.addr = fmt.Sprintf("localhost:%s", API_SVC_PORT)
	s.logger.Info("Starting", zap.String("service", API_SVC_NAME), zap.String("port", API_SVC_PORT))
	s.logger.Info(fmt.Sprintf("Swagger docs available at http://%s/docs/index.html", s.addr))
	if err := s.router.Run(fmt.Sprintf(":%s", API_SVC_PORT)); err != nil {
		s.logger.Fatal("Failed to start server:", zap.Error(err))
	}
}
