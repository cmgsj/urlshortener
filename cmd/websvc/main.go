package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	urlv1 "github.com/cmgsj/urlshortener/pkg/gen/proto/url/v1"
	"github.com/cmgsj/urlshortener/pkg/websvc"
	_ "github.com/cmgsj/urlshortener/pkg/websvc/docs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var (
		logger             = zap.Must(zap.NewDevelopment())
		apiSvcPort         = os.Getenv("WEB_SVC_PORT")
		apiSvcUrl          = os.Getenv("WEB_SVC_URL")
		apiSvcCacheTimeout = os.Getenv("WEB_SVC_CACHE_TIMEOUT")
		redisAddr          = os.Getenv("REDIS_URL")
		redisPassword      = os.Getenv("REDIS_PASSWORD")
		redisDb            = os.Getenv("REDIS_DB")
		urlSvcAddr         = os.Getenv("URL_SVC_ADDR")
	)

	urlConn, err := grpc.Dial(urlSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("failed to dial url service:", zap.Error(err))
	}

	cacheTimeout, err := time.ParseDuration(apiSvcCacheTimeout)
	if err != nil {
		logger.Fatal("failed to parse WEB_SVC_CACHE_TIMEOUT:", zap.Error(err))
	}

	redisDbNum, err := strconv.Atoi(redisDb)
	if err != nil {
		logger.Fatal("failed to convert REDIS_DB to int:", zap.Error(err))
	}

	opt := websvc.Options{
		Addr:           apiSvcUrl,
		TrustedProxies: []string{"127.0.0.1"},
		Logger:         logger,
		UrlClient:      urlv1.NewUrlServiceClient(urlConn),
		CacheTimeout:   cacheTimeout,
	}
	svc := websvc.New(opt)

	svc.InitRedisDB(redisAddr, redisPassword, redisDbNum)
	svc.RegisterEndpoints()
	svc.RegisterTrustedProxies()

	svc.Logger.Info("starting", zap.String("service", "web-svc"), zap.String("port", apiSvcPort))
	svc.Logger.Info(fmt.Sprintf("swagger docs available at http://%s/docs/index.html", svc.Addr))

	if err := svc.Router.Run(fmt.Sprintf(":%s", apiSvcPort)); err != nil {
		svc.Logger.Fatal("failed to start server:", zap.Error(err))
	}
}
