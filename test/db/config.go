package test

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetTestDBConnection() (*pgxpool.Pool, error) {
	dbURL := `postgres://neighborconnect:neighborconnect@158.220.93.168:5432/neighborconnect?sslmode=require`

	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}
	return dbPool, err
}