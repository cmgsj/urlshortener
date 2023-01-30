package api

import (
	_ "api_service/pkg/docs"
	"api_service/pkg/grpc_util/grpc_health"
	"api_service/pkg/proto/urlspb"
	"strconv"
	"time"

	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	API_SVC_PORT           = os.Getenv("API_SVC_PORT")
	API_SVC_ADDR           = os.Getenv("API_SVC_ADDR")
	API_SVC_NAME           = os.Getenv("API_SVC_NAME")
	API_SVC_CACHE_EXP_TIME = os.Getenv("API_SVC_CACHE_EXP_TIME")
	REDIS_ADDR             = os.Getenv("REDIS_ADDR")
	REDIS_PASSWORD         = os.Getenv("REDIS_PASSWORD")
	REDIS_DB               = os.Getenv("REDIS_DB")
	URLS_SVC_ADDR          = os.Getenv("URLS_SVC_ADDR")
	URLS_SVC_NAME          = os.Getenv("URLS_SVC_NAME")
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
	logger := zap.Must(zap.NewDevelopment())
	urlsConn, err := grpc.Dial(URLS_SVC_ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("failed to dial urls service:", zap.Error(err))
	}
	cacheExpTime, err := strconv.Atoi(API_SVC_CACHE_EXP_TIME)
	if err != nil {
		logger.Fatal("failed to convert API_SVC_CACHE_EXP_TIME to int:", zap.Error(err))
	}
	redisdb, err := strconv.Atoi(REDIS_DB)
	if err != nil {
		logger.Fatal("failed to convert REDIS_DB to int:", zap.Error(err))
	}
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())
	service := &Service{
		name:             API_SVC_NAME,
		addr:             API_SVC_ADDR,
		trustedProxies:   []string{"127.0.0.1"},
		router:           router,
		logger:           logger,
		urlsClient:       urlspb.NewUrlsServiceClient(urlsConn),
		urlsHealthClient: healthpb.NewHealthClient(urlsConn),
		urlsServiceName:  URLS_SVC_NAME,
		cacheExpTime:     time.Duration(cacheExpTime) * time.Second,
	}
	service.initRedisDB(REDIS_ADDR, REDIS_PASSWORD, redisdb)
	service.RegisterEndpoints()
	service.RegisterTrustedProxies()
	return service
}

func (s *Service) Run() {
	go grpc_health.WatchServices(s.name, s.logger, time.Second, []*grpc_health.HealthClient{
		{HealthClient: s.urlsHealthClient, Active: &s.urlsServiceOk, Name: s.urlsServiceName},
	})
	s.addr = fmt.Sprintf("localhost:%s", API_SVC_PORT) // TODO: delete this line
	s.logger.Info("Starting", zap.String("service", API_SVC_NAME), zap.String("port", API_SVC_PORT))
	s.logger.Info(fmt.Sprintf("Swagger docs available at http://%s/docs/index.html", s.addr))
	if err := s.router.Run(fmt.Sprintf(":%s", API_SVC_PORT)); err != nil {
		s.logger.Fatal("Failed to start server:", zap.Error(err))
	}
}
