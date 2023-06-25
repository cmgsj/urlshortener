package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"

	"github.com/cmgsj/urlshortener/pkg/database"
	urlshortenerv1 "github.com/cmgsj/urlshortener/pkg/gen/proto/urlshortener/v1"
	urlsv1 "github.com/cmgsj/urlshortener/pkg/gen/sqlc/urls/v1"
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
	urlshortenerv1.UnimplementedURLShortenerServer
	DB database.DB
}

func (s *Service) ListURLs(ctx context.Context, req *urlshortenerv1.ListURLsRequest) (*urlshortenerv1.ListURLsResponse, error) {
	urls, err := s.DB.ListUrls(ctx)
	if err != nil {
		return nil, ErrInternal
	}
	urlsv1 := make([]*urlshortenerv1.URL, len(urls))
	for i, u := range urls {
		urlsv1[i] = &urlshortenerv1.URL{
			UrlId:       u.UrlID,
			RedirectURL: u.RedirectUrl,
		}
	}
	return &urlshortenerv1.ListURLsResponse{Urls: urlsv1}, nil
}

func (s *Service) GetURL(ctx context.Context, req *urlshortenerv1.GetURLRequest) (*urlshortenerv1.GetURLResponse, error) {
	u, err := s.DB.GetUrl(ctx, req.GetUrlId())
	if err != nil {
		return nil, ErrUrlNotFound
	}
	return &urlshortenerv1.GetURLResponse{Url: &urlshortenerv1.URL{UrlId: u.UrlID, RedirectURL: u.RedirectUrl}}, nil
}

func (s *Service) CreateURL(ctx context.Context, req *urlshortenerv1.CreateURLRequest) (*urlshortenerv1.CreateURLResponse, error) {
	urlId, err := generateUrlId(UrlIdLength)
	if err != nil {
		return nil, ErrInternal
	}
	if !isValidUrl(req.GetRedirectURL()) {
		return nil, ErrInvalidUrl
	}
	u, err := s.DB.CreateUrl(ctx, urlsv1.CreateUrlParams{
		UrlID:       urlId,
		RedirectUrl: req.GetRedirectURL(),
	})
	if err != nil {
		return nil, ErrUrlAlreadyExists
	}
	return &urlshortenerv1.CreateURLResponse{UrlId: u.UrlID}, nil
}

func (s *Service) UpdateURL(ctx context.Context, req *urlshortenerv1.UpdateURLRequest) (*urlshortenerv1.UpdateURLResponse, error) {
	if !isValidUrl(req.GetUrl().GetRedirectURL()) {
		return nil, ErrInvalidUrl
	}
	if err := s.DB.UpdateUrl(ctx, urlsv1.UpdateUrlParams{
		UrlID:       req.GetUrl().GetUrlId(),
		RedirectUrl: req.GetUrl().GetRedirectURL(),
	}); err != nil {
		return nil, ErrUrlAlreadyExists
	}
	return &urlshortenerv1.UpdateURLResponse{}, nil
}

func (s *Service) DeleteURL(ctx context.Context, req *urlshortenerv1.DeleteURLRequest) (*urlshortenerv1.DeleteURLResponse, error) {
	if err := s.DB.DeleteUrl(ctx, req.GetUrlId()); err != nil {
		return nil, ErrUrlNotFound
	}
	return &urlshortenerv1.DeleteURLResponse{}, nil
}

func (s *Service) RedirectURL(prefix string) http.Handler {
	return http.StripPrefix(prefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := s.DB.GetUrl(r.Context(), strings.TrimPrefix(r.URL.Path, "/"))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, u.RedirectUrl, http.StatusMovedPermanently)
	}))
}

func (s *Service) SeedDB(ctx context.Context) error {
	params := []urlsv1.CreateUrlParams{
		{UrlID: "google", RedirectUrl: "https://www.google.com"},
		{UrlID: "youtube", RedirectUrl: "https://www.youtube.com"},
		{UrlID: "apple", RedirectUrl: "https://www.apple.com"},
	}
	for _, param := range params {
		if _, err := s.DB.CreateUrl(ctx, param); err != nil {
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
