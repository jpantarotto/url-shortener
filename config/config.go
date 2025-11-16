package config

import (
	"os"
	// this will automatically load your .env file:
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	DB   PostgresConfig
	Port string
}

type PostgresConfig struct {
	Username string
	Password string
	URL      string
	Port     string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port: os.Getenv("PORT"),
		DB: PostgresConfig{
			Username: os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PWD"),
			URL:      os.Getenv("POSTGRES_URL"),
			Port:     os.Getenv("POSTGRES_PORT"),
		},
	}

	return cfg, nil
}
