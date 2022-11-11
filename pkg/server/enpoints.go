package server

import (
	"log"
	"net/http"
	"urlshortener/pkg/cache"
	"urlshortener/pkg/urls"

	"github.com/gin-gonic/gin"
)

type CreateUrlDTO struct {
	RedirectUrl string `json:"redirectUrl" binding:"required"`
}

type UrlDTO struct {
	UrlId       string `json:"urlId"`
	RedirectUrl string `json:"redirectUrl"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// @Summary     Ping the server
// @ID          ping
// @Tags        ping
// @Description Ping the server
// @Produce     text/plain
// @Success     200 {string} string "pong"
// @Failure     500 {object} ErrorResponse
// @Router      /ping [get]
func (server *Server) Ping(ctx *gin.Context) {
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
func (server *Server) GetUrl(ctx *gin.Context) {
	urlId := ctx.Param("urlId")
	c, cancel := makeCtx()
	defer cancel()
	cacheRes, err := server.cacheClient.GetUrl(c, &cache.GetUrlRequest{UrlId: urlId})
	if err == nil {
		ctx.JSON(http.StatusOK, UrlDTO{UrlId: urlId, RedirectUrl: cacheRes.RedirectUrl})
		return
	}
	c, cancel = makeCtx()
	defer cancel()
	urlRes, err := server.urlClient.GetUrl(c, &urls.GetUrlRequest{UrlId: urlId})
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, UrlDTO{UrlId: urlId, RedirectUrl: urlRes.RedirectUrl})
}

// @Summary     Create a new url
// @ID          create-url
// @Tags        url
// @Description Create a new url
// @Consume     application/json
// @Produce     application/json
// @Param       url body     CreateUrlDTO true "url"
// @Success     200 {object} UrlDTO
// @Failure     400 {object} ErrorResponse
// @Failure     500 {object} ErrorResponse
// @Router      /url [post]
func (server *Server) PostUrl(ctx *gin.Context) {
	var body CreateUrlDTO
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c, cancel := makeCtx()
	defer cancel()
	urlRes, err := server.urlClient.CreateUrl(c, &urls.CreateUrlRequest{RedirectUrl: body.RedirectUrl})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c, cancel = makeCtx()
	defer cancel()
	_, err = server.cacheClient.SetUrl(c, &cache.SetUrlRequest{UrlId: urlRes.UrlId, RedirectUrl: body.RedirectUrl})
	if err != nil {
		log.Println(err)
	}
	ctx.JSON(http.StatusOK, UrlDTO{UrlId: urlRes.UrlId, RedirectUrl: body.RedirectUrl})
}

func (server *Server) RedirectToUrl(ctx *gin.Context) {
	urlId := ctx.Param("urlId")
	c, cancel := makeCtx()
	defer cancel()
	cacheRes, err := server.cacheClient.GetUrl(c, &cache.GetUrlRequest{UrlId: urlId})
	if err == nil {
		ctx.Redirect(http.StatusTemporaryRedirect, cacheRes.RedirectUrl)
		return
	}
	c, cancel = makeCtx()
	defer cancel()
	urlRes, err := server.urlClient.GetUrl(c, &urls.GetUrlRequest{UrlId: urlId})
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	ctx.Redirect(http.StatusTemporaryRedirect, urlRes.RedirectUrl)
}
