package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cmgsj/go-env/env"
	urlv1 "github.com/cmgsj/urlshortener/pkg/gen/proto/url/v1"
	"github.com/cmgsj/urlshortener/pkg/websvc"
	_ "github.com/cmgsj/urlshortener/pkg/websvc/docs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var (
		logger = zap.Must(zap.NewDevelopment())
		// TODO:  change to env.GetDefault()
		apiSvcPort         = env.MustGet("WEB_SVC_PORT")
		apiSvcUrl          = env.MustGet("WEB_SVC_URL")
		apiSvcCacheTimeout = env.MustGet("WEB_SVC_CACHE_TIMEOUT")
		redisAddr          = env.MustGet("REDIS_URL")
		redisPassword      = env.MustGet("REDIS_PASSWORD")
		redisDb            = env.MustGet("REDIS_DB")
		urlSvcAddr         = env.MustGet("URL_SVC_ADDR")
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

	if err := svc.Router.Run(":" + apiSvcPort); err != nil {
		svc.Logger.Fatal("failed to start server:", zap.Error(err))
	}
}
