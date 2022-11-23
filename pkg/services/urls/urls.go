package urls

import (
	"context"
	"flag"
	"fmt"
	"net"

	"urlshortener/pkg/db"
	"urlshortener/pkg/interceptor"
	"urlshortener/pkg/logger"
	"urlshortener/pkg/proto/urlspb"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	port   = flag.Int("port", 8081, "the port to serve on")
	urlsDB = flag.String("urls_db", "urls.sqlite", "the urls db")
)

func NewService() *Service {
	flag.Parse()
	service := &Service{
		db: intiDB(*urlsDB),
	}
	db.CreateTables(context.Background(), service.db)
	db.SeedDB(context.Background(), service.db)
	return service
}

func (s *Service) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		logger.Fatal("failed to listen:", err)
	}

	loggerInterceptor := interceptor.NewLoggerInterceptor()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggerInterceptor.Unary),
		grpc.StreamInterceptor(loggerInterceptor.Stream),
	)

	reflection.Register(grpcServer) // for testing

	healthServer := health.NewServer()

	healthServer.SetServingStatus("api.service", healthpb.HealthCheckResponse_SERVING)

	healthpb.RegisterHealthServer(grpcServer, healthServer)
	urlspb.RegisterUrlsServiceServer(grpcServer, s)

	logger.Info("Starting urls.service at:", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("failed to serve:", err)
	}
}
