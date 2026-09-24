package server

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"go.uber.org/zap"
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
	s.setupMiddlewares()
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
func (s *Server) setupMiddlewares() {
	logger, _ := zap.NewProduction()
	s.echo.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:          true,
		LogMethod:       true,
		LogLatency:      true,
		LogStatus:       true,
		LogResponseSize: true,
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {

			logger.Info("request",
				zap.String("Method", v.Method),
				zap.String("URI", v.URI),
				zap.Duration("latency", v.Latency),
				zap.Int("status", v.Status),
				zap.Int("response length", int(v.ResponseSize)),
			)
			return nil
		},
	}))
}
