package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/cmgsj/urlshortener/pkg/grpcutil/interceptor"
	"github.com/cmgsj/urlshortener/pkg/proto/urlpb"
	"github.com/cmgsj/urlshortener/pkg/urlsvc"
	"github.com/cmgsj/urlshortener/pkg/urlsvc/db"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var (
		ctx        = context.Background()
		urlSvcPort = os.Getenv("URL_SVC_PORT")
		urlDbUri   = os.Getenv("URL_DB_URI")
		logger     = zap.Must(zap.NewDevelopment())
		querier    = db.MustPrepare(ctx, db.Must(db.Migrate(ctx, db.Must(db.Connect("sqlite3", urlDbUri)))))
		svc        = urlsvc.New(logger, querier)
	)

	if err := svc.SeedDB(ctx); err != nil {
		svc.Logger.Fatal("failed to seed DB:", zap.Error(err))
	}

	logInterceptor := interceptor.NewLogger(svc.Logger)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(logInterceptor.Unary),
		grpc.StreamInterceptor(logInterceptor.Stream),
	)

	reflection.Register(grpcServer)
	urlpb.RegisterUrlServiceServer(grpcServer, svc)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", urlSvcPort))
	if err != nil {
		svc.Logger.Fatal("failed to listen:", zap.Error(err))
	}

	svc.Logger.Info("starting", zap.String("service", "url-svc"), zap.String("address", lis.Addr().String()))
	if err := grpcServer.Serve(lis); err != nil {
		svc.Logger.Fatal("failed to serve:", zap.Error(err))
	}
}
