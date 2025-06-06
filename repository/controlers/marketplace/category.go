package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/marketplace"
)

func CreateCategory(category models.Category, dbPool *pgxpool.Pool) error {
	query := `INSERT INTO marketplace.category (name) 
	VALUES ($1)`

	_, err := dbPool.Exec(context.Background(), query,
		category.Name,
	)

	if err != nil {
		return err
	}

	return nil
}

func GetAllCategories(dbPool *pgxpool.Pool) ([]models.Category, error) {
	var categories []models.Category

	query := `SELECT id, name FROM marketplace.category`

	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var category models.Category

		err := rows.Scan(
			&category.Id,
			&category.Name,
		)

		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return categories, nil
}

func GetCategoryById(category_id int, dbPool *pgxpool.Pool) (models.Category, error) {
	var category models.Category

	query := `SELECT id, name FROM marketplace.category WHERE id = $1`

	err := dbPool.QueryRow(context.Background(), query, category_id).Scan(
		&category.Id,
		&category.Name,
	)

	if err != nil {
		return models.Category{}, nil
	}

	return category, nil
}
