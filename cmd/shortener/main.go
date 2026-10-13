package main

import (
	"log"

	"github.com/presswintox/practicum-shortener/internal/config"
	"github.com/presswintox/practicum-shortener/internal/config/db"
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

	postgres, err := db.NewPostgresDatabase(cfg.Database.Dsn)
	if err != nil {
		return err
	}
	defer postgres.Close()

	repo, err := repository.NewFileRepository(cfg.ShorterService.FileStoragePath)
	if err != nil {
		return err
	}
	defer repo.Close()

	shortService := service.NewShorterService(repo, cfg.ShorterService.ShortURLAddr)
	shorterAPI := handler.NewShorterAPI(shortService)

	srv := server.NewServer(cfg.Server.Port, shorterAPI, postgres)

	return srv.Start()
}
