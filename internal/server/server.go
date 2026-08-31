package server

import (
	"github.com/labstack/echo/v5"
	"github.com/presswintox/practicum-shortener/internal/handler"
)

type Server struct {
	echo       *echo.Echo
	port       string
	shorterApi *handler.ShorterApi
}

func NewServer(port string, shorterApi *handler.ShorterApi) *Server {
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
