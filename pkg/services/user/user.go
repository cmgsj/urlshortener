package auth

import (
	"flag"
	"fmt"
	"net"

	"urlshortener/pkg/interceptor"
	"urlshortener/pkg/logger"
	"urlshortener/pkg/proto/userpb"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
)

var (
	port   = flag.Int("port", 8081, "the port to serve on")
	urlsDB = flag.String("urls_db", "urls.sqlite", "the urls db")
)

func NewService() *Service {
	flag.Parse()
	service := &Service{
		db: initSqliteDB(*urlsDB),
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
	userpb.RegisterUserServiceServer(grpcServer, service)
	logger.Info("Starting urls_service at:", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("failed to serve:", err)
	}
}
