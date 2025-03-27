package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func connect() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), "postgres://username:password@localhost:5432/mydb")
	if err != nil {
		return nil, err
	}
	return conn, nil
}
