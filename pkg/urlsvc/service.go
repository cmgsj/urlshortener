package urlsvc

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"net/url"

	"github.com/cmgsj/urlshortener/pkg/proto/urlpb"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidUrl       = status.Error(codes.InvalidArgument, "invalid url argument")
	ErrUrlNotFound      = status.Error(codes.NotFound, "url not found")
	ErrUrlAlreadyExists = status.Error(codes.AlreadyExists, "url already exists")
	ErrInternal         = status.Error(codes.Internal, "internal error")
)

var (
	UrlIdLength = 8
)

type Service struct {
	urlpb.UnimplementedUrlServiceServer
	Logger *zap.Logger
	Db     *sql.DB
}

func (s *Service) GetUrl(ctx context.Context, req *urlpb.GetUrlRequest) (*urlpb.GetUrlResponse, error) {
	urlEntity, err := getUrlById(ctx, s.Db, req.GetUrlId())
	if err != nil {
		return nil, ErrUrlNotFound
	}
	return &urlpb.GetUrlResponse{Url: &urlpb.Url{UrlId: urlEntity.UrlId, RedirectUrl: urlEntity.RedirectUrl}}, nil
}

func (s *Service) CreateUrl(ctx context.Context, req *urlpb.CreateUrlRequest) (*urlpb.CreateUrlResponse, error) {
	urlId, err := generateUrlId(UrlIdLength)
	if err != nil {
		return nil, ErrInternal
	}
	if !isValidUrl(req.GetRedirectUrl()) {
		return nil, ErrInvalidUrl
	}
	err = createUrl(ctx, s.Db, urlId, req.GetRedirectUrl())
	if err != nil {
		return nil, ErrUrlAlreadyExists
	}
	return &urlpb.CreateUrlResponse{UrlId: urlId}, nil
}

func (s *Service) UpdateUrl(ctx context.Context, req *urlpb.UpdateUrlRequest) (*urlpb.UpdateUrlResponse, error) {
	if !isValidUrl(req.GetUrl().GetRedirectUrl()) {
		return nil, ErrInvalidUrl
	}
	err := updateUrl(ctx, s.Db, req.GetUrl().GetUrlId(), req.GetUrl().GetRedirectUrl())
	if err != nil {
		return nil, ErrUrlAlreadyExists
	}
	return &urlpb.UpdateUrlResponse{}, nil
}

func (s *Service) DeleteUrl(ctx context.Context, req *urlpb.DeleteUrlRequest) (*urlpb.DeleteUrlResponse, error) {
	err := deleteUrl(ctx, s.Db, req.GetUrlId())
	if err != nil {
		return nil, ErrUrlNotFound
	}
	return &urlpb.DeleteUrlResponse{}, nil
}

func generateUrlId(n int) (string, error) {
	data := make([]byte, n)
	if _, err := rand.Read(data); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}

func isValidUrl(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}
