package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/marketplace"
)

func CreateTransaction(transaction models.Transaction, dbPool *pgxpool.Pool) (int, error) {
	query := `INSERT INTO marketplace.transaction 
		(final_price, transaction_time, transaction_type, seller_id, buyer_id, listing_id, payment_status, payment_due_time) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id`

	var insertedId int
	err := dbPool.QueryRow(context.Background(), query,
		transaction.Final_Price,
		transaction.Transaction_Time,
		transaction.Transaction_Type,
		transaction.Seller_Id,
		transaction.Buyer_Id,
		transaction.Listing_Id,
		transaction.Payment_Status,
		transaction.Payment_Due_time,
	).Scan(&insertedId)

	if err != nil {
		return 0, err
	}

	return insertedId, nil
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

func GetTransactionsByBuyerId(buyer_id int, dbPool *pgxpool.Pool) ([]models.Transaction, error) {
	var transactions []models.Transaction

	query := `SELECT id, final_price, transaction_time, transaction_type, buyer_id, seller_id, listing_id, payment_status, payment_due_time FROM marketplace.transaction WHERE buyer_id = $1 AND payment_status = 'unpaid'`

	rows, err := dbPool.Query(context.Background(), query, buyer_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
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
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return transactions, nil
}

func GetAllTransactions(dbPool *pgxpool.Pool) ([]models.Transaction, error) {
	var transactions []models.Transaction

	query := `SELECT id, final_price, transaction_time, transaction_type, buyer_id, seller_id, listing_id, payment_status, payment_due_time FROM marketplace.transaction`

	rows, err := dbPool.Query(context.Background(), query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
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
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return transactions, nil
}

func DeleteTransactionById(transaction_id int, dbPool *pgxpool.Pool) error {
	query := `DELETE FROM marketplace.transaction WHERE id = $1`

	_, err := dbPool.Exec(context.Background(), query, transaction_id)
	if err != nil {
		return err
	}

	return nil
}
