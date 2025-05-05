package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
)

func CreateManagerActivity(manager_activity models.Manager_Activity, dbPool *pgxpool.Pool) error {
	query := `INSERT INTO users.manager_activity (type, description, created_at, manager_id)
	VALUES ($1, $2, $3, $4)`

	_, err := dbPool.Exec(context.Background(), query,
		manager_activity.Type,
		manager_activity.Description,
		manager_activity.Created_At,
		manager_activity.Manager_Id,
	)

	if err != nil {
		return err
	}
	return nil
}

func GetManagerActivityByManagerId(manager_id int, dbPool *pgxpool.Pool) ([]models.Manager_Activity, error) {
	var manager_activities []models.Manager_Activity

	query := `SELECT id, type, description, created_at, manager_id
	FROM users.manager_activity WHERE manager_id = $1`

	rows, err := dbPool.Query(context.Background(), query, manager_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var manager_activity models.Manager_Activity

		err := rows.Scan(
			&manager_activity.Id,
			&manager_activity.Type,
			&manager_activity.Description,
			&manager_activity.Created_At,
			&manager_activity.Manager_Id,
		)

		if err != nil {
			return nil, err
		}
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return manager_activities, nil
}