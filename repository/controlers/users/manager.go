package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
)

func CreateManager(manager models.Manager, dbPool *pgxpool.Pool) error {
	query := `INSERT INTO users.manager (name, email, password, phone) VALUES ($1, $2, $3, $4)`

	_, err := dbPool.Exec(context.Background(), query,
		manager.Name,
		manager.Email,
		manager.Password,
		manager.Phone,
	)

	if err != nil {
		return err
	}

	return nil
}

func GetManagerByEmail(email string, dbPool *pgxpool.Pool) (models.Manager, error) {
	var manager models.Manager
	query := `SELECT id, name, email, password, phone FROM users.manager WHERE email = $1`

	err := dbPool.QueryRow(context.Background(), query, email).Scan(
		&manager.Id,
		&manager.Name,
		&manager.Email,
		&manager.Password,
		&manager.Phone,
	)

	if err != nil {
		return models.Manager{}, err
	}

	return manager, nil
}

func GetManagerIdByApartmentId(apartment_id int, dbPool *pgxpool.Pool) (*int, error) {
	query := `SELECT manager_id FROM users.apartment WHERE id = $1`

	var manager_id int

	err := dbPool.QueryRow(context.Background(), query, apartment_id).Scan(
		&manager_id,
	)

	if err != nil {
		return nil, err
	}

	return &manager_id, nil
}
