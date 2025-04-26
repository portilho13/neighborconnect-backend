package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/marketplace"
)

func CreateCategory(category models.Category, dbPool *pgxpool.Pool) error {
	query := `INSERT INTO marketplace.category (name, url) 
	VALUES ($1, $2)`

	_, err := dbPool.Exec(context.Background(), query,
		category.Name,
		category.Url,
	)

	if err != nil {
		return err
	}

	return nil
}

func GetAllCategories(dbPool *pgxpool.Pool) ([]models.Category, error) {
	var categories []models.Category

	query := `SELECT id, name, url FROM marketplace.category`

	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var category models.Category

		err := rows.Scan(
			&category.Id,
			&category.Name,
			&category.Url,
		)

		if err != nil {
			return nil, err
		}


		categories = append(categories, category)
	}

	return categories, nil
}