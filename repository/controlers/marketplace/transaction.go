package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/marketplace"
)

func CreateTransaction(transaction models.Transaction, dbPool *pgxpool.Pool) error {
	query := `INSERT INTO marketplace.transaction (final_price, transaction_time, transaction_type, seller_id, buyer_id, listing_id, payment_status, payment_due_time) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := dbPool.Exec(context.Background(), query,
		transaction.Final_Price,
		transaction.Transaction_Time,
		transaction.Transaction_Type,
		transaction.Seller_Id,
		transaction.Buyer_Id,
		transaction.Listing_Id,
		transaction.Payment_Status,
		transaction.Payment_Due_time,
	)

	if err != nil {
		return err
	}

	return nil

}

func GetTransactionById(id int, dbPool *pgxpool.Pool) (models.Transaction, error) {
	query := `SELECT id, final_price, transaction_time, transaction_type, buyer_id, seller_id, listing_id, payment_status, payment_due_time FROM marketplace.transaction WHERE id = $1`

	var transaction models.Transaction

	err := dbPool.QueryRow(context.Background(), query, id).Scan(
		&transaction.Id,
		&transaction.Final_Price,
		&transaction.Transaction_Time,
		&transaction.Transaction_Type,
		&transaction.Buyer_Id,
		&transaction.Seller_Id,
		&transaction.Listing_Id,
		&transaction.Payment_Status,
		&transaction.Payment_Due_time,
	)

	if err != nil {
		return models.Transaction{}, err
	}

	return transaction, nil
}

func UpdateTransactionStatus(status string, id int, dbPool *pgxpool.Pool) error {
	query := `UPDATE marketplace.transaction SET payment_status = $1 WHERE id = $2`

	_, err := dbPool.Exec(context.Background(), query, status, id)
	if err != nil {
		return err
	}

	return nil
}
