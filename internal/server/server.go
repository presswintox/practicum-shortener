package server

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/presswintox/practicum-shortener/internal/config/db"
	"go.uber.org/zap"
)

type ShorterAPI interface {
	DoShortURLHandler(c *echo.Context) error
	GetURLHandler(c *echo.Context) error
	ShortenHandler(c *echo.Context) error
}

type Server struct {
	echo       *echo.Echo
	port       string
	shorterAPI ShorterAPI
	db         *db.Database
}

func NewServer(port string, shorterAPI ShorterAPI, db *db.Database) *Server {
	e := echo.New()

	s := &Server{
		echo:       e,
		port:       port,
		shorterAPI: shorterAPI,
		db:         db,
	}
	s.setupMiddlewares()
	s.setupRouters()
	return s
}

func (s *Server) Start() error {
	return s.echo.Start(s.port)
}

func (s *Server) setupRouters() {
	s.echo.GET("/ping", s.pingHandler)

	s.echo.GET("/:id", s.shorterAPI.GetURLHandler)
	s.echo.POST("/", s.shorterAPI.DoShortURLHandler)

	api := s.echo.Group("/api")
	api.POST("/shorten", s.shorterAPI.ShortenHandler)
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

	s.echo.Use(middleware.Gzip())
	s.echo.Use(middleware.Decompress())
}

func (s *Server) pingHandler(c *echo.Context) error {
	err := s.db.Ping()
	if err != nil {
		c.Logger().Error(err.Error())
		return err
	}
	return c.String(http.StatusOK, "database is ready")
}
