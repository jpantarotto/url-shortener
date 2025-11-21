package config

import (
	"fmt"
	"os"
	"strconv"

	// this will automatically load your .env file:
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	DB           PostgresConfig
	Port         string
	CounterStart int64
}

type PostgresConfig struct {
	URL string
}

func Load() (*Config, error) {
	counterStr := os.Getenv("COUNTER_START")
	counter, err := strconv.ParseInt(counterStr, 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting counter start to int: %v", counterStr)
	}
	cfg := &Config{
		Port:         os.Getenv("PORT"),
		CounterStart: counter,
		DB: PostgresConfig{
			URL: os.Getenv("POSTGRES_URL"),
		},
	}

	return cfg, nil
}
