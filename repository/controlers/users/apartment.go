package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
)

func CreateApartment(apartment models.Apartment, dbPool *pgxpool.Pool) error {
	query := `INSERT INTO users.apartment (n_bedrooms, floor, manager_id) VALUES ($1, $2, $3)`

	_, err := dbPool.Exec(context.Background(), query,
		apartment.N_bedrooms,
		apartment.Floor,
		apartment.Manager_id)

	if err != nil {
		return err // Return the actual error
	}
	return nil
}

func GetAllApartments(dbPool *pgxpool.Pool) ([]models.Apartment, error) {
	var apartments []models.Apartment

	query := `SELECT id, n_bedrooms, floor, rent, manager_id FROM users.apartment`

	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var apartment models.Apartment

		err := rows.Scan(
			&apartment.Id,
			&apartment.N_bedrooms,
			&apartment.Floor,
			&apartment.Rent,
			&apartment.Manager_id,
		)

		if err != nil {
			return nil, err
		}

		apartments = append(apartments, apartment)
	}

	return apartments, nil
}

func GetAllApartmentsByManagerId(manager_id int, dbPool *pgxpool.Pool) ([]models.Apartment, error) {
	var apartments []models.Apartment

	query := `SELECT id, n_bedrooms, floor, rent, manager_id FROM users.apartment WHERE manager_id = $1`

	rows, err := dbPool.Query(context.Background(), query, manager_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var apartment models.Apartment

		err := rows.Scan(
			&apartment.Id,
			&apartment.N_bedrooms,
			&apartment.Floor,
			&apartment.Rent,
			&apartment.Manager_id,
		)

		if err != nil {
			return nil, err
		}

		apartments = append(apartments, apartment)
	}

	return apartments, nil
}