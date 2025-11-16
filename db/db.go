package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jpantarotto/url-shortener/config"
)

func Connect(pgConfig config.PostgresConfig) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), pgConfig.URL)
	return conn, err
}
