package main

import (
	"fmt"
	"grpc_util/pkg/grpc_interceptor"
	"net"
	"os"
	"proto/pkg/urlpb"
	"url_service/pkg/url"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	API_SVC_NAME = os.Getenv("API_SVC_NAME")
	URL_SVC_NAME = os.Getenv("URL_SVC_NAME")
	URL_SVC_PORT = os.Getenv("URL_SVC_PORT")
	URL_DB_URI   = os.Getenv("URL_DB_URI")
)

func main() {
	svc := &url.Service{
		HealthServer: health.NewServer(),
		Logger:       zap.Must(zap.NewDevelopment()),
	}
	svc.IntiDB(URL_DB_URI)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", URL_SVC_PORT))
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
	svc.HealthServer.SetServingStatus(API_SVC_NAME, healthpb.HealthCheckResponse_SERVING)

	svc.Logger.Info("Starting", zap.String("service", URL_SVC_NAME), zap.String("address", lis.Addr().String()))
	if err := grpcServer.Serve(lis); err != nil {
		svc.Logger.Fatal("failed to serve:", zap.Error(err))
	}
}
