package config

import (
	"flag"
	"os"
)

const (
	StoragePostgres string = "postgres"
	StorageFile            = "file"
	StorageMemory          = "memory"
)

type Config struct {
	Server         *ServerConfig
	ShorterService *ShorterServiceConfig
	Database       *DatabaseConfig
}
type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Dsn string
}

type ShorterServiceConfig struct {
	ShortURLAddr    string
	FileStoragePath string
}

func NewConfig() *Config {
	cfg := &Config{
		Server:         &ServerConfig{},
		ShorterService: &ShorterServiceConfig{},
		Database:       &DatabaseConfig{},
	}
	flag.StringVar(&cfg.Server.Port, "a", ":8080", "address and port to run server")
	flag.StringVar(&cfg.ShorterService.ShortURLAddr, "b", "http://localhost:8080", "address and port for short url")
	flag.StringVar(&cfg.ShorterService.FileStoragePath, "f", "", "path to file storage path")
	flag.StringVar(&cfg.Database.Dsn, "d", "", "database connection string")
	flag.Parse()

	if envServerPort := os.Getenv("SERVER_ADDRESS"); envServerPort != "" {
		cfg.Server.Port = envServerPort
	}
	if envShorterURL := os.Getenv("BASE_URL"); envShorterURL != "" {
		cfg.ShorterService.ShortURLAddr = envShorterURL
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		cfg.ShorterService.FileStoragePath = envFileStoragePath
	}
	if envDatabaseConfig := os.Getenv("DATABASE_DSN"); envDatabaseConfig != "" {
		cfg.Database.Dsn = envDatabaseConfig
	}
	return cfg
}

func (c *Config) StorageIs() string {
	if c.Database.Dsn != "" {
		return StoragePostgres
	}
	if c.ShorterService.FileStoragePath != "" {
		return StorageFile
	}
	return StorageMemory
}
