package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateRent(dbPool *pgxpool.Pool) error {
	apartments, err := GetAllApartments(dbPool)
	if err != nil {
		return err
	}

	now := time.Now()

	month := now.Month()

	year := now.Year()

	query := `
	INSERT INTO users.rent (month, year, base_amount, reduction, final_amount, apartment_id, status)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	for _, apartment := range apartments {

		_, err := dbPool.Exec(context.Background(), query, 
			month,
			year,
			apartment.Rent,
			0,
			apartment.Rent,
			apartment.Id,
			"unpaid",
		)

		if err != nil {
			return err
		}
	}

	return nil
}