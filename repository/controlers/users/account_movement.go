package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
)

func CreateAccountMovement(account_movement models.Account_Movement, dbPool *pgxpool.Pool) error {
	query := `INSERT INTO users.account_movement (amount, created_at, account_id, type) VALUES ($1, $2, $3, $4)`
	_, err := dbPool.Exec(context.Background(), query,
	account_movement.Ammount,
	account_movement.Created_at,
	account_movement.Account_id,
	account_movement.Type, 
	)

	if err != nil {
		return err
	}

	return nil
}