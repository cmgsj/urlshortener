package main

import (
	"context"
	"net"

	"net/http"

	"github.com/cmgsj/go-env/env"
	urlv1 "github.com/cmgsj/urlshortener/pkg/gen/proto/url/v1"
	"github.com/cmgsj/urlshortener/pkg/grpcutil/interceptor"
	"github.com/cmgsj/urlshortener/pkg/urlsvc"
	"github.com/cmgsj/urlshortener/pkg/urlsvc/database"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	var (
		logger         = zap.Must(zap.NewDevelopment())
		urlSvcGrpcPort = env.MustGet("URL_SVC_PORT")
		urlSvcHttpPort = env.MustGet("URL_SVC_HTTP_PORT")
		urlDbUri       = env.MustGet("URL_DB_URI")

		opts = database.Options{
			Driver:      "sqlite3",
			ConnString:  urlDbUri,
			AutoMigrate: true,
		}
		svc = urlsvc.New(logger, database.Must(database.New(opts)))
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

	lis, err := net.Listen("tcp", ":"+urlSvcGrpcPort)
	if err != nil {
		svc.Logger.Fatal("failed to listen:", zap.Error(err))
	}

	go runHTTPServer(logger, urlSvcHttpPort, urlSvcGrpcPort)

	svc.Logger.Info("starting grpc server", zap.String("service", "url-svc"), zap.String("port", urlSvcGrpcPort))
	if err := grpcServer.Serve(lis); err != nil {
		svc.Logger.Fatal("failed to serve:", zap.Error(err))
	}
}

func runHTTPServer(logger *zap.Logger, urlSvcGrpcPort, urlSvcHttpPort string) {
	conn, err := grpc.Dial(":"+urlSvcGrpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("failed to dial grpc", zap.String("port", urlSvcGrpcPort), zap.Error(err))
	}
	defer conn.Close()

	ctx := context.Background()
	rmux := runtime.NewServeMux()
	client := urlv1.NewUrlServiceClient(conn)
	if err = urlv1.RegisterUrlServiceHandlerClient(ctx, rmux, client); err != nil {
		logger.Fatal("failed register grpc-gateway", zap.Error(err))
	}

	mux := http.NewServeMux()
	mux.Handle("/", rmux)
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("./swagger-ui"))))
	mux.Handle("/swagger-ui/swagger.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "./swagger.json") }))

	logger.Info("starting http server", zap.String("service", "url-svc"), zap.String("port", urlSvcHttpPort))
	if err = http.ListenAndServe(":"+urlSvcHttpPort, mux); err != nil {
		logger.Fatal("failed to serve http", zap.String("port", urlSvcHttpPort), zap.Error(err))
	}
}
