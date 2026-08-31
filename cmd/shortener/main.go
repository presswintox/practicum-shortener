package main

import (
	"github.com/labstack/echo/v5"
	"github.com/presswintox/practicum-shortener/internal/handler"
	"github.com/presswintox/practicum-shortener/internal/repository"
	"github.com/presswintox/practicum-shortener/internal/service"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {

	db := repository.NewMemoryRepository()
	shortService := service.NewShorterService(db)
	server := handler.NewServer(shortService)

	e := echo.New()
	e.GET("/:id", server.GetUrlHandler)
	e.POST("/", server.DoShortUrlHandler)

	return e.Start(":8080")
}
