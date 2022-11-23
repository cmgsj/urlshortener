package urls

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"net/url"
	"urlshortener/pkg/proto/urlspb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InvalidUrlError       = status.Error(codes.InvalidArgument, "invalid url argument")
	UrlNotFoundError      = status.Error(codes.NotFound, "url not found")
	UrlAlreadyExistsError = status.Error(codes.AlreadyExists, "url already exists")
	InternalServerError   = status.Error(codes.Internal, "internal server error")
)

type Service struct {
	urlspb.UnimplementedUrlsServiceServer
	db *sql.DB
}

func (s *Service) GetUrl(ctx context.Context, req *urlspb.GetUrlRequest) (*urlspb.GetUrlResponse, error) {
	urlEntity, err := getUrlById(ctx, s.db, req.GetUrlId())
	if err != nil {
		return nil, UrlNotFoundError
	}
	return &urlspb.GetUrlResponse{Url: &urlspb.Url{UrlId: urlEntity.UrlId, RedirectUrl: urlEntity.RedirectUrl}}, nil
}

func (s *Service) CreateUrl(ctx context.Context, req *urlspb.CreateUrlRequest) (*urlspb.CreateUrlResponse, error) {
	urlId, err := generateID()
	if err != nil {
		return nil, InternalServerError
	}
	if !validateUrl(req.GetRedirectUrl()) {
		return nil, InvalidUrlError
	}
	var userId int64 = 1 //TODO: get user id from auth service
	err = createUrl(ctx, s.db, urlId, req.GetRedirectUrl(), userId)
	if err != nil {
		return nil, UrlAlreadyExistsError
	}
	return &urlspb.CreateUrlResponse{UrlId: urlId}, nil
}

func (s *Service) UpdateUrl(ctx context.Context, req *urlspb.UpdateUrlRequest) (*urlspb.NoContent, error) {
	if !validateUrl(req.GetUrl().GetRedirectUrl()) {
		return nil, InvalidUrlError
	}
	var userId int64 = 1 //TODO: get user id from auth service
	err := updateUrl(ctx, s.db, req.GetUrl().GetUrlId(), req.GetUrl().GetRedirectUrl(), userId)
	if err != nil {
		return nil, UrlAlreadyExistsError
	}
	return &urlspb.NoContent{}, nil
}

func (s *Service) DeleteUrl(ctx context.Context, req *urlspb.DeleteUrlRequest) (*urlspb.NoContent, error) {
	var userId int64 = 1 //TODO: get user id from auth service
	err := deleteUrl(ctx, s.db, req.GetUrlId(), userId)
	if err != nil {
		return nil, UrlNotFoundError
	}
	return &urlspb.NoContent{}, nil
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
