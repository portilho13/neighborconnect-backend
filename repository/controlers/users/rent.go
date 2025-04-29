package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
)

func getLastDayOfMonth(year int, month time.Month) int {
	// Go to the first day of the *next* month
	firstOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	// Subtract one day to get the last day of the current month
	lastDay := firstOfNextMonth.AddDate(0, 0, -1)
	return lastDay.Day()
}

func CreateRent(dbPool *pgxpool.Pool) error {
	apartments, err := GetAllApartments(dbPool)
	if err != nil {
		return err
	}

	now := time.Now()

	month := now.Month()

	year := now.Year()

	query := `
	INSERT INTO users.rent (month, year, base_amount, reduction, final_amount, apartment_id, status, due_day)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
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
			getLastDayOfMonth(year, month),
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func GetRentByApartmentId(apartment_id int, dbPool *pgxpool.Pool) ([]models.Rent, error) {
	var rents []models.Rent

	query := `SELECT id, month, year, base_amount, reduction, final_amount, apartment_id, status, due_day FROM users.rent WHERE apartment_id = $1 ORDER BY year DESC, month DESC`

	rows, err := dbPool.Query(context.Background(), query, apartment_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var rent models.Rent
		err := rows.Scan(
			&rent.Id,
			&rent.Month,
			&rent.Year,
			&rent.Base_Amount,
			&rent.Reduction,
			&rent.Final_Amount,
			&rent.Apartment_Id,
			&rent.Status,
			&rent.Due_day,
		)

		if err != nil {
			return nil, err
		}

		rents = append(rents, rent)

	}

	return rents, nil
}
