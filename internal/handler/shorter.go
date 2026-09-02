package handler

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v5"
)

type ShorterService interface {
	DoShortUrl(url string) (string, string, error)
	GetUrl(shortUrl string) (string, error)
}
type ShorterApi struct {
	service ShorterService
}

func NewShorterApi(service ShorterService) *ShorterApi {
	return &ShorterApi{service: service}
}

func (s *ShorterApi) DoShortUrlHandler(c *echo.Context) error {
	r := c.Request()

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		c.Logger().Info(err.Error())
		return echo.ErrBadRequest
	}

	url := string(body)
	if url == "" {
		c.Logger().Info("empty url")
		return echo.ErrBadRequest
	}

	_, shortUrl, err := s.service.DoShortUrl(url)
	if err != nil {
		c.Logger().Error(err.Error())
		return echo.ErrInternalServerError
	}

	return c.String(http.StatusCreated, shortUrl)
}

func (s *ShorterApi) GetUrlHandler(c *echo.Context) error {
	id := c.Param("id")
	url, err := s.service.GetUrl(id)
	if err != nil {
		c.Logger().Info(err.Error())
		return echo.ErrBadRequest
	}

	return c.Redirect(http.StatusTemporaryRedirect, url)
}
