package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

const (
	// Default Server Fields

	DefaultReadTimeout     = 5 * time.Second
	DefaultWriteTimeout    = 10 * time.Second
	DefaultShutdownTimeout = 15 * time.Second
	DefaultHTTPHost        = "0.0.0.0"
	DefaultHTTPPort        = "8080"
)

type Config struct {
	HTTPServer HTTPServerConfig
	Database   DatabaseConfig
}

type DatabaseConfig struct {
	ConnStr           string
	Retries           int
	RetrySecondsDelay uint
}

type HTTPServerConfig struct {
	Host string
	Port string
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
			ConnStr:           viper.GetString("APP_DATABASE_CONN_URL"),
			Retries:           viper.GetInt("APP_DATABASE_RETRIES"),
			RetrySecondsDelay: viper.GetUint("APP_DATABASE_RETRIES_SECONDS_DELAY"),
		},
	}, nil
}

func valueOrDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}
