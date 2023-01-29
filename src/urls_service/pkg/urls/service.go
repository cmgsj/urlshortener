package urls

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"net/url"
	"urls_service/pkg/proto/urlspb"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidUrl       = status.Error(codes.InvalidArgument, "invalid url argument")
	ErrUrlNotFound      = status.Error(codes.NotFound, "url not found")
	ErrUrlAlreadyExists = status.Error(codes.AlreadyExists, "url already exists")
	ErrInternal         = status.Error(codes.Internal, "internal error")
)

type Service struct {
	urlspb.UnimplementedUrlsServiceServer
	healthServer *health.Server
	logger       *zap.Logger
	db           *sql.DB
}

func (s *Service) GetUrl(ctx context.Context, req *urlspb.GetUrlRequest) (*urlspb.GetUrlResponse, error) {
	urlEntity, err := getUrlById(ctx, s.db, req.GetUrlId())
	if err != nil {
		return nil, ErrUrlNotFound
	}
	return &urlspb.GetUrlResponse{Url: &urlspb.Url{UrlId: urlEntity.UrlId, RedirectUrl: urlEntity.RedirectUrl}}, nil
}

func (s *Service) CreateUrl(ctx context.Context, req *urlspb.CreateUrlRequest) (*urlspb.CreateUrlResponse, error) {
	urlId, err := generateUrlId()
	if err != nil {
		return nil, ErrInternal
	}
	if !isValidUrl(req.GetRedirectUrl()) {
		return nil, ErrInvalidUrl
	}
	err = createUrl(ctx, s.db, urlId, req.GetRedirectUrl())
	if err != nil {
		return nil, ErrUrlAlreadyExists
	}
	return &urlspb.CreateUrlResponse{UrlId: urlId}, nil
}

func (s *Service) UpdateUrl(ctx context.Context, req *urlspb.UpdateUrlRequest) (*urlspb.UpdateUrlResponse, error) {
	if !isValidUrl(req.GetUrl().GetRedirectUrl()) {
		return nil, ErrInvalidUrl
	}
	err := updateUrl(ctx, s.db, req.GetUrl().GetUrlId(), req.GetUrl().GetRedirectUrl())
	if err != nil {
		return nil, ErrUrlAlreadyExists
	}
	return &urlspb.UpdateUrlResponse{}, nil
}

func (s *Service) DeleteUrl(ctx context.Context, req *urlspb.DeleteUrlRequest) (*urlspb.DeleteUrlResponse, error) {
	err := deleteUrl(ctx, s.db, req.GetUrlId())
	if err != nil {
		return nil, ErrUrlNotFound
	}
	return &urlspb.DeleteUrlResponse{}, nil
}

func generateUrlId() (string, error) {
	var data [8]byte
	if _, err := rand.Read(data[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data[:]), nil
}

func isValidUrl(urlStr string) bool {
	_, err := url.ParseRequestURI(urlStr)
	return err == nil
}
