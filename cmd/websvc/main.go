package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/cmgsj/urlshortener/pkg/proto/urlpb"
	"github.com/cmgsj/urlshortener/pkg/websvc"
	_ "github.com/cmgsj/urlshortener/pkg/websvc/docs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ApiSvcPort         = os.Getenv("WEB_SVC_PORT")
	ApiSvcUrl          = os.Getenv("WEB_SVC_URL")
	ApiSvcCacheTimeout = os.Getenv("WEB_SVC_CACHE_TIMEOUT")
	RedisAddr          = os.Getenv("REDIS_URL")
	RedisPassword      = os.Getenv("REDIS_PASSWORD")
	RedisDb            = os.Getenv("REDIS_DB")
	UrlSvcAddr         = os.Getenv("URL_SVC_ADDR")
)

func main() {
	logger := zap.Must(zap.NewDevelopment())

	urlConn, err := grpc.Dial(UrlSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("failed to dial url service:", zap.Error(err))
	}

	cacheTimeout, err := time.ParseDuration(ApiSvcCacheTimeout)
	if err != nil {
		logger.Fatal("failed to parse WEB_SVC_CACHE_TIMEOUT:", zap.Error(err))
	}

	redisDb, err := strconv.Atoi(RedisDb)
	if err != nil {
		logger.Fatal("failed to convert REDIS_DB to int:", zap.Error(err))
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	svc := &websvc.Service{
		Addr:           ApiSvcUrl,
		TrustedProxies: []string{"127.0.0.1"},
		Router:         router,
		Logger:         logger,
		UrlClient:      urlpb.NewUrlServiceClient(urlConn),
		CacheTimeout:   cacheTimeout,
	}

	svc.InitRedisDB(RedisAddr, RedisPassword, redisDb)
	svc.RegisterEndpoints()
	svc.RegisterTrustedProxies()

	svc.Logger.Info("starting", zap.String("service", "web-svc"), zap.String("port", ApiSvcPort))
	svc.Logger.Info(fmt.Sprintf("swagger docs available at http://%s/docs/index.html", svc.Addr))

	if err := svc.Router.Run(fmt.Sprintf(":%s", ApiSvcPort)); err != nil {
		svc.Logger.Fatal("failed to start server:", zap.Error(err))
	}
}
