package main

import (
	"fmt"
	"net"
	"os"

	"github.com/cmgsj/url-shortener/pkg/grpc_util/grpc_interceptor"
	"github.com/cmgsj/url-shortener/pkg/proto/urlpb"
	"github.com/cmgsj/url-shortener/pkg/url_service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	UrlSvcPort = os.Getenv("URL_SVC_PORT")
	UrlDbUri   = os.Getenv("URL_DB_URI")
)

func main() {
	svc := &url_service.Service{
		Logger: zap.Must(zap.NewDevelopment()),
	}
	svc.IntiDB(UrlDbUri)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", UrlSvcPort))
	if err != nil {
		svc.Logger.Fatal("failed to listen:", zap.Error(err))
	}

	loggerInterceptor := grpc_interceptor.NewLogger(svc.Logger)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggerInterceptor.Unary),
		grpc.StreamInterceptor(loggerInterceptor.Stream),
	)

	reflection.Register(grpcServer)
	urlpb.RegisterUrlServiceServer(grpcServer, svc)

	svc.Logger.Info("starting", zap.String("service", "url-svc"), zap.String("address", lis.Addr().String()))
	if err := grpcServer.Serve(lis); err != nil {
		svc.Logger.Fatal("failed to serve:", zap.Error(err))
	}
}
