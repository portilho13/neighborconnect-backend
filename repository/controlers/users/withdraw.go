package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
)

func CreateWithdraw(withdaw models.Withdraw, dbPool *pgxpool.Pool) error {
	query := `INSERT INTO users.withdraw (amount, created_at, users_id) VALUES ($1, $2, $3)`
	_, err := dbPool.Exec(context.Background(), query,
	withdaw.Ammount,
	withdaw.Created_at,
	withdaw.User_Id, 
	)

	if err != nil {
		return err
	}

	return nil
}