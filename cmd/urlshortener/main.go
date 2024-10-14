package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/cmgsj/go-lib/swagger"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	reflectionv1 "google.golang.org/grpc/reflection/grpc_reflection_v1"
	reflectionv1alpha "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/cmgsj/urlshortener/pkg/database"
	"github.com/cmgsj/urlshortener/pkg/docs"
	urlshortenerv1 "github.com/cmgsj/urlshortener/pkg/gen/proto/urlshortener/v1"
	urlshortenerserver "github.com/cmgsj/urlshortener/pkg/urlshortener/server"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	pflag.String("grpc-address", "127.0.0.1:9090", "urlshortener server grpc address")
	pflag.String("http-address", "127.0.0.1:8080", "urlshortener server http address")

	pflag.Parse()

	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)
	viper.SetEnvPrefix("urlshortener")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.BindPFlags(pflag.CommandLine)

	grpcAddress := viper.GetString("grpc-address")
	httpAddress := viper.GetString("http-address")

	db := database.Must(database.New(database.Options{
		Driver:      "sqlite3",
		ConnString:  ":memory:",
		AutoMigrate: true,
	}))

	urlshortenerServer := urlshortenerserver.NewServer(db)

	err := urlshortenerServer.SeedDB(ctx)
	if err != nil {
		return err
	}

	healthServer := health.NewServer()
	healthServer.SetServingStatus(urlshortenerv1.URLShortenerService_ServiceDesc.ServiceName, healthv1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus(healthv1.Health_ServiceDesc.ServiceName, healthv1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus(reflectionv1.ServerReflection_ServiceDesc.ServiceName, healthv1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus(reflectionv1alpha.ServerReflection_ServiceDesc.ServiceName, healthv1.HealthCheckResponse_SERVING)

	grpcServer := grpc.NewServer()

	urlshortenerv1.RegisterURLShortenerServiceServer(grpcServer, urlshortenerServer)
	healthv1.RegisterHealthServer(grpcServer, healthServer)
	reflection.Register(grpcServer)

	rmux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	err = urlshortenerv1.RegisterURLShortenerServiceHandlerFromEndpoint(ctx, rmux, grpcAddress, opts)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	mux.Handle("/", rmux)
	mux.Handle("/r/", http.StripPrefix("/r/", urlshortenerServer.RedirectUrl()))
	mux.Handle("/docs/", swagger.Docs("/docs/", docs.SwaggerSchema()))

	httpServer := &http.Server{
		Handler: mux,
	}

	slog.Info("starting urlshortener server", "grpc_address", grpcAddress, "http_address", httpAddress)

	errch := make(chan error)

	go func() {
		lis, err := net.Listen("tcp", grpcAddress)
		if err != nil {
			errch <- err
		}
		errch <- grpcServer.Serve(lis)
	}()

	go func() {
		lis, err := net.Listen("tcp", httpAddress)
		if err != nil {
			errch <- err
		}
		errch <- httpServer.Serve(lis)
	}()

	return <-errch
}
