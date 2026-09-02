package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {

	tests := []struct {
		name string
		args []string
		want *Config
	}{
		{
			name: "default values",
			args: nil,
			want: &Config{
				Server: &ServerConfig{
					Port: ":8080",
				},
				ShorterService: &ShorterServiceConfig{
					ShortUrlAddr: "http://localhost:8080",
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
					ShortUrlAddr: "https://google.com",
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
					ShortUrlAddr: "http://example.com",
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
					ShortUrlAddr: "http://localhost:8080",
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
					ShortUrlAddr: "http://short.ly",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
