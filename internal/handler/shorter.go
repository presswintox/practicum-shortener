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

type Server struct {
	service ShorterService
}

func NewServer(s ShorterService) *Server {
	return &Server{service: s}
}

func (s *Server) DoShortUrlHandler(c *echo.Context) error {
	r := c.Request()
	if r.Method != http.MethodPost {
		return c.String(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return c.String(http.StatusBadRequest, "Не удалось прочитать тело запроса")
	}
	defer r.Body.Close()

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

func (s *Server) GetUrlHandler(c *echo.Context) error {
	r := c.Request()
	if r.Method != http.MethodGet {
		return c.String(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}

	id := c.Param("id")
	url, err := s.service.GetUrl(id)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.Redirect(http.StatusTemporaryRedirect, url)
}
