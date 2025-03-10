package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
)

func CreateManager(manager models.Manager, dbPool *pgxpool.Pool) error {
	query := `INSERT INTO users.manager (name, email, password, phone) VALUES ($1, $2, $3, $4)`
	
	_, err := dbPool.Exec(context.Background(), query,
	manager.Name,
	manager.Email,
	manager.Password,
	manager.Phone,
	)

	if err != nil {
		return err
	}

	return nil
}