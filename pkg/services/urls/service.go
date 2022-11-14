package urls

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"net/url"
	"time"
	"urlshortener/pkg/logger"
	"urlshortener/pkg/proto/healthpb"
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
	healthpb.UnimplementedHealthServiceServer
	db *sql.DB
}

func (service *Service) GetUrl(ctx context.Context, req *urlspb.GetUrlRequest) (*urlspb.GetUrlResponse, error) {
	urlEntity, err := getUrl(service.db, ctx, req.GetUrlId())
	if err != nil {
		return nil, UrlNotFoundError
	}
	return &urlspb.GetUrlResponse{Url: &urlspb.Url{UrlId: urlEntity.UrlId, RedirectUrl: urlEntity.RedirectUrl}}, nil
}

func (service *Service) CreateUrl(ctx context.Context, req *urlspb.CreateUrlRequest) (*urlspb.CreateUrlResponse, error) {
	urlId, err := generateID()
	if err != nil {
		return nil, InternalServerError
	}
	if !validateUrl(req.GetRedirectUrl()) {
		return nil, InvalidUrlError
	}
	err = createUrl(service.db, ctx, urlId, req.GetRedirectUrl())
	if err != nil {
		return nil, UrlAlreadyExistsError
	}
	return &urlspb.CreateUrlResponse{UrlId: urlId}, nil
}

func (service *Service) UpdateUrl(ctx context.Context, req *urlspb.UpdateUrlRequest) (*urlspb.NoContent, error) {
	if !validateUrl(req.GetUrl().GetRedirectUrl()) {
		return nil, InvalidUrlError
	}
	err := updateUrl(service.db, ctx, req.GetUrl().GetUrlId(), req.GetUrl().GetRedirectUrl())
	if err != nil {
		return nil, UrlAlreadyExistsError
	}
	return &urlspb.NoContent{}, nil
}

func (service *Service) DeleteUrl(ctx context.Context, req *urlspb.DeleteUrlRequest) (*urlspb.NoContent, error) {
	err := deleteUrl(service.db, ctx, req.GetUrlId())
	if err != nil {
		return nil, UrlNotFoundError
	}
	return &urlspb.NoContent{}, nil
}

func (service *Service) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (service *Service) Watch(req *healthpb.HealthCheckRequest, stream healthpb.HealthService_WatchServer) error {
	for {
		select {
		case <-stream.Context().Done():
			return nil
		default:
			err := stream.Send(&healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING})
			if err != nil {
				return err
			}
			time.Sleep(time.Minute)
		}
	}
}

func (service *Service) seedUrls() {
	redirectUrls := []string{
		"https://www.google.com",
		"https://www.youtube.com",
		"https://www.wikipedia.org",
		"https://www.reddit.com",
		"https://www.amazon.com",
	}
	for _, redirectUrl := range redirectUrls {
		res, err := service.CreateUrl(context.Background(), &urlspb.CreateUrlRequest{RedirectUrl: redirectUrl})
		if err != nil && err != UrlAlreadyExistsError {
			logger.Error("failed to seed url:", redirectUrl, "error:", err)
		} else if err == UrlAlreadyExistsError {
			urlEntity, err := getUrlByRedirectUrl(service.db, context.Background(), redirectUrl)
			if err != nil {
				logger.Error("failed to get url id for url:", redirectUrl, "error:", err)
			} else {
				logger.Info("url: id=", urlEntity.UrlId, "url=", urlEntity.RedirectUrl)
			}
		} else {
			logger.Info("seeded url: id=", res.GetUrlId(), "url=", redirectUrl)
		}
	}
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
