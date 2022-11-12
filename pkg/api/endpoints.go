package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"urlshortener/pkg/protobuf/cachepb"
	"urlshortener/pkg/protobuf/urlspb"

	"github.com/gin-gonic/gin"
)

var (
	UrlIdParam               = "urlId"
	ServicesUnavailableError = ErrorResponse{Error: "services unavailable"}
)

// @Summary     Ping the server
// @ID          ping
// @Tags        ping
// @Description Ping the server
// @Produce     text/plain
// @Success     200 {string} string "pong"
// @Failure     500 {object} ErrorResponse
// @Router      /ping [get]
func (server *apiServer) Pong(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong")
}

// @Summary     Get the redirect url
// @ID          get-url
// @Tags        url
// @Description Get the redirect url
// @Consume     application/json
// @Produce     application/json
// @Param       urlId path     string true "url id"
// @Success     200   {object} UrlDTO
// @Failure     404   {object} ErrorResponse
// @Failure     500   {object} ErrorResponse
// @Router      /url/{urlId} [get]
func (server *apiServer) GetUrl(ctx *gin.Context) {
	urlId := ctx.Param(UrlIdParam)
	if server.cacheServiceActive {
		c, cancel := makeCtx()
		defer cancel()
		cacheRes, err := server.cacheClient.GetUrl(c, &cachepb.GetUrlRequest{UrlId: urlId})
		if err == nil {
			ctx.JSON(http.StatusOK, UrlDTO{UrlId: urlId, RedirectUrl: cacheRes.RedirectUrl})
			return
		}
	}
	if server.urlsServiceActive {
		c, cancel := makeCtx()
		defer cancel()
		urlRes, err := server.urlsClient.GetUrl(c, &urlspb.GetUrlRequest{UrlId: urlId})
		if err == nil {
			ctx.JSON(http.StatusOK, UrlDTO{UrlId: urlId, RedirectUrl: urlRes.RedirectUrl})
		} else {
			ctx.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusInternalServerError, ServicesUnavailableError)
}

// @Summary     Create a new url
// @ID          create-url
// @Tags        url
// @Description Create a new url
// @Consume     application/json
// @Produce     application/json
// @Param       url body     CreateUrlRequest true "url"
// @Success     200 {object} UrlDTO
// @Failure     400 {object} ErrorResponse
// @Failure     500 {object} ErrorResponse
// @Router      /url [post]
func (server *apiServer) PostUrl(ctx *gin.Context) {
	var body CreateUrlRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	if server.urlsServiceActive {
		c, cancel := makeCtx()
		defer cancel()
		urlRes, err := server.urlsClient.CreateUrl(c, &urlspb.CreateUrlRequest{RedirectUrl: body.RedirectUrl})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
		if server.cacheServiceActive {
			c, cancel = makeCtx()
			defer cancel()
			_, err = server.cacheClient.SetUrl(c, &cachepb.SetUrlRequest{UrlId: urlRes.UrlId, RedirectUrl: body.RedirectUrl})
			if err != nil {
				log.Println(err)
			}
		}
		ctx.JSON(http.StatusOK, UrlDTO{UrlId: urlRes.UrlId, RedirectUrl: body.RedirectUrl})
		return
	}
	ctx.JSON(http.StatusInternalServerError, ServicesUnavailableError)
}

func (server *apiServer) RedirectToUrl(ctx *gin.Context) {
	urlId := ctx.Param(UrlIdParam)
	if server.cacheServiceActive {
		c, cancel := makeCtx()
		defer cancel()
		cacheRes, err := server.cacheClient.GetUrl(c, &cachepb.GetUrlRequest{UrlId: urlId})
		if err == nil {
			ctx.Redirect(http.StatusTemporaryRedirect, cacheRes.RedirectUrl)
			return
		}
	}
	if server.urlsServiceActive {
		c, cancel := makeCtx()
		defer cancel()
		urlRes, err := server.urlsClient.GetUrl(c, &urlspb.GetUrlRequest{UrlId: urlId})
		if err == nil {
			ctx.Redirect(http.StatusTemporaryRedirect, urlRes.RedirectUrl)
		} else {
			ctx.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusInternalServerError, ServicesUnavailableError)
}

func makeCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second)
}
