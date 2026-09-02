package main

import (
	"log"

	"github.com/presswintox/practicum-shortener/internal/config"
	"github.com/presswintox/practicum-shortener/internal/handler"
	"github.com/presswintox/practicum-shortener/internal/repository"
	"github.com/presswintox/practicum-shortener/internal/server"
	"github.com/presswintox/practicum-shortener/internal/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run init all dependencies and run server
func run() error {
	cfg := config.NewConfig()
	db := repository.NewMemoryRepository()
	shortService := service.NewShorterService(db, cfg.ShorterService.ShortUrlAddr)
	shorterApi := handler.NewShorterApi(shortService)

	srv := server.NewServer(cfg.Server.Port, shorterApi)

	return srv.Start()
}
