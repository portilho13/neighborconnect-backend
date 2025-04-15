package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/marketplace"
)

func CreateItem(item models.Item, dbPool *pgxpool.Pool) error {
	query := `INSERT INTO marketplace.item (name, category, description) 
	VALUES ($1, $2, $3)`

	_, err := dbPool.Exec(context.Background(), query,
		item.Name,
		item.Category,
		item.Description,
	)

	if err != nil {
		return nil
	}
	return nil
}

func GetItemById(id int, dbPool *pgxpool.Pool) (models.Item, error) {
	query := `SELECT (id, name, category, description) FROM marketplace.item WHERE id = $1`

	var item models.Item

	err := dbPool.QueryRow(context.Background(), query, id).Scan(
		&item.Id,
		&item.Name,
		&item.Category,
		&item.Description,
	)

	if err != nil {
		return models.Item{}, nil
	}

	return item, nil

}
