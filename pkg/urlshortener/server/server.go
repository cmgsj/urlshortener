package server

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
)

type Server struct {
	urlshortenerv1.UnimplementedURLShortenerServiceServer
	db database.DB
}

func NewServer(db database.DB) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) ListURLs(ctx context.Context, req *urlshortenerv1.ListURLsRequest) (*urlshortenerv1.ListURLsResponse, error) {
	urlList, err := s.db.ListUrls(ctx)
	if err != nil {
		return nil, ErrInternal
	}

	urls := make([]*urlshortenerv1.URL, len(urlList))

	for i, u := range urlList {
		urls[i] = &urlshortenerv1.URL{
			UrlId:       u.UrlID,
			RedirectUrl: u.RedirectUrl,
		}
	}

	return &urlshortenerv1.ListURLsResponse{
		Urls: urls,
	}, nil
}

func (s *Server) GetURL(ctx context.Context, req *urlshortenerv1.GetURLRequest) (*urlshortenerv1.GetURLResponse, error) {
	u, err := s.db.GetUrl(ctx, req.GetUrlId())
	if err != nil {
		return nil, ErrUrlNotFound
	}

	return &urlshortenerv1.GetURLResponse{
		Url: &urlshortenerv1.URL{
			UrlId:       u.UrlID,
			RedirectUrl: u.RedirectUrl,
		},
	}, nil
}

func (s *Server) CreateURL(ctx context.Context, req *urlshortenerv1.CreateURLRequest) (*urlshortenerv1.CreateURLResponse, error) {
	urlId, err := generateUrlId()
	if err != nil {
		return nil, ErrInternal
	}

	if !isValidUrl(req.GetRedirectUrl()) {
		return nil, ErrInvalidUrl
	}

	u, err := s.db.CreateUrl(ctx, urlsv1.CreateUrlParams{
		UrlID:       urlId,
		RedirectUrl: req.GetRedirectUrl(),
	})
	if err != nil {
		return nil, ErrUrlAlreadyExists
	}

	return &urlshortenerv1.CreateURLResponse{
		UrlId: u.UrlID,
	}, nil
}

func (s *Server) UpdateURL(ctx context.Context, req *urlshortenerv1.UpdateURLRequest) (*urlshortenerv1.UpdateURLResponse, error) {
	if !isValidUrl(req.GetUrl().GetRedirectUrl()) {
		return nil, ErrInvalidUrl
	}

	err := s.db.UpdateUrl(ctx, urlsv1.UpdateUrlParams{
		UrlID:       req.GetUrl().GetUrlId(),
		RedirectUrl: req.GetUrl().GetRedirectUrl(),
	})
	if err != nil {
		return nil, ErrUrlAlreadyExists
	}

	return &urlshortenerv1.UpdateURLResponse{}, nil
}

func (s *Server) DeleteURL(ctx context.Context, req *urlshortenerv1.DeleteURLRequest) (*urlshortenerv1.DeleteURLResponse, error) {
	err := s.db.DeleteUrl(ctx, req.GetUrlId())
	if err != nil {
		return nil, ErrUrlNotFound
	}

	return &urlshortenerv1.DeleteURLResponse{}, nil
}

func (s *Server) RedirectUrl() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := s.db.GetUrl(r.Context(), strings.TrimPrefix(r.URL.Path, "/"))
		if err != nil {
			http.NotFound(w, r)
			return
		}

		http.Redirect(w, r, u.RedirectUrl, http.StatusMovedPermanently)
	})
}

func (s *Server) SeedDB(ctx context.Context) error {
	params := []urlsv1.CreateUrlParams{
		{UrlID: "google", RedirectUrl: "https://www.google.com"},
		{UrlID: "youtube", RedirectUrl: "https://www.youtube.com"},
		{UrlID: "apple", RedirectUrl: "https://www.apple.com"},
	}

	for _, param := range params {
		_, err := s.db.CreateUrl(ctx, param)
		if err != nil {
			return err
		}
	}

	return nil
}

const urlIdLength = 8

func generateUrlId() (string, error) {
	data := make([]byte, urlIdLength)

	if _, err := rand.Read(data); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}

func isValidUrl(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}