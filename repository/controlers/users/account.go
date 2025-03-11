package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
)

func CreateAccount(account models.Account, dbPool *pgxpool.Pool) error {
	query := `INSERT INTO users.account (account_number, balance, currency, users_id) VALUES ($1, $2, $3)`
	_, err := dbPool.Exec(context.Background(), query,
	account.Account_number,
	account.Balance,
	account.Currency,
	account.Users_id, 
	)

	if err != nil {
		return err
	}

	return nil
}