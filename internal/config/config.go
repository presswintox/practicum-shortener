package config

import "flag"

type Config struct {
	Server         *ServerConfig
	ShorterService *ShorterServiceConfig
}
type ServerConfig struct {
	Port string
}
type ShorterServiceConfig struct {
	ShortUrlAddr string
}

func NewConfig() *Config {
	cfg := &Config{
		Server:         &ServerConfig{},
		ShorterService: &ShorterServiceConfig{},
	}
	flag.StringVar(&cfg.Server.Port, "a", ":8080", "address and port to run server")
	flag.StringVar(&cfg.ShorterService.ShortUrlAddr, "b", "http://localhost:8080", "address and port for short url")
	flag.Parse()

	return cfg
}
