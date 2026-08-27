package main

import (
	"net/http"

	"github.com/presswintox/practicum-shortener/internal/handler"
	"github.com/presswintox/practicum-shortener/internal/service"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {

	shortService := service.NewShorterService()
	server := handler.NewServer(shortService)

	mux := http.NewServeMux()
	mux.HandleFunc("/", server.DoShortUrl)
	mux.HandleFunc("/{id}", server.GetUrl)
	return http.ListenAndServe(`:8080`, mux)
}
