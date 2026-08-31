package main

import (
	"net/http"

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

	mux := http.NewServeMux()
	mux.HandleFunc("/", server.DoShortUrlHandler)
	mux.HandleFunc("/{id}", server.GetUrlHandler)
	return http.ListenAndServe(`:8080`, mux)
}
