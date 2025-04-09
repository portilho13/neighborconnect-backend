package repository

import (
	"context"
	"errors"

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

// Função para obter a password do utilizador
func GetUserPassword(userID int, dbPool *pgxpool.Pool) (string, error) {
	var hashedPassword string
	query := `SELECT password FROM users.users WHERE id = $1`

	err := dbPool.QueryRow(context.Background(), query, userID).Scan(&hashedPassword)
	if err != nil {
		return "", err
	}
	return hashedPassword, nil
}

// Função para atualizar a password do utilizador
func UpdateUserPassword(userID int, newPassword string, dbPool *pgxpool.Pool) error {
	query := `UPDATE users.users SET password = $1 WHERE id = $2`
	cmd, err := dbPool.Exec(context.Background(), query, newPassword, userID)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

// Adiciona esta função no teu ficheiro repository
func UpdateUserPasswordByEmail(email string, newPassword string, dbPool *pgxpool.Pool) error {
	query := `UPDATE users.users SET password = $1 WHERE email = $2`
	cmd, err := dbPool.Exec(context.Background(), query, newPassword, email)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("no user found with the given email")
	}

	return nil
}
