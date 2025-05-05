package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
)

func CreateAccount(account models.Account, dbPool *pgxpool.Pool) error {
	query := `INSERT INTO users.account (account_number, balance, currency, users_id) VALUES ($1, $2, $3, $4)`
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

func GetAccountByUserId(id int, dbPool *pgxpool.Pool) (models.Account, error) {
	var account models.Account

	query := `SELECT id, account_number, balance, currency, users_id FROM users.account WHERE users_id = $1`
	err := dbPool.QueryRow(context.Background(), query, id).Scan(
		&account.Id,
		&account.Account_number,
		&account.Balance,
		&account.Currency,
		&account.Users_id,
	)

	if err != nil {
		return models.Account{}, err
	}

	return account, nil
}

func UpdateAccountBalance(account_id int, balance float64, dbPool *pgxpool.Pool) error {
	query := `UPDATE users.account SET balance = $1 WHERE id = $2`

	_, err := dbPool.Exec(context.Background(), query, balance, account_id)
	if err != nil {
		return err
	}

	return nil
}