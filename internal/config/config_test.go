package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {

	tests := []struct {
		name          string
		args          []string
		serverAddress string
		baseURL       string
		filePath      string
		want          *Config
	}{
		{
			name: "default values",
			args: nil,
			want: &Config{
				Server: &ServerConfig{
					Port: ":8080",
				},
				ShorterService: &ShorterServiceConfig{
					ShortURLAddr:    "http://localhost:8080",
					FileStoragePath: "storage.json",
				},
			},
		},
		{
			name: "both flags",
			args: []string{"-a=1111", "-b=https://google.com"},
			want: &Config{
				Server: &ServerConfig{
					Port: "1111",
				},
				ShorterService: &ShorterServiceConfig{
					ShortURLAddr:    "https://google.com",
					FileStoragePath: "storage.json",
				},
			},
		},
		{
			name: "value as separate argument",
			args: []string{"-a", ":9090", "-b", "http://example.com"},
			want: &Config{
				Server: &ServerConfig{
					Port: ":9090",
				},
				ShorterService: &ShorterServiceConfig{
					ShortURLAddr:    "http://example.com",
					FileStoragePath: "storage.json",
				},
			},
		},
		{
			name: "only port",
			args: []string{"-a=:3000"},
			want: &Config{
				Server: &ServerConfig{
					Port: ":3000",
				},
				ShorterService: &ShorterServiceConfig{
					ShortURLAddr:    "http://localhost:8080",
					FileStoragePath: "storage.json",
				},
			},
		},
		{
			name: "only short url address",
			args: []string{"-b=http://short.ly"},
			want: &Config{
				Server: &ServerConfig{
					Port: ":8080",
				},
				ShorterService: &ShorterServiceConfig{
					ShortURLAddr:    "http://short.ly",
					FileStoragePath: "storage.json",
				},
			},
		},
		{
			name:          "only SERVER_ADDRESS environment variable",
			serverAddress: ":9090",
			want: &Config{
				Server:         &ServerConfig{Port: ":9090"},
				ShorterService: &ShorterServiceConfig{ShortURLAddr: "http://localhost:8080", FileStoragePath: "storage.json"},
			},
		},
		{
			name:    "only BASE_URL environment variable",
			baseURL: "https://short.example.com",
			want: &Config{
				Server:         &ServerConfig{Port: ":8080"},
				ShorterService: &ShorterServiceConfig{ShortURLAddr: "https://short.example.com", FileStoragePath: "storage.json"},
			},
		},
		{
			name:          "both environment variables",
			serverAddress: ":9090",
			baseURL:       "https://short.example.com",
			want: &Config{
				Server:         &ServerConfig{Port: ":9090"},
				ShorterService: &ShorterServiceConfig{ShortURLAddr: "https://short.example.com", FileStoragePath: "storage.json"},
			},
		},
		{
			name:          "environment variables override flags",
			args:          []string{"-a=:3000", "-b=http://short.ly"},
			serverAddress: ":9090",
			baseURL:       "https://short.example.com",
			want: &Config{
				Server:         &ServerConfig{Port: ":9090"},
				ShorterService: &ShorterServiceConfig{ShortURLAddr: "https://short.example.com", FileStoragePath: "storage.json"},
			},
		},
		{
			name: "file storage flag",
			args: []string{"-f=/tmp/shortener.json"},
			want: &Config{Server: &ServerConfig{Port: ":8080"}, ShorterService: &ShorterServiceConfig{ShortURLAddr: "http://localhost:8080", FileStoragePath: "/tmp/shortener.json"}},
		},
		{
			name:     "file storage environment variable overrides flag",
			args:     []string{"-f=/tmp/from-flag.json"},
			filePath: "/tmp/from-env.json",
			want:     &Config{Server: &ServerConfig{Port: ":8080"}, ShorterService: &ShorterServiceConfig{ShortURLAddr: "http://localhost:8080", FileStoragePath: "/tmp/from-env.json"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SERVER_ADDRESS", tt.serverAddress)
			t.Setenv("BASE_URL", tt.baseURL)
			t.Setenv("FILE_STORAGE_PATH", tt.filePath)
			setArgs(t, tt.args)
			assert.Equal(t, tt.want, NewConfig())
		})
	}
}

// setArgs помогает обновить аргументы при множественных тестах
func setArgs(t *testing.T, args []string) {
	oldArgs, oldFlags := os.Args, flag.CommandLine
	os.Args = append([]string{"app_flags"}, args...)
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	t.Cleanup(func() {
		os.Args, flag.CommandLine = oldArgs, oldFlags
	})
}
