package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(databaseURL string) (*pgxpool.Pool, error) {
	var err error
	DBPool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}

	log.Println("Database pool created sucessfully")
	return DBPool, nil
}

func CloseDB(DBPool *pgxpool.Pool) {
	if DBPool != nil {
		DBPool.Close()
		log.Println("Database pool closed")
	}
}