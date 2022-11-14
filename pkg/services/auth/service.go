package auth

import (
	"context"
	"database/sql"
	"time"
	"urlshortener/pkg/proto/authpb"
	"urlshortener/pkg/proto/healthpb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InvalidUrlArgumentError = status.Error(codes.InvalidArgument, "invalid url argument")
	UrlNotFoundError        = status.Error(codes.NotFound, "url not found")
	UrlAlreadyExistsError   = status.Error(codes.AlreadyExists, "url already exists")
	InternalServerError     = status.Error(codes.Internal, "internal server error")
)

type Service struct {
	authpb.UnimplementedAuthServiceServer
	healthpb.UnimplementedHealthServiceServer
	db *sql.DB
}

func (service *Service) Login(context.Context, *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}

func (service *Service) VerifyToken(context.Context, *authpb.VerifyTokenRequest) (*authpb.VerifyTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyToken not implemented")
}

func (service *Service) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (service *Service) Watch(req *healthpb.HealthCheckRequest, stream healthpb.HealthService_WatchServer) error {
	for {
		select {
		case <-stream.Context().Done():
			return nil
		default:
			err := stream.Send(&healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING})
			if err != nil {
				return err
			}
			time.Sleep(time.Minute)
		}
	}
}
