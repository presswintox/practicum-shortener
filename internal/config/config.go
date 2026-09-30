package config

import (
	"flag"
	"os"
)

type Config struct {
	Server         *ServerConfig
	ShorterService *ShorterServiceConfig
}
type ServerConfig struct {
	Port string
}
type ShorterServiceConfig struct {
	ShortURLAddr string
}

func NewConfig() *Config {
	cfg := &Config{
		Server:         &ServerConfig{},
		ShorterService: &ShorterServiceConfig{},
	}
	flag.StringVar(&cfg.Server.Port, "a", ":8080", "address and port to run server")
	flag.StringVar(&cfg.ShorterService.ShortURLAddr, "b", "http://localhost:8080", "address and port for short url")
	flag.Parse()

	if envServerPort := os.Getenv("SERVER_ADDRESS"); envServerPort != "" {
		cfg.Server.Port = envServerPort
	}
	if envShorterURL := os.Getenv("BASE_URL"); envShorterURL != "" {
		cfg.ShorterService.ShortURLAddr = envShorterURL
	}
	return cfg
}
