package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
)

func CreateManagerTransaction(manager_transaction models.Manager_Transaction, dbPool *pgxpool.Pool) error {
	query := `INSERT INTO users.manager_transaction (type, amount, date, description, users_id, manager_id)
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := dbPool.Exec(context.Background(), query,
		manager_transaction.Type,
		manager_transaction.Amount,
		manager_transaction.Date,
		manager_transaction.Description,
		manager_transaction.Users_Id,
		manager_transaction.Manager_Id,
	)

	if err != nil {
		return err
	}

	return nil
}

func GetManagerTransactionsByManagerId(manager_id int, dbPool *pgxpool.Pool) ([]models.Manager_Transaction, error) {
	var manager_transactions []models.Manager_Transaction

	query := `SELECT id, type, amount, date, description, users_id, manager_id FROM users.manager_transaction
	WHERE manager_id = $1`

	rows, err := dbPool.Query(context.Background(), query, manager_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var manager_transaction models.Manager_Transaction

		err := rows.Scan(
			&manager_transaction.Id,
			&manager_transaction.Type,
			&manager_transaction.Amount,
			&manager_transaction.Date,
			&manager_transaction.Description,
			&manager_transaction.Users_Id,
			&manager_transaction.Manager_Id,
		)

		if err != nil {
			return nil, err
		}
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return manager_transactions, nil
}
