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

func (c HTTPServerConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func valueOrDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}

	return value
}
