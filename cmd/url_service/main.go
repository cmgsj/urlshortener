package main

import (
	"fmt"
	"net"
	"os"
	"urlshortener/pkg/grpc_util/grpc_interceptor"
	"urlshortener/pkg/proto/urlpb"
	url_service "urlshortener/pkg/url_service"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	ApiSvcName = os.Getenv("API_SVC_NAME")
	UrlSvcName = os.Getenv("URL_SVC_NAME")
	UrlSvcPort = os.Getenv("URL_SVC_PORT")
	UrlDbUri   = os.Getenv("URL_DB_URI")
)

func main() {
	svc := &url_service.Service{
		HealthServer: health.NewServer(),
		Logger:       zap.Must(zap.NewDevelopment()),
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
	healthpb.RegisterHealthServer(grpcServer, svc.HealthServer)
	urlpb.RegisterUrlServiceServer(grpcServer, svc)
	svc.HealthServer.SetServingStatus(ApiSvcName, healthpb.HealthCheckResponse_SERVING)

	svc.Logger.Info("Starting", zap.String("service", UrlSvcName), zap.String("address", lis.Addr().String()))
	if err := grpcServer.Serve(lis); err != nil {
		svc.Logger.Fatal("failed to serve:", zap.Error(err))
	}
}
