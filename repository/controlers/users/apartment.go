package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
)

func CreateApartment(apartment models.Apartment, dbPool *pgxpool.Pool) error {
	query := `INSERT INTO users.apartment (n_bedrooms, floor, rent, manager_id, status) VALUES ($1, $2, $3, $4, $5)`

	_, err := dbPool.Exec(context.Background(), query,
		apartment.N_bedrooms,
		apartment.Floor,
		apartment.Rent,
		apartment.Manager_id,
		apartment.Status,
	)

	if err != nil {
		return err
	}

	return nil
}

func UpdateApartmentStatus(apartment_id int, dbPool *pgxpool.Pool) error {
	apartment, err := GetApartmentById(apartment_id, dbPool)
	if err != nil {
		return err
	}

	if apartment.Status == "occupied" {
		return nil
	}

	query := `UPDATE users.apartment SET status = 'occupied' WHERE id = $1`

	_, err = dbPool.Exec(context.Background(), query, apartment_id)
	if err != nil {
		return err
	}

	err = CreateRentForApartmentById(apartment_id, dbPool)
	if err != nil {
		return err
	}

	return nil
}

func GetAllApartments(dbPool *pgxpool.Pool) ([]models.Apartment, error) {
	var apartments []models.Apartment

	query := `SELECT id, n_bedrooms, floor, rent, manager_id, status FROM users.apartment`

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
			&apartment.Status,
		)

		if err != nil {
			return nil, err
		}

		apartments = append(apartments, apartment)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return apartments, nil
}

func GetAllOccupiedApartments(dbPool *pgxpool.Pool) ([]models.Apartment, error) {
	var apartments []models.Apartment

	query := `SELECT id, n_bedrooms, floor, rent, manager_id, status FROM users.apartment WHERE status = 'occupied'`

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
			&apartment.Status,
		)

		if err != nil {
			return nil, err
		}

		apartments = append(apartments, apartment)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return apartments, nil
}

func GetApartmentById(apartment_id int, dbPool *pgxpool.Pool) (models.Apartment, error) {
	var apartment models.Apartment

	query := `SELECT id, n_bedrooms, floor, rent, manager_id, status FROM users.apartment WHERE id = $1`

	err := dbPool.QueryRow(context.Background(), query, apartment_id).Scan(
		&apartment.Id,
		&apartment.N_bedrooms,
		&apartment.Floor,
		&apartment.Rent,
		&apartment.Manager_id,
		&apartment.Status,
	)

	if err != nil {
		return models.Apartment{}, nil
	}

	return apartment, nil
}

func GetAllApartmentsByManagerId(manager_id int, dbPool *pgxpool.Pool) ([]models.Apartment, error) {
	var apartments []models.Apartment

	query := `SELECT id, n_bedrooms, floor, rent, manager_id, status FROM users.apartment WHERE manager_id = $1`

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
			&apartment.Status,
		)

		if err != nil {
			return nil, err
		}

		apartments = append(apartments, apartment)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return apartments, nil
}
