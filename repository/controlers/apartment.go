package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
)

func CreateApartment(apartment models.Apartment, dbPool *pgxpool.Pool) error {
	query := `INSERT INTO users.appartment (n_bedrooms, floor, rent, manager_id) VALUES ($1, $2, $3, $4)`

	_, err := dbPool.Exec(context.Background(), query, 
	apartment.N_bedrooms, 
	apartment.Floor, 
	apartment.Rent, 
	apartment.Manager_id)

	if err != nil {
		return nil
	}
	return nil
}