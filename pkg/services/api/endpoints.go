package api

import (
	"fmt"
	"net/http"

	"urlshortener/pkg/logger"
	"urlshortener/pkg/proto/cachepb"
	"urlshortener/pkg/proto/urlspb"

	"github.com/gin-gonic/gin"
)

var (
	BaseUrl                  = "http://localhost:8080"
	UrlIdParam               = "urlId"
	ServicesUnavailableError = ErrorResponse{Error: "services unavailable"}
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
func (service *Service) Pong(ctx *gin.Context) {
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
func (service *Service) GetUrl(ctx *gin.Context) {
	urlId := ctx.Param(UrlIdParam)
	if service.cacheServiceOk {
		c, cancel := makeCtx()
		defer cancel()
		cacheRes, err := service.cacheClient.GetUrl(c, &cachepb.GetUrlRequest{UrlId: urlId})
		if err == nil {
			urlDTO := UrlDTO{
				UrlId:       urlId,
				RedirectUrl: cacheRes.RedirectUrl,
				NewUrl:      fmt.Sprintf("%s/%s", BaseUrl, urlId),
			}
			ctx.JSON(http.StatusOK, urlDTO)
			return
		}
	}
	if service.urlsServiceOk {
		c, cancel := makeCtx()
		defer cancel()
		urlRes, err := service.urlsClient.GetUrl(c, &urlspb.GetUrlRequest{UrlId: urlId})
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
	ctx.JSON(http.StatusInternalServerError, ServicesUnavailableError)
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
func (service *Service) PostUrl(ctx *gin.Context) {
	var body CreateUrlRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	if service.urlsServiceOk {
		c, cancel := makeCtx()
		defer cancel()
		urlRes, err := service.urlsClient.CreateUrl(c, &urlspb.CreateUrlRequest{RedirectUrl: body.RedirectUrl})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		if service.cacheServiceOk {
			c, cancel = makeCtx()
			defer cancel()
			_, err = service.cacheClient.SetUrl(c, &cachepb.SetUrlRequest{Url: &urlspb.Url{UrlId: urlRes.UrlId, RedirectUrl: body.RedirectUrl}})
			if err != nil {
				logger.Error(err)
			}
		}
		urlDTO := UrlDTO{
			UrlId:       urlRes.UrlId,
			RedirectUrl: body.RedirectUrl,
			NewUrl:      fmt.Sprintf("%s/%s", BaseUrl, urlRes.UrlId),
		}
		ctx.JSON(http.StatusOK, urlDTO)
		return
	}
	ctx.JSON(http.StatusInternalServerError, ServicesUnavailableError)
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
func (service *Service) RedirectToUrl(ctx *gin.Context) {
	urlId := ctx.Param(UrlIdParam)
	if service.cacheServiceOk {
		c, cancel := makeCtx()
		defer cancel()
		cacheRes, err := service.cacheClient.GetUrl(c, &cachepb.GetUrlRequest{UrlId: urlId})
		if err == nil {
			ctx.Redirect(http.StatusFound, cacheRes.RedirectUrl)
			return
		}
	}
	if service.urlsServiceOk {
		c, cancel := makeCtx()
		defer cancel()
		urlRes, err := service.urlsClient.GetUrl(c, &urlspb.GetUrlRequest{UrlId: urlId})
		if err == nil {
			ctx.Redirect(http.StatusFound, urlRes.GetUrl().GetRedirectUrl())
		} else {
			ctx.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusInternalServerError, ServicesUnavailableError)
}
