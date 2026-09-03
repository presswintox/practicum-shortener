package server

import (
	"github.com/labstack/echo/v5"
)

type ShorterAPI interface {
	DoShortURLHandler(c *echo.Context) error
	GetURLHandler(c *echo.Context) error
}

type Server struct {
	echo       *echo.Echo
	port       string
	shorterAPI ShorterAPI
}

func NewServer(port string, shorterAPI ShorterAPI) *Server {
	e := echo.New()

	s := &Server{
		echo:       e,
		port:       port,
		shorterAPI: shorterAPI,
	}
	s.setupRouters()
	return s
}

func (s *Server) Start() error {
	return s.echo.Start(s.port)
}

func (s *Server) setupRouters() {
	s.echo.GET("/:id", s.shorterAPI.GetURLHandler)
	s.echo.POST("/", s.shorterAPI.DoShortURLHandler)
}
