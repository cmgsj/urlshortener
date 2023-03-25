package websvc

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	urlv1 "github.com/cmgsj/urlshortener/pkg/gen/proto/url/v1"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	UrlIdParam             = "urlId"
	ErrServicesUnavailable = errors.New("services unavailable")
)

// Pong
//
//	@Summary		Ping the server
//	@ID				ping
//	@Tags			ping
//	@Description	Ping the server
//	@Produce		text/plain
//	@Success		200	{string}	string	"pong"
//	@Failure		500	{object}	ErrorResponse
//	@Router			/ping [GET]
func (s *Service) Pong(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// RedirectToUrl
//
//	@Summary		Redirect to url
//	@ID				redirect-url
//	@Tags			url
//	@Description	Redirect to url
//	@Param			urlId	path	string	true	"url id"
//	@Success		302
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/{urlId} [GET]
func (s *Service) RedirectToUrl(c *gin.Context) {
	s.makeGetUrlResponse(c, true)
}

// GetUrl
//
//	@Summary		Get url
//	@ID				get-url
//	@Tags			url
//	@Description	Get url
//	@Consume		application/json
//	@Produce		application/json
//	@Param			urlId	path		string	true	"url id"
//	@Success		200		{object}	UrlDto
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/url/{urlId} [GET]
func (s *Service) GetUrl(c *gin.Context) {
	s.makeGetUrlResponse(c, false)
}

// PostUrl
//
//	@Summary		Create a new url redirect
//	@ID				create-url
//	@Tags			url
//	@Description	Create a new url redirect
//	@Consume		application/json
//	@Produce		application/json
//	@Param			url	body		CreateUrlRequest	true	"url"
//	@Success		200	{object}	UrlDto
//	@Failure		400	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/url [POST]
func (s *Service) PostUrl(c *gin.Context) {
	var body CreateUrlRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	urlRes, err := s.UrlClient.CreateUrl(ctx, &urlv1.CreateUrlRequest{RedirectUrl: body.RedirectUrl})
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	s.putInCache(urlRes.GetUrlId(), body.RedirectUrl)
	UrlDto := UrlDto{
		UrlId:       urlRes.GetUrlId(),
		RedirectUrl: body.RedirectUrl,
		NewUrl:      fmt.Sprintf("%s/%s", s.Addr, urlRes.UrlId),
	}
	c.JSON(http.StatusOK, UrlDto)
}

func (s *Service) makeGetUrlResponse(c *gin.Context, redirect bool) {
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
		UrlDto := UrlDto{
			UrlId:       urlId,
			RedirectUrl: redirectUrl,
			NewUrl:      fmt.Sprintf("%s/%s", s.Addr, urlId),
		}
		c.JSON(http.StatusOK, UrlDto)
	}
}

func (s *Service) getRedirectUrl(urlId string) (string, error) {
	if redirectUrl, err := s.getFromCache(urlId); err == nil {
		return redirectUrl, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	urlRes, err := s.UrlClient.GetUrl(ctx, &urlv1.GetUrlRequest{UrlId: urlId})
	if err != nil {
		return "", err
	}
	s.putInCache(urlId, urlRes.GetUrl().GetRedirectUrl())
	return urlRes.GetUrl().GetRedirectUrl(), nil
}
