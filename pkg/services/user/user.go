package auth

import (
	"context"
	"flag"
	"fmt"
	"net"

	"urlshortener/pkg/db"
	"urlshortener/pkg/interceptor"
	"urlshortener/pkg/logger"
	"urlshortener/pkg/proto/userpb"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
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

	healthServer := health.NewServer()

	healthpb.RegisterHealthServer(grpcServer, healthServer)
	userpb.RegisterUserServiceServer(grpcServer, s)

	logger.Info("Starting user.service at:", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("failed to serve:", err)
	}
}
