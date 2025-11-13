package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

const (
	DefaultShutdownTimeout = 10 * time.Second

	DefaultHTTPHost = "0.0.0.0"
	DefaultHTTPPort = "8080"
)

type Config struct {
	HTTPServer HTTPServerConfig
	Database   DatabaseConfig
}

type DatabaseConfig struct {
	DSN string
}

type HTTPServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func LoadConfig() (*Config, error) {

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, fallback to system env")
	}
	viper.AutomaticEnv()

	return &Config{
		HTTPServer: HTTPServerConfig{
			Host: valueOrDefault(viper.GetString("APP_HTTP_HOST"), DefaultHTTPHost),
			Port: valueOrDefault(viper.GetString("APP_HTTP_PORT"), DefaultHTTPPort),
		},
		Database: DatabaseConfig{
			DSN: viper.GetString("APP_DATABASE_DSN"),
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
