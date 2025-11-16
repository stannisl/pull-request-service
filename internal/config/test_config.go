package config

func TestConfig() *Config {
	return &Config{
		HTTPServer: HTTPServerConfig{
			Host: "localhost",
			Port: "8080",
		},
		Database: DatabaseConfig{
			ConnStr:           "postgres://test:test@localhost:5432/test?sslmode=disable",
			Retries:           5,
			RetrySecondsDelay: 2,
			DriverName:        "postgres",
		},
	}
}
