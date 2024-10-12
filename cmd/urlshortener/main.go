package main

import (
	"context"
	"net"
	"net/http"

	"github.com/cmgsj/go-env/env"
	"github.com/cmgsj/go-lib/swagger"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/cmgsj/urlshortener/pkg/database"
	"github.com/cmgsj/urlshortener/pkg/docs"
	urlshortenerv1 "github.com/cmgsj/urlshortener/pkg/gen/proto/urlshortener/v1"
	"github.com/cmgsj/urlshortener/pkg/service"
)

func main() {
	var (
		grpcAddr = ":" + env.GetDefault("GRPC_PORT", "9090")
		httpAddr = ":" + env.GetDefault("HTTP_PORT", "8080")
	)

	svc := &service.Service{
		DB: database.Must(database.New(database.Options{
			Driver:      "sqlite3",
			ConnString:  ":memory:",
			AutoMigrate: true,
		})),
	}

	err := svc.SeedDB(context.Background())
	check(err)

	hs := health.NewServer()
	hs.SetServingStatus(service.ServiceName, healthv1.HealthCheckResponse_SERVING)

	gs := grpc.NewServer()
	reflection.Register(gs)
	healthv1.RegisterHealthServer(gs, hs)
	urlshortenerv1.RegisterURLShortenerServer(gs, svc)

	rmux := runtime.NewServeMux()
	ctx := context.Background()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	err = urlshortenerv1.RegisterURLShortenerHandlerFromEndpoint(ctx, rmux, grpcAddr, opts)
	check(err)

	mux := http.NewServeMux()
	mux.Handle("/", rmux)
	mux.Handle("/r/", http.StripPrefix("/r/", svc.RedirectUrl()))
	mux.Handle("/docs/", swagger.Docs("/docs/", docs.SwaggerSchema()))

	go func() {
		hl, err := net.Listen("tcp", httpAddr)
		check(err)
		err = http.Serve(hl, mux)
		check(err)
	}()

	go func() {
		gl, err := net.Listen("tcp", grpcAddr)
		check(err)
		err = gs.Serve(gl)
		check(err)
	}()

	select {}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
