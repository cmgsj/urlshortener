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
		urlSvcGrpcPort = env.GetDefault("URL_SVC_PORT", "9090")
		urlSvcHttpPort = env.GetDefault("URL_SVC_HTTP_PORT", "8080")
		urlDbUri       = env.GetDefault("URL_DB_URI", ":memory:")

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

	go serveHTTP(logger, urlSvcHttpPort, urlSvcGrpcPort)

	svc.Logger.Info("starting", zap.String("service", "url-svc"), zap.String("address", lis.Addr().String()))
	if err := grpcServer.Serve(lis); err != nil {
		svc.Logger.Fatal("failed to serve:", zap.Error(err))
	}
}

func serveHTTP(logger *zap.Logger, httpPort, grpcPort string) {
	conn, err := grpc.Dial(":"+grpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("failed to dial grpc", zap.String("port", grpcPort), zap.Error(err))
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
	mux.Handle("/swagger-ui/swagger.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "pkg/gen/proto/url/v1/url.swagger.json")
	}))
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("./swagger-ui"))))

	if err = http.ListenAndServe(":"+httpPort, mux); err != nil {
		logger.Fatal("failed to serve http", zap.String("port", grpcPort), zap.Error(err))
	}
}
