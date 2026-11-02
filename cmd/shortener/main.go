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

	var (
		repo         service.ShorterRepository
		infraService server.InfraService
	)

	switch cfg.StorageIs() {
	case config.StoragePostgres:
		postgresRepo, err := repository.NewPostgresRepository(cfg.Database.Dsn)
		if err != nil {
			return err
		}
		defer postgresRepo.Close()
		if err := postgresRepo.Migrate(); err != nil {
			return err
		}

		repo = postgresRepo
		infraService = service.NewInfraService(postgresRepo)
	case config.StorageFile:
		fileRepo, err := repository.NewFileRepository(cfg.ShorterService.FileStoragePath)
		if err != nil {
			return err
		}
		defer fileRepo.Close()

		repo = fileRepo
	default:
		repo = repository.NewMemoryRepository()
	}

	shortService := service.NewShorterService(repo, cfg.ShorterService.ShortURLAddr)
	shorterAPI := handler.NewShorterAPI(shortService)

	srv := server.NewServer(cfg.Server.Port, shorterAPI, infraService)

	return srv.Start()
}
