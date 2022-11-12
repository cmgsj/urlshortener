package api

import (
	"context"
	"github.com/mike9107/urlshortener/pkg/protobuf/apipb"

	"google.golang.org/grpc"
)

type UrlDTO struct {
	UrlId       string `json:"urlId"`
	RedirectUrl string `json:"redirectUrl"`
}

type CreateUrlRequest struct {
	RedirectUrl string `json:"redirectUrl" binding:"required"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type client struct {
	name    string
	service pingCallable
	active  *bool
}

type pingCallable interface {
	Ping(ctx context.Context, in *apipb.PingRequest, opts ...grpc.CallOption) (*apipb.PingResponse, error)
}
