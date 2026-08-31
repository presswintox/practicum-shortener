package config

import "flag"

type Config struct {
	Port         string
	ShortUrlAddr string
}

func NewConfig() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.Port, "a", ":8080", "address and port to run server")
	flag.StringVar(&cfg.ShortUrlAddr, "b", "http://localhost:8080", "address and port for short url")
	flag.Parse()

	return cfg
}
