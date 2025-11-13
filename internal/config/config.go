package config

import (
	"fmt"
	"os"
	"time"
)

const (
	DefaultHTTPHost = "0.0.0.0"
	DefaultHTTPPort = "8080"
)

type Config struct {
	HTTPServer HTTPServerConfig
}

type HTTPServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func MustLoad() Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}

	return cfg
}

func Load() (Config, error) {
	host := valueOrDefault(os.Getenv("HTTP_HOST"), DefaultHTTPHost)
	port := valueOrDefault(os.Getenv("HTTP_PORT"), DefaultHTTPPort)

	return Config{
		HTTPServer: HTTPServerConfig{
			Host:         host,
			Port:         port,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}, nil
}

func (c HTTPServerConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func valueOrDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}
