package api

type UrlDTO struct {
	UrlId       string `json:"urlId"`
	RedirectUrl string `json:"redirectUrl"`
}

type CreateUrlRequest struct {
	RedirectUrl string `json:"redirectUrl" binding:"required"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
