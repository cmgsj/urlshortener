package main

import (
	"api_service/pkg/api"
	_ "api_service/pkg/docs"
	"grpc_util/pkg/grpc_health"
	"proto/pkg/urlspb"
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

func main() {
	logger := zap.Must(zap.NewDevelopment())
	logger.Info("", zap.String("URLS_SVC_ADDR", URLS_SVC_ADDR))
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
	svc := &api.Service{
		Name:             API_SVC_NAME,
		Addr:             API_SVC_ADDR,
		TrustedProxies:   []string{"127.0.0.1"},
		Router:           router,
		Logger:           logger,
		UrlsClient:       urlspb.NewUrlsServiceClient(urlsConn),
		UrlsHealthClient: healthpb.NewHealthClient(urlsConn),
		UrlsServiceName:  URLS_SVC_NAME,
		CacheExpTime:     time.Duration(cacheExpTime) * time.Second,
	}
	svc.InitRedisDB(REDIS_ADDR, REDIS_PASSWORD, redisdb)
	svc.RegisterEndpoints()
	svc.RegisterTrustedProxies()
	go grpc_health.WatchServices(svc.Name, svc.Logger, time.Second, []*grpc_health.HealthClient{
		{HealthClient: svc.UrlsHealthClient, Active: &svc.UrlsServiceOk, Name: svc.UrlsServiceName},
	})
	svc.Addr = fmt.Sprintf("localhost:%s", API_SVC_PORT) // TODO: delete this line
	svc.Logger.Info("Starting", zap.String("service", API_SVC_NAME), zap.String("port", API_SVC_PORT))
	svc.Logger.Info(fmt.Sprintf("Swagger docs available at http://%s/docs/index.html", svc.Addr))
	if err := svc.Router.Run(fmt.Sprintf(":%s", API_SVC_PORT)); err != nil {
		svc.Logger.Fatal("Failed to start server:", zap.Error(err))
	}
}
