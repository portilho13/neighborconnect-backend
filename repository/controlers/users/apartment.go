package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
)

func CreateApartment(apartment models.Apartment, db *pgxpool.Pool) error {
	query := `INSERT INTO users.apartment (n_bedrooms, floor, manager_id) VALUES ($1, $2, $3)`

	_, err := db.Exec(context.Background(), query,
		apartment.N_bedrooms,
		apartment.Floor,
		apartment.Manager_id)

	if err != nil {
		return err // Return the actual error
	}
	return nil
}