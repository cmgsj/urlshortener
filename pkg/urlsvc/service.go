package urlsvc

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/url"

	urlv1 "github.com/cmgsj/urlshortener/pkg/gen/proto/url/v1"
	sqlc "github.com/cmgsj/urlshortener/pkg/gen/sqlc/url/v1"
	"github.com/cmgsj/urlshortener/pkg/urlsvc/database"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidUrl       = status.Error(codes.InvalidArgument, "invalid url argument")
	ErrUrlNotFound      = status.Error(codes.NotFound, "url not found")
	ErrUrlAlreadyExists = status.Error(codes.AlreadyExists, "url already exists")
	ErrInternal         = status.Error(codes.Internal, "internal error")

	UrlIdLength = 8
)

type Service struct {
	urlv1.UnimplementedUrlServiceServer
	Logger *zap.Logger
	DB     *database.DB
}

func New(logger *zap.Logger, db *database.DB) *Service {
	return &Service{
		Logger: logger,
		DB:     db,
	}
}

func (s *Service) GetUrl(ctx context.Context, req *urlv1.GetUrlRequest) (*urlv1.GetUrlResponse, error) {
	u, err := s.DB.Query.GetUrl(ctx, req.GetUrlId())
	if err != nil {
		return nil, ErrUrlNotFound
	}
	return &urlv1.GetUrlResponse{Url: &urlv1.Url{UrlId: u.UrlID, RedirectUrl: u.RedirectUrl}}, nil
}

func (s *Service) CreateUrl(ctx context.Context, req *urlv1.CreateUrlRequest) (*urlv1.CreateUrlResponse, error) {
	urlId, err := generateUrlId(UrlIdLength)
	if err != nil {
		return nil, ErrInternal
	}
	if !isValidUrl(req.GetRedirectUrl()) {
		return nil, ErrInvalidUrl
	}
	param := &sqlc.CreateUrlParams{
		UrlID:       urlId,
		RedirectUrl: req.GetRedirectUrl(),
	}
	u, err := s.DB.Query.CreateUrl(ctx, param)
	if err != nil {
		return nil, ErrUrlAlreadyExists
	}
	return &urlv1.CreateUrlResponse{UrlId: u.UrlID}, nil
}

func (s *Service) UpdateUrl(ctx context.Context, req *urlv1.UpdateUrlRequest) (*urlv1.UpdateUrlResponse, error) {
	if !isValidUrl(req.GetUrl().GetRedirectUrl()) {
		return nil, ErrInvalidUrl
	}
	param := &sqlc.UpdateUrlParams{
		UrlID:       req.GetUrl().GetUrlId(),
		RedirectUrl: req.GetUrl().GetRedirectUrl(),
	}
	if err := s.DB.Query.UpdateUrl(ctx, param); err != nil {
		return nil, ErrUrlAlreadyExists
	}
	return &urlv1.UpdateUrlResponse{}, nil
}

func (s *Service) DeleteUrl(ctx context.Context, req *urlv1.DeleteUrlRequest) (*urlv1.DeleteUrlResponse, error) {
	if err := s.DB.Query.DeleteUrl(ctx, req.GetUrlId()); err != nil {
		return nil, ErrUrlNotFound
	}
	return &urlv1.DeleteUrlResponse{}, nil
}

func (s *Service) SeedDB(ctx context.Context) error {
	params := []*sqlc.CreateUrlParams{
		{UrlID: "abcdef01", RedirectUrl: "https://www.google.com"},
		{UrlID: "abcdef02", RedirectUrl: "https://www.youtube.com"},
		{UrlID: "abcdef03", RedirectUrl: "https://www.apple.com"},
	}
	for _, param := range params {
		if _, err := s.DB.Query.CreateUrl(ctx, param); err != nil {
			return err
		}
	}
	return nil
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
