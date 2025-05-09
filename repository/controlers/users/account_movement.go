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

func GetAccountMovementsByAccountId(account_id int, dbPool *pgxpool.Pool) ([]models.Account_Movement, error) {
	var account_movements []models.Account_Movement

	query := `SELECT id, amount, created_at, account_id, type FROM users.account_movement
	WHERE account_id = $1`

	rows, err := dbPool.Query(context.Background(), query, account_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var account_movement models.Account_Movement

		err := rows.Scan(
			&account_movement.Id,
			&account_movement.Ammount,
			&account_movement.Created_at,
			&account_movement.Account_id,
			&account_movement.Type,
		)

		if err != nil {
			return nil, err
		}

		account_movements = append(account_movements, account_movement)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return account_movements, nil
}
