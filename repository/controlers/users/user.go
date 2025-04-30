package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
)

func CreateUser(user models.User, dbPool *pgxpool.Pool) error {
	query := `INSERT INTO users.users (name, email, password, phone, apartment_id, profile_picture) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := dbPool.Exec(context.Background(), query,
		user.Name,
		user.Email,
		user.Password,
		user.Phone,
		user.Apartment_id,
		user.Profile_Picture,
	)
	
	if err != nil {
		return err
	}

	return nil
}

func GetUserByEmail(email string, dbPool *pgxpool.Pool) (models.User, error) {
	var user models.User

	query := `SELECT id, name, email, password, phone, apartment_id, profile_picture FROM users.users WHERE email = $1`

	err := dbPool.QueryRow(context.Background(), query, email).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Phone,
		&user.Apartment_id,
		&user.Profile_Picture,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}


func GetUsersByApartmentId(apartment_id int, dbPool *pgxpool.Pool) ([]models.User, error) {
	var users []models.User

	query := `SELECT id, name, email, password, phone, apartment_id, profile_picture
	FROM users.users WHERE apartment_id = $1`

	rows, err := dbPool.Query(context.Background(), query, apartment_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User

		err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Phone,
			&user.Apartment_id,
			&user.Profile_Picture,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}


	if rows.Err() != nil {
		return nil, rows.Err()
	}


	return users, nil
}
