package utils

import (
	"github.com/jackc/pgx/v5/pgxpool"
	repositoryControllersUsers "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
)

func GetManagerIdByUserId(user_id int, dbPool *pgxpool.Pool) (*int, error) {
	user, err := repositoryControllersUsers.GetUsersById(user_id, dbPool)
	if err != nil {
		return nil, err
	}

	apartment, err := repositoryControllersUsers.GetApartmentById(*user.Apartment_id, dbPool)
	if err != nil {
		return nil, err
	}

	return &apartment.Manager_id, nil
}
