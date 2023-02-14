package main

import (
	"strconv"
	"time"
	"urlshortener/pkg/api_service"
	_ "urlshortener/pkg/api_service/docs"
	"urlshortener/pkg/grpc_util/grpc_health"
	"urlshortener/pkg/proto/urlpb"

	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	ApiSvcPort         = os.Getenv("API_SVC_PORT")
	ApiSvcUrl          = os.Getenv("API_SVC_URL")
	ApiSvcName         = os.Getenv("API_SVC_NAME")
	ApiSvcCacheTimeout = os.Getenv("API_SVC_CACHE_TIMEOUT")
	RedisAddr          = os.Getenv("REDIS_URL")
	RedisPassword      = os.Getenv("REDIS_PASSWORD")
	RedisDb            = os.Getenv("REDIS_DB")
	UrlSvcAddr         = os.Getenv("URL_SVC_ADDR")
	UrlSvcName         = os.Getenv("URL_SVC_NAME")
)

func main() {
	logger := zap.Must(zap.NewDevelopment())

	urlConn, err := grpc.Dial(UrlSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("failed to dial url service:", zap.Error(err))
	}

	cacheTimeout, err := time.ParseDuration(ApiSvcCacheTimeout)
	if err != nil {
		logger.Fatal("failed to parse API_SVC_CACHE_TIMEOUT:", zap.Error(err))
	}

	redisDb, err := strconv.Atoi(RedisDb)
	if err != nil {
		logger.Fatal("failed to convert REDIS_DB to int:", zap.Error(err))
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	svc := &api_service.Service{
		Name:            ApiSvcName,
		Addr:            ApiSvcUrl,
		TrustedProxies:  []string{"127.0.0.1"},
		Router:          router,
		Logger:          logger,
		UrlClient:       urlpb.NewUrlServiceClient(urlConn),
		UrlHealthClient: healthpb.NewHealthClient(urlConn),
		UrlServiceName:  UrlSvcName,
		CacheTimeout:    cacheTimeout,
	}

	svc.InitRedisDB(RedisAddr, RedisPassword, redisDb)
	svc.RegisterEndpoints()
	svc.RegisterTrustedProxies()

	go grpc_health.WatchServices(svc.Name, svc.Logger, time.Second, []*grpc_health.HealthClient{
		{HealthClient: svc.UrlHealthClient, Active: &svc.UrlServiceOk, Name: svc.UrlServiceName},
	})

	svc.Logger.Info("Starting", zap.String("service", ApiSvcName), zap.String("port", ApiSvcPort))
	svc.Logger.Info(fmt.Sprintf("Swagger docs available at http://%s/docs/index.html", svc.Addr))

	if err := svc.Router.Run(fmt.Sprintf(":%s", ApiSvcPort)); err != nil {
		svc.Logger.Fatal("Failed to start server:", zap.Error(err))
	}
}
