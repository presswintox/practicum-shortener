package main

import (
	"github.com/labstack/echo/v5"
	"github.com/presswintox/practicum-shortener/internal/config"
	"github.com/presswintox/practicum-shortener/internal/handler"
	"github.com/presswintox/practicum-shortener/internal/repository"
	"github.com/presswintox/practicum-shortener/internal/service"
)

func main() {
	cfg := config.NewConfig()

	if err := run(cfg); err != nil {
		panic(err)
	}
}

func run(cfg *config.Config) error {

	db := repository.NewMemoryRepository()
	shortService := service.NewShorterService(db, cfg)
	server := handler.NewServer(shortService)

	e := echo.New()
	e.GET("/:id", server.GetUrlHandler)
	e.POST("/", server.DoShortUrlHandler)

	return e.Start(cfg.Port)
}
