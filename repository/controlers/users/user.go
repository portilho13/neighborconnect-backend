package repository

import (
	"context"
	// "golang.org/x/crypto/argon2"
	// "encoding/base64"
	// "strings"

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
