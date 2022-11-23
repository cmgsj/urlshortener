package auth

import (
	"context"
	"database/sql"
	"urlshortener/pkg/proto/authpb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	authpb.UnimplementedAuthServiceServer
	db *sql.DB
}

func (s *Service) Login(context.Context, *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
