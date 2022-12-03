package api

import (
	"errors"
	"fmt"
	"net/http"

	"urlshortener/pkg/grpcutil"
	"urlshortener/pkg/proto/cachepb"
	"urlshortener/pkg/proto/urlspb"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	BaseUrl                = "http://localhost:8080"
	UrlIdParam             = "urlId"
	ErrServicesUnavailable = errors.New("services unavailable")
	ServicesUnavailable    = ErrorResponse{Error: ErrServicesUnavailable.Error()}
)

// Pong
// @Summary     Ping the server
// @ID          ping
// @Tags        ping
// @Description Ping the server
// @Produce     text/plain
// @Success     200 {string} string "pong"
// @Failure     500 {object} ErrorResponse
// @Router      /ping [get]
func (s *Service) Pong(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}

// GetUrl
// @Summary     Get url
// @ID          get-url
// @Tags        url
// @Description Get url
// @Consume     application/json
// @Produce     application/json
// @Param       urlId path     string true "url id"
// @Success     200   {object} UrlDTO
// @Failure     404   {object} ErrorResponse
// @Failure     500   {object} ErrorResponse
// @Router      /url/{urlId} [get]
func (s *Service) GetUrl(ctx *gin.Context) {
	urlId := ctx.Param(UrlIdParam)
	if s.cacheServiceOk {
		c, cancel := grpcutil.MakeUnaryCtx()
		defer cancel()
		cacheRes, err := s.cacheClient.Get(c, &cachepb.GetRequest{Key: urlId})
		if err == nil {
			urlDTO := UrlDTO{
				UrlId:       urlId,
				RedirectUrl: cacheRes.GetValue(),
				NewUrl:      fmt.Sprintf("%s/%s", BaseUrl, urlId),
			}
			ctx.JSON(http.StatusOK, urlDTO)
			return
		}
	}
	if s.urlsServiceOk {
		c, cancel := grpcutil.MakeUnaryCtx()
		defer cancel()
		urlRes, err := s.urlsClient.GetUrl(c, &urlspb.GetUrlRequest{UrlId: urlId})
		if err == nil {
			urlDTO := UrlDTO{
				UrlId:       urlId,
				RedirectUrl: urlRes.GetUrl().GetRedirectUrl(),
				NewUrl:      fmt.Sprintf("%s/%s", BaseUrl, urlId),
			}
			ctx.JSON(http.StatusOK, urlDTO)
		} else {
			ctx.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusInternalServerError, ServicesUnavailable)
}

// PostUrl
// @Summary     Create a new url redirect
// @ID          create-url
// @Tags        url
// @Description Create a new url redirect
// @Consume     application/json
// @Produce     application/json
// @Param       url body     CreateUrlRequest true "url"
// @Success     200 {object} UrlDTO
// @Failure     400 {object} ErrorResponse
// @Failure     500 {object} ErrorResponse
// @Router      /url [post]
func (s *Service) PostUrl(ctx *gin.Context) {
	var body CreateUrlRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	if s.urlsServiceOk {
		c, cancel := grpcutil.MakeUnaryCtx()
		defer cancel()
		urlRes, err := s.urlsClient.CreateUrl(c, &urlspb.CreateUrlRequest{RedirectUrl: body.RedirectUrl})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		if s.cacheServiceOk {
			c, cancel = grpcutil.MakeUnaryCtx()
			defer cancel()
			_, err = s.cacheClient.Set(c, &cachepb.SetRequest{Key: urlRes.GetUrlId(), Value: body.RedirectUrl})
			if err != nil {
				s.logger.Error("failed to set cache", zap.Error(err))
			}
		}
		urlDTO := UrlDTO{
			UrlId:       urlRes.GetUrlId(),
			RedirectUrl: body.RedirectUrl,
			NewUrl:      fmt.Sprintf("%s/%s", BaseUrl, urlRes.UrlId),
		}
		ctx.JSON(http.StatusOK, urlDTO)
		return
	}
	ctx.JSON(http.StatusInternalServerError, ServicesUnavailable)
}

// RedirectToUrl
// @Summary     Redirect to url
// @ID          redirect-url
// @Tags        url
// @Description Redirect to url
// @Param       urlId path string true "url id"
// @Success     302
// @Failure     404 {object} ErrorResponse
// @Failure     500 {object} ErrorResponse
// @Router      /{urlId} [get]
func (s *Service) RedirectToUrl(ctx *gin.Context) {
	urlId := ctx.Param(UrlIdParam)
	if s.cacheServiceOk {
		c, cancel := grpcutil.MakeUnaryCtx()
		defer cancel()
		cacheRes, err := s.cacheClient.Get(c, &cachepb.GetRequest{Key: urlId})
		if err == nil {
			ctx.Redirect(http.StatusFound, cacheRes.GetValue())
			return
		}
	}
	if s.urlsServiceOk {
		c, cancel := grpcutil.MakeUnaryCtx()
		defer cancel()
		urlRes, err := s.urlsClient.GetUrl(c, &urlspb.GetUrlRequest{UrlId: urlId})
		if err == nil {
			ctx.Redirect(http.StatusFound, urlRes.GetUrl().GetRedirectUrl())
		} else {
			ctx.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusInternalServerError, ServicesUnavailable)
}
