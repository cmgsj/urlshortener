package urls

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"net/url"
	"urlshortener/pkg/protobuf/apipb"
	"urlshortener/pkg/protobuf/urlspb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InvalidUrlArgumentError = status.Error(codes.InvalidArgument, "invalid url argument")
	UrlNotFoundError        = status.Error(codes.NotFound, "url not found")
	UrlAlreadyExistsError   = status.Error(codes.AlreadyExists, "url already exists")
	InternalServerError     = status.Error(codes.Internal, "internal server error")
)

type urlServer struct {
	urlspb.UnimplementedUrlsServer
	db *sql.DB
}

func (server *urlServer) GetUrl(ctx context.Context, req *urlspb.GetUrlRequest) (*urlspb.GetUrlResponse, error) {
	urlEntity, err := getUrl(server.db, ctx, req.GetUrlId())
	if err != nil {
		return nil, UrlNotFoundError
	}
	return &urlspb.GetUrlResponse{UrlId: urlEntity.UrlId, RedirectUrl: urlEntity.RedirectUrl}, nil
}

func (server *urlServer) CreateUrl(ctx context.Context, req *urlspb.CreateUrlRequest) (*urlspb.CreateUrlResponse, error) {
	urlId, err := generateID()
	if err != nil {
		return nil, InternalServerError
	}
	if !validateUrl(req.GetRedirectUrl()) {
		return nil, InvalidUrlArgumentError
	}
	err = createUrl(server.db, ctx, urlId, req.GetRedirectUrl())
	if err != nil {
		return nil, UrlAlreadyExistsError
	}
	return &urlspb.CreateUrlResponse{UrlId: urlId}, nil
}

func (server *urlServer) UpdateUrl(ctx context.Context, req *urlspb.UpdateUrlRequest) (*apipb.NoContent, error) {
	if !validateUrl(req.GetRedirectUrl()) {
		return nil, InvalidUrlArgumentError
	}
	err := updateUrl(server.db, ctx, req.GetUrlId(), req.GetRedirectUrl())
	if err != nil {
		return nil, UrlAlreadyExistsError
	}
	return &apipb.NoContent{}, nil
}

func (server *urlServer) DeleteUrl(ctx context.Context, req *urlspb.DeleteUrlRequest) (*apipb.NoContent, error) {
	err := deleteUrl(server.db, ctx, req.GetUrlId())
	if err != nil {
		return nil, UrlNotFoundError
	}
	return &apipb.NoContent{}, nil
}

func (server *urlServer) Ping(ctx context.Context, req *apipb.PingRequest) (*apipb.PingResponse, error) {
	return &apipb.PingResponse{Message: "pong"}, nil
}

func generateID() (string, error) {
	var data [8]byte
	if _, err := rand.Read(data[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data[:]), nil
}

func validateUrl(urlStr string) bool {
	_, err := url.ParseRequestURI(urlStr)
	return err == nil
}
