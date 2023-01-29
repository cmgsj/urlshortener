package auth

import (
	"fmt"
	"grpc_util/pkg/grpc_interceptor"

	"net"
	"os"
	"proto/pkg/authpb"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	AUTH_SVC_NAME = os.Getenv("AUTH_SVC_NAME")
	AUTH_SVC_PORT = os.Getenv("AUTH_SVC_PORT")
	DB_URI        = os.Getenv("DB_URI")
)

func NewService() *Service {

	service := &Service{
		healthServer: health.NewServer(),
		logger:       zap.Must(zap.NewDevelopment()),
	}
	service.intiDB(DB_URI)
	return service
}

func (s *Service) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", AUTH_SVC_PORT))
	if err != nil {
		s.logger.Fatal("failed to listen:", zap.Error(err))
	}

	loggerInterceptor := grpc_interceptor.NewLogger(s.logger)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggerInterceptor.Unary),
		grpc.StreamInterceptor(loggerInterceptor.Stream),
	)
	reflection.Register(grpcServer)
	healthpb.RegisterHealthServer(grpcServer, s.healthServer)
	authpb.RegisterAuthServiceServer(grpcServer, s)

	s.logger.Info("Starting", zap.String("service", AUTH_SVC_NAME), zap.String("address", lis.Addr().String()))
	if err := grpcServer.Serve(lis); err != nil {
		s.logger.Fatal("failed to serve:", zap.Error(err))
	}
}
