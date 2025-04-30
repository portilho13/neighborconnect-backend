package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/events"
)

func CreateCommunityEvent(event models.Community_Event, dbPool *pgxpool.Pool) error {
	query := `INSERT INTO events.community_event (name, percentage, code, capacity, date_time, manager_id, event_image, duration, local, current_ocupation) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := dbPool.Exec(context.Background(), query,
		event.Name,
		event.Percentage,
		event.Code,
		event.Capacity,
		event.Date_Time,
		event.Manager_Id,
		event.Event_Image,
		event.Duration,
		event.Local,
		event.Current_Ocupation,
	)

	if err != nil {
		return err
	}
	return nil
}

func AddUserToCommunityEvent(userId int, eventId int, dbPool *pgxpool.Pool) error {
	query := `INSERT INTO events.many_community_event_has_many_users (community_event_id, users_id) VALUES ($1, $2)`
	_, err := dbPool.Exec(context.Background(), query,
		eventId,
		userId,
	)

	if err != nil {
		return err
	}
	return nil

}

func GetAllEvents(dbPool *pgxpool.Pool) ([]models.Community_Event, error) {
	var events []models.Community_Event

	query := `SELECT id, name, percentage, code, capacity, date_time, manager_id, event_image, duration, local, current_ocupation
	FROM events.community_event`

	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var event models.Community_Event

		err := rows.Scan(
			&event.Id,
			&event.Name,
			&event.Percentage,
			&event.Code,
			&event.Capacity,
			&event.Date_Time,
			&event.Manager_Id,
			&event.Event_Image,
			&event.Duration,
			&event.Local,
			&event.Current_Ocupation,
		)

		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return events, nil
}

func GetAllEventsByManagerId(manager_id int, dbPool *pgxpool.Pool) ([]models.Community_Event, error) {
	var events []models.Community_Event

	query := `SELECT id, name, percentage, code, capacity, date_time, manager_id, event_image, duration, local, current_ocupation
	FROM events.community_event WHERE manager_id = $1`

	rows, err := dbPool.Query(context.Background(), query, manager_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var event models.Community_Event

		err := rows.Scan(
			&event.Id,
			&event.Name,
			&event.Percentage,
			&event.Code,
			&event.Capacity,
			&event.Date_Time,
			&event.Manager_Id,
			&event.Event_Image,
			&event.Duration,
			&event.Local,
			&event.Current_Ocupation,
		)

		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return events, nil
}

func GetEventById(event_id int, dbPool *pgxpool.Pool) (models.Community_Event, error) {
	var event models.Community_Event

	query := `SELECT id, name, percentage, code, capacity, date_time, manager_id, event_image, duration, local, current_ocupation
	FROM events.community_event WHERE id = $1`

	err := dbPool.QueryRow(context.Background(), query, event_id).Scan(
		&event.Id,
		&event.Name,
		&event.Percentage,
		&event.Code,
		&event.Capacity,
		&event.Date_Time,
		&event.Manager_Id,
		&event.Event_Image,
		&event.Duration,
		&event.Local,
		&event.Current_Ocupation,
	)

	if err != nil {
		return models.Community_Event{}, err
	}

	return event, nil
}

func GetEventsByUserId(user_id int, dbPool *pgxpool.Pool) ([]models.Community_Event, error) {
	var users_community_events []models.Community_Event

	query := `SELECT community_event_id
	FROM events.many_community_event_has_many_users WHERE users_id = $1`

	rows, err := dbPool.Query(context.Background(), query, user_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var communityEventId int

		err = rows.Scan(
			&communityEventId,
		)

		if err != nil {
			return nil, err
		}

		event, err := GetEventById(communityEventId, dbPool)
		if err != nil {
			return nil, err
		}

		users_community_events = append(users_community_events, event)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return users_community_events, nil
}
