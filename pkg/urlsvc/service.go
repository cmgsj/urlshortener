package urlsvc

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/url"

	"github.com/cmgsj/urlshortener/pkg/proto/urlpb"
	"github.com/cmgsj/urlshortener/pkg/urlsvc/db"
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
	urlpb.UnimplementedUrlServiceServer
	Logger *zap.Logger
	Query  db.Querier
}

func New(logger *zap.Logger, query db.Querier) *Service {
	return &Service{
		Logger: logger,
		Query:  query,
	}
}

func (s *Service) GetUrl(ctx context.Context, req *urlpb.GetUrlRequest) (*urlpb.GetUrlResponse, error) {
	u, err := s.Query.GetUrl(ctx, req.GetUrlId())
	if err != nil {
		return nil, ErrUrlNotFound
	}
	return &urlpb.GetUrlResponse{Url: &urlpb.Url{UrlId: u.UrlID, RedirectUrl: u.RedirectUrl}}, nil
}

func (s *Service) CreateUrl(ctx context.Context, req *urlpb.CreateUrlRequest) (*urlpb.CreateUrlResponse, error) {
	urlId, err := generateUrlId(UrlIdLength)
	if err != nil {
		return nil, ErrInternal
	}
	if !isValidUrl(req.GetRedirectUrl()) {
		return nil, ErrInvalidUrl
	}
	arg := db.CreateUrlParams{
		UrlID:       urlId,
		RedirectUrl: req.GetRedirectUrl(),
	}
	u, err := s.Query.CreateUrl(ctx, arg)
	if err != nil {
		return nil, ErrUrlAlreadyExists
	}
	return &urlpb.CreateUrlResponse{UrlId: u.UrlID}, nil
}

func (s *Service) UpdateUrl(ctx context.Context, req *urlpb.UpdateUrlRequest) (*urlpb.UpdateUrlResponse, error) {
	if !isValidUrl(req.GetUrl().GetRedirectUrl()) {
		return nil, ErrInvalidUrl
	}
	arg := db.UpdateUrlParams{
		UrlID:       req.GetUrl().GetUrlId(),
		RedirectUrl: req.GetUrl().GetRedirectUrl(),
	}
	if err := s.Query.UpdateUrl(ctx, arg); err != nil {
		return nil, ErrUrlAlreadyExists
	}
	return &urlpb.UpdateUrlResponse{}, nil
}

func (s *Service) DeleteUrl(ctx context.Context, req *urlpb.DeleteUrlRequest) (*urlpb.DeleteUrlResponse, error) {
	if err := s.Query.DeleteUrl(ctx, req.GetUrlId()); err != nil {
		return nil, ErrUrlNotFound
	}
	return &urlpb.DeleteUrlResponse{}, nil
}

func (s *Service) SeedDB(ctx context.Context) error {
	args := []db.CreateUrlParams{
		{UrlID: "abcdef01", RedirectUrl: "https://www.google.com"},
		{UrlID: "abcdef02", RedirectUrl: "https://www.youtube.com"},
		{UrlID: "abcdef03", RedirectUrl: "https://www.apple.com"},
	}
	for _, arg := range args {
		if _, err := s.Query.CreateUrl(ctx, arg); err != nil {
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
