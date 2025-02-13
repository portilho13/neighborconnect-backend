package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DBPool *pgxpool.Pool

func InitDB(databaseURL string) {
	var err error
	DBPool, err = pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	log.Println("Database pool created sucessfully")
}

func CloseDB() {
	if DBPool != nil {
		DBPool.Close()
		log.Println("Database pool closed")
	}
}