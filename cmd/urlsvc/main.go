package main

import (
	"context"
	"net"
	"os"

	urlv1 "github.com/cmgsj/urlshortener/pkg/gen/proto/url/v1"
	"github.com/cmgsj/urlshortener/pkg/grpcutil/interceptor"
	"github.com/cmgsj/urlshortener/pkg/urlsvc"
	"github.com/cmgsj/urlshortener/pkg/urlsvc/database"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var (
		urlSvcPort = os.Getenv("URL_SVC_PORT")
		urlDbUri   = os.Getenv("URL_DB_URI")
		logger     = zap.Must(zap.NewDevelopment())
		opts       = database.Options{
			Driver:  "sqlite3",
			URI:     urlDbUri,
			Migrate: true,
		}
		db  = database.Must(database.New(opts))
		svc = urlsvc.New(logger, db)
	)

	if err := svc.SeedDB(context.Background()); err != nil {
		svc.Logger.Fatal("failed to seed DB:", zap.Error(err))
	}

	logInterceptor := interceptor.NewLogger(svc.Logger)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(logInterceptor.Unary),
		grpc.StreamInterceptor(logInterceptor.Stream),
	)

	reflection.Register(grpcServer)
	urlv1.RegisterUrlServiceServer(grpcServer, svc)

	lis, err := net.Listen("tcp", ":"+urlSvcPort)
	if err != nil {
		svc.Logger.Fatal("failed to listen:", zap.Error(err))
	}

	svc.Logger.Info("starting", zap.String("service", "url-svc"), zap.String("address", lis.Addr().String()))
	if err := grpcServer.Serve(lis); err != nil {
		svc.Logger.Fatal("failed to serve:", zap.Error(err))
	}
}
