package urls

import (
	"fmt"
	"grpc_util/pkg/grpc_interceptor"
	"net"
	"os"
	"proto/pkg/urlspb"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	API_SVC_NAME  = os.Getenv("API_SVC_NAME")
	URLS_SVC_NAME = os.Getenv("URLS_SVC_NAME")
	URLS_SVC_PORT = os.Getenv("URLS_SVC_PORT")
	DB_URI        = os.Getenv("DB_URI")
)

func NewService() *Service {
	service := &Service{
		healthServer: health.NewServer(),
		logger:       zap.Must(zap.NewDevelopment()),
	}
	service.intiDB(DB_URI)
	service.healthServer.SetServingStatus(API_SVC_NAME, healthpb.HealthCheckResponse_SERVING)
	return service
}

func (s *Service) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", URLS_SVC_PORT))
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
	urlspb.RegisterUrlsServiceServer(grpcServer, s)

	s.logger.Info("Starting", zap.String("service", URLS_SVC_NAME), zap.String("address", lis.Addr().String()))
	if err := grpcServer.Serve(lis); err != nil {
		s.logger.Fatal("failed to serve:", zap.Error(err))
	}
}
