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
		return c.String(http.StatusBadRequest, "Не удалось прочитать тело запроса")
	}

	url := string(body)
	if url == "" {
		return c.String(http.StatusBadRequest, "Пустой URL")
	}

	_, shortUrl, err := s.service.DoShortUrl(url)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Не удалось сохранить URL")
	}

	return c.String(http.StatusCreated, shortUrl)
}

func (s *ShorterApi) GetUrlHandler(c *echo.Context) error {
	id := c.Param("id")
	url, err := s.service.GetUrl(id)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.Redirect(http.StatusTemporaryRedirect, url)
}
