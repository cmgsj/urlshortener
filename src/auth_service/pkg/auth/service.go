package auth

import (
	"auth_service/pkg/proto/authpb"
	"context"
	"database/sql"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/status"
)

var (
	ErrUserNotFound       = status.Error(codes.NotFound, "user not found")
	ErrInvalidCredentials = status.Error(codes.Unauthenticated, "invalid credentials")
	ErrInternal           = status.Error(codes.Internal, "internal error")
)

type Service struct {
	authpb.UnimplementedAuthServiceServer
	healthServer *health.Server
	logger       *zap.Logger
	db           *sql.DB
	jwtKey       string
}

func (s *Service) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	user, err := getUserByEmail(ctx, s.db, req.GetEmail())
	if err != nil {
		return nil, ErrUserNotFound
	}
	if err = CompareHash(user.Password, req.GetPassword()); err != nil {
		return nil, ErrInvalidCredentials
	}
	token, err := GenerateJwt(s.jwtKey, user.UserId, user.Email)
	if err != nil {
		return nil, ErrInternal
	}
	return &authpb.LoginResponse{Token: token}, nil
}

func (s *Service) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	hash, err := HashString(req.GetPassword())
	if err != nil {
		return nil, ErrInternal
	}
	user, err := createUser(ctx, s.db, req.GetEmail(), hash)
	if err != nil {
		return nil, ErrInternal
	}
	token, err := GenerateJwt(s.jwtKey, user.UserId, user.Email)
	if err != nil {
		return nil, ErrInternal
	}
	return &authpb.RegisterResponse{Token: token}, nil
}
