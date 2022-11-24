package urls

import (
	"flag"
	"fmt"
	"net"
	"urlshortener/pkg/interceptor"
	"urlshortener/pkg/proto/urlspb"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
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
		healthServer: health.NewServer(),
		logger:       zap.Must(zap.NewDevelopment()),
	}
	service.intiDB(*urlsDB)
	service.healthServer.SetServingStatus("api.service", healthpb.HealthCheckResponse_SERVING)
	return service
}

func (s *Service) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		s.logger.Fatal("failed to listen:", zap.Error(err))
	}

	loggerInterceptor := interceptor.NewLogger()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggerInterceptor.Unary),
		grpc.StreamInterceptor(loggerInterceptor.Stream),
	)
	reflection.Register(grpcServer)

	healthpb.RegisterHealthServer(grpcServer, s.healthServer)
	urlspb.RegisterUrlsServiceServer(grpcServer, s)

	s.logger.Info("Starting urls.service at:", zap.String("address", lis.Addr().String()))
	if err := grpcServer.Serve(lis); err != nil {
		s.logger.Fatal("failed to serve:", zap.Error(err))
	}
}
