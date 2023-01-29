package api

import (
	"errors"
	"fmt"
	"grpc_util/pkg/grpc_ctx"
	"net/http"
	"proto/pkg/cachepb"
	"proto/pkg/urlspb"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	UrlIdParam             = "urlId"
	ErrServicesUnavailable = errors.New("services unavailable")
)

// Pong
// @Summary     Ping the server
// @ID          ping
// @Tags        ping
// @Description Ping the server
// @Produce     text/plain
// @Success     200 {string} string "pong"
// @Failure     500 {object} ErrorResponse
// @Router      /ping [GET]
func (s *Service) Pong(c *gin.Context) {
	c.String(http.StatusOK, "pong")
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
// @Router      /{urlId} [GET]
func (s *Service) RedirectToUrl(c *gin.Context) {
	s.makeUrlResponse(c, true)
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
// @Router      /url/{urlId} [GET]
func (s *Service) GetUrl(c *gin.Context) {
	s.makeUrlResponse(c, false)
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
// @Router      /url [POST]
func (s *Service) PostUrl(c *gin.Context) {
	var body CreateUrlRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	if s.urlsServiceOk.Load() {
		ctx, cancel := grpc_ctx.MakeUnaryCtx()
		defer cancel()
		urlRes, err := s.urlsClient.CreateUrl(ctx, &urlspb.CreateUrlRequest{RedirectUrl: body.RedirectUrl})
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		if s.cacheServiceOk.Load() {
			ctx, cancel = grpc_ctx.MakeUnaryCtx()
			defer cancel()
			_, err = s.cacheClient.Set(ctx, &cachepb.SetRequest{Key: urlRes.GetUrlId(), Value: body.RedirectUrl})
			if err != nil {
				s.logger.Error("failed to set cache", zap.Error(err))
			}
		}
		urlDTO := UrlDTO{
			UrlId:       urlRes.GetUrlId(),
			RedirectUrl: body.RedirectUrl,
			NewUrl:      fmt.Sprintf("%s/%s", s.addr, urlRes.UrlId),
		}
		c.JSON(http.StatusOK, urlDTO)
		return
	}
	c.JSON(http.StatusInternalServerError, ErrorResponse{Error: ErrServicesUnavailable.Error()})
}

func (s *Service) makeUrlResponse(c *gin.Context, redirect bool) {
	urlId := c.Param(UrlIdParam)
	redirectUrl, err := s.getRedirectUrl(urlId)
	if err != nil {
		code := status.Code(err)
		if code == codes.NotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		}
		return
	}
	if redirect {
		c.Redirect(http.StatusFound, redirectUrl)
	} else {
		urlDTO := UrlDTO{
			UrlId:       urlId,
			RedirectUrl: redirectUrl,
			NewUrl:      fmt.Sprintf("%s/%s", s.addr, urlId),
		}
		c.JSON(http.StatusOK, urlDTO)
	}
}

func (s *Service) getRedirectUrl(urlId string) (string, error) {
	if s.cacheServiceOk.Load() {
		ctx, cancel := grpc_ctx.MakeUnaryCtx()
		defer cancel()
		cacheRes, err := s.cacheClient.Get(ctx, &cachepb.GetRequest{Key: urlId})
		if err == nil {
			return cacheRes.GetValue(), nil
		}
	}
	if s.urlsServiceOk.Load() {
		ctx, cancel := grpc_ctx.MakeUnaryCtx()
		defer cancel()
		urlRes, err := s.urlsClient.GetUrl(ctx, &urlspb.GetUrlRequest{UrlId: urlId})
		if err != nil {
			return "", err
		}
		return urlRes.GetUrl().GetRedirectUrl(), nil
	}
	return "", ErrServicesUnavailable
}
