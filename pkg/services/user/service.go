package auth

import (
	"context"
	"database/sql"
	"time"
	"urlshortener/pkg/proto/healthpb"
	"urlshortener/pkg/proto/userpb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	UserNotFoundError           = status.Error(codes.NotFound, "user not found")
	InvalidEmailError           = status.Error(codes.InvalidArgument, "invalid email argument")
	InvalidEmailOrPasswordError = status.Error(codes.InvalidArgument, "invalid email argument")
)

type Service struct {
	userpb.UnimplementedUserServiceServer
	healthpb.UnimplementedHealthServiceServer
	db *sql.DB
}

func (service *Service) GetUserById(ctx context.Context, req *userpb.GetUserByIdRequest) (*userpb.GetUserResponse, error) {
	userEntity, err := getUserById(service.db, ctx, req.GetUserId())
	if err != nil {
		return nil, UserNotFoundError
	}
	return &userpb.GetUserResponse{User: &userpb.User{UserId: userEntity.UserId, Email: userEntity.Email}}, nil
}

func (service *Service) GetUserByEmail(ctx context.Context, req *userpb.GetUserByEmailRequest) (*userpb.GetUserResponse, error) {
	userEntity, err := getUserByEmail(service.db, ctx, req.GetEmail())
	if err != nil {
		return nil, UserNotFoundError
	}
	return &userpb.GetUserResponse{User: &userpb.User{UserId: userEntity.UserId, Email: userEntity.Email}}, nil
}

func (service *Service) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.NoContent, error) {
	email := req.GetEmail()
	if !isValidEmail(email) {
		return nil, InvalidEmailError
	}
	// TODO: hash password
	err := createUser(service.db, ctx, email, req.GetPassword())
	if err != nil {
		return nil, InvalidEmailOrPasswordError
	}
	return &userpb.NoContent{}, nil
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

func isValidEmail(email string) bool {
	return email != ""
}
