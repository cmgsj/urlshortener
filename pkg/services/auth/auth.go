package auth

import (
	"flag"
	"fmt"
	"net"

	"urlshortener/pkg/interceptor"
	"urlshortener/pkg/logger"
	"urlshortener/pkg/proto/authpb"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
)

var (
	port   = flag.Int("port", 8083, "the port to serve on")
	dbAddr = flag.String("db_addr", "postgres_db:8085", "the db address")
)

func NewService() *Service {
	flag.Parse()
	service := &Service{
		db: initSqliteDB(*dbAddr),
	}
	return service
}

func (service *Service) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		logger.Fatal("failed to listen:", err)
	}
	loggerInterceptor := interceptor.NewLoggerInterceptor()
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggerInterceptor.Unary),
		grpc.StreamInterceptor(loggerInterceptor.Stream),
	)
	authpb.RegisterAuthServiceServer(grpcServer, service)
	logger.Info("Starting auth_service at:", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("failed to serve:", err)
	}
}
