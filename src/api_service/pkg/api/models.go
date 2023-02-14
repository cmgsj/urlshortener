package api

type UrlDto struct {
	UrlId       string `json:"urlId"`
	RedirectUrl string `json:"redirectUrl"`
	NewUrl      string `json:"newUrl"`
}

type CreateUrlRequest struct {
	RedirectUrl string `json:"redirectUrl" binding:"required"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
