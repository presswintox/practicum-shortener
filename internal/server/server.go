package server

import (
	"github.com/labstack/echo/v5"
)

type ShorterApi interface {
	DoShortUrlHandler(c *echo.Context) error
	GetUrlHandler(c *echo.Context) error
}

type Server struct {
	echo       *echo.Echo
	port       string
	shorterApi ShorterApi
}

func NewServer(port string, shorterApi ShorterApi) *Server {
	e := echo.New()

	s := &Server{
		echo:       e,
		port:       port,
		shorterApi: shorterApi,
	}
	s.setupRouters()
	return s
}

func (s *Server) Start() error {
	return s.echo.Start(s.port)
}

func (s *Server) setupRouters() {
	s.echo.GET("/:id", s.shorterApi.GetUrlHandler)
	s.echo.POST("/", s.shorterApi.DoShortUrlHandler)
}
