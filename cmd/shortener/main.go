package main

import (
	"log"
	"os"

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

	file, err := os.OpenFile(cfg.ShorterService.FileStoragePath, os.O_RDWR|os.O_CREATE, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()

	db := repository.NewFileRepository(file)
	if err = db.Load(); err != nil {
		return err
	}

	shortService := service.NewShorterService(db, cfg.ShorterService.ShortURLAddr)
	shorterAPI := handler.NewShorterAPI(shortService)

	srv := server.NewServer(cfg.Server.Port, shorterAPI)

	return srv.Start()
}
