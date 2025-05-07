package test

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetTestDBConnection() (*pgxpool.Pool, error) {
	dbURL := `postgres://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable`

	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}
	return dbPool, err
}

// CleanDatabase removes all records from a table
func CleanDatabase(dbPool *pgxpool.Pool, tableName string) {
	_, _ = dbPool.Exec(context.Background(), "TRUNCATE "+tableName+" RESTART IDENTITY CASCADE")
}