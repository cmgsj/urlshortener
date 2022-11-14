package api

import (
	"context"
	"urlshortener/pkg/proto/healthpb"

	"google.golang.org/grpc"
)

type UrlDTO struct {
	UrlId       string `json:"urlId"`
	RedirectUrl string `json:"redirectUrl"`
	NewUrl      string `json:"newUrl"`
}

type CreateUrlRequest struct {
	RedirectUrl string `json:"redirectUrl" binding:"required"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type healthServer interface {
	Check(ctx context.Context, in *healthpb.HealthCheckRequest, opts ...grpc.CallOption) (*healthpb.HealthCheckResponse, error)
}
