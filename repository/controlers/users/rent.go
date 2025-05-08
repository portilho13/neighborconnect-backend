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

func CreateRentForAllApartments(dbPool *pgxpool.Pool) error {
	apartments, err := GetAllOccupiedApartments(dbPool) // Only occupied apartments should have rent
	if err != nil {
		return err
	}

	now := time.Now().UTC()

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

func CreateRentForApartmentById(apartment_id int, dbPool *pgxpool.Pool) error {
	apartment, err := GetApartmentById(apartment_id, dbPool)
	if err != nil {
		return err
	}

	now := time.Now().UTC()

	month := now.Month()

	year := now.Year()

	query := `
	INSERT INTO users.rent (month, year, base_amount, reduction, final_amount, apartment_id, status, due_day)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err = dbPool.Exec(context.Background(), query,
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

	return nil
}

func GetRentByApartmentId(apartment_id int, dbPool *pgxpool.Pool) ([]models.Rent, error) {
	var rents []models.Rent

	query := `SELECT id, month, year, base_amount, reduction, final_amount, apartment_id, status, due_day, pay_day FROM users.rent WHERE apartment_id = $1 ORDER BY year DESC, month DESC`

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
			&rent.Pay_Day,
		)

		if err != nil {
			return nil, err
		}

		rents = append(rents, rent)

	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return rents, nil
}

func GetRentById(rent_id int, dbPool *pgxpool.Pool) (models.Rent, error) {
	var rent models.Rent
	query := `SELECT id, month, year, base_amount, reduction, final_amount, apartment_id, status, due_day, pay_day FROM users.rent WHERE id = $1`

	err := dbPool.QueryRow(context.Background(), query, rent_id).Scan(
		&rent.Id,
		&rent.Month,
		&rent.Year,
		&rent.Base_Amount,
		&rent.Reduction,
		&rent.Final_Amount,
		&rent.Apartment_Id,
		&rent.Status,
		&rent.Due_day,
		&rent.Pay_Day,
	)

	if err != nil {
		return models.Rent{}, nil
	}

	return rent, nil

}

func UpdateRentReductionAndFinalAmount(rent_id int, new_reduction float64, new_amount float64, dbPool *pgxpool.Pool) error {
	query := `UPDATE users.rent SET reduction = $1, final_amount = $2 WHERE id = $3`

	_, err := dbPool.Exec(context.Background(), query, new_reduction, new_amount, rent_id)
	if err != nil {
		return err
	}

	return nil

}

func UpdateRentStatus(status string, rent_id int, dbPool *pgxpool.Pool) error {
	query := `UPDATE users.rent SET status = $1 WHERE id = $2`

	_, err := dbPool.Exec(context.Background(), query, status, rent_id)
	if err != nil {
		return err
	}

	return nil

}

func UpdateRentPayday(pay_day time.Time, rent_id int, dbPool *pgxpool.Pool) error {
	query := `UPDATE users.rent SET pay_day = $1 WHERE id = $2`

	_, err := dbPool.Exec(context.Background(), query, pay_day, rent_id)
	if err != nil {
		return err
	}

	return nil

}
