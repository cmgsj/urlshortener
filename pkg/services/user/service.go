package auth

import (
	"context"
	"database/sql"
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
	db *sql.DB
}

func (s *Service) GetUserById(ctx context.Context, req *userpb.GetUserByIdRequest) (*userpb.GetUserResponse, error) {
	userEntity, err := getUserById(ctx, s.db, req.GetUserId())
	if err != nil {
		return nil, UserNotFoundError
	}
	return &userpb.GetUserResponse{User: &userpb.User{UserId: userEntity.UserId, Email: userEntity.Email}}, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, req *userpb.GetUserByEmailRequest) (*userpb.GetUserResponse, error) {
	userEntity, err := getUserByEmail(ctx, s.db, req.GetEmail())
	if err != nil {
		return nil, UserNotFoundError
	}
	return &userpb.GetUserResponse{User: &userpb.User{UserId: userEntity.UserId, Email: userEntity.Email}}, nil
}

func (s *Service) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.NoContent, error) {
	email := req.GetEmail()
	if !isValidEmail(email) {
		return nil, InvalidEmailError
	}
	// TODO: hash password
	err := createUser(ctx, s.db, email, req.GetPassword())
	if err != nil {
		return nil, InvalidEmailOrPasswordError
	}
	return &userpb.NoContent{}, nil
}

func isValidEmail(email string) bool {
	return email != ""
}
