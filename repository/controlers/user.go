package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
)

func CreateUser(user models.User, dbPool *pgxpool.Pool) error {
	query := `INSERT INTO users.users (name, email, password, phone) VALUES ($1, $2, $3, $4)`

	_, err := dbPool.Exec(context.Background(), query,
		user.Name,
		user.Email,
		user.Password,
		user.Phone)
	
	if err != nil {
		return err
	}

	return nil
}