package handler

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v5"
)

type ShorterService interface {
	DoShortURL(url string) (string, string, error)
	GetURL(shortURL string) (string, error)
}
type ShorterAPI struct {
	service ShorterService
}

func NewShorterAPI(service ShorterService) *ShorterAPI {
	return &ShorterAPI{service: service}
}

func (s *ShorterAPI) DoShortURLHandler(c *echo.Context) error {
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

	_, shortURL, err := s.service.DoShortURL(url)
	if err != nil {
		c.Logger().Error(err.Error())
		return echo.ErrInternalServerError
	}

	return c.String(http.StatusCreated, shortURL)
}

func (s *ShorterAPI) GetURLHandler(c *echo.Context) error {
	id := c.Param("id")
	url, err := s.service.GetURL(id)
	if err != nil {
		c.Logger().Info(err.Error())
		return echo.ErrBadRequest
	}

	return c.Redirect(http.StatusTemporaryRedirect, url)
}

type ShortenBody struct {
	URL string `json:"url,omitempty"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

func (s *ShorterAPI) ShortenHandler(c *echo.Context) error {
	var body ShortenBody
	if err := c.Bind(&body); err != nil {
		c.Logger().Info(err.Error())
		return echo.ErrBadRequest
	}

	if body.URL == "" {
		c.Logger().Info("empty url")
		return echo.ErrBadRequest
	}

	_, shortURL, err := s.service.DoShortURL(body.URL)
	if err != nil {
		c.Logger().Error(err.Error())
		return echo.ErrInternalServerError
	}
	return c.JSON(http.StatusCreated, ShortenResponse{Result: shortURL})
}
