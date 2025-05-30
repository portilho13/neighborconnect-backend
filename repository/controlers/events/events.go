package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/events"
)

func CreateCommunityEvent(event models.Community_Event, dbPool *pgxpool.Pool) error {
	query := `INSERT INTO events.community_event (name, percentage, code, capacity, date_time, manager_id, event_image, duration, local, current_ocupation, status, expiration_date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
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
		event.Status,
		event.Expiration_Date,
	)

	if err != nil {
		return err
	}
	return nil
}

func AddUserToCommunityEvent(user_event models.User_Event, dbPool *pgxpool.Pool) error {
	query := `INSERT INTO events.many_community_event_has_many_users (community_event_id, users_id, isrewarded, claimedreward) VALUES ($1, $2, $3, $4)`
	_, err := dbPool.Exec(context.Background(), query,
		user_event.Event_Id,
		user_event.User_Id,
		user_event.IsRewarded,
		user_event.ClaimedReward,
	)

	if err != nil {
		return err
	}
	return nil

}

func GetAllEvents(dbPool *pgxpool.Pool) ([]models.Community_Event, error) {
	var events []models.Community_Event

	query := `SELECT id, name, percentage, code, capacity, date_time, manager_id, event_image, duration, local, current_ocupation, status, expiration_date
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
			&event.Status,
			&event.Expiration_Date,
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

	query := `SELECT id, name, percentage, code, capacity, date_time, manager_id, event_image, duration, local, current_ocupation, status, expiration_date
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
			&event.Status,
			&event.Expiration_Date,
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

func GetEventByCodeReward(code string, dbPool *pgxpool.Pool) (*models.Community_Event, error) {
	var event models.Community_Event

	query := `SELECT id, name, percentage, code, capacity, date_time, manager_id, event_image, duration, local, current_ocupation, status, expiration_date
	FROM events.community_event WHERE code = $1`

	err := dbPool.QueryRow(context.Background(), query, code).Scan(
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
		&event.Status,
		&event.Expiration_Date,
	)

	if err != nil {
		return nil, err
	}

	return &event, err

}

func GetEventById(event_id int, dbPool *pgxpool.Pool) (models.Community_Event, error) {
	var event models.Community_Event

	query := `SELECT id, name, percentage, code, capacity, date_time, manager_id, event_image, duration, local, current_ocupation, status, expiration_date
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
		&event.Status,
		&event.Expiration_Date,
	)

	if err != nil {
		return models.Community_Event{}, err
	}

	return event, nil
}

func GetEventsByUserId(user_id int, dbPool *pgxpool.Pool) ([]models.Community_Event, error) {
	var users_community_events []models.Community_Event

	query := `SELECT mcehmu.community_event_id
FROM events.many_community_event_has_many_users mcehmu
JOIN events.community_event ce ON mcehmu.community_event_id = ce.id
WHERE mcehmu.users_id = $1 AND ce.status = 'active'`

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

func DeleteEventById(event_id int, dbPool *pgxpool.Pool) error {
	query := `DELETE FROM events.community_event WHERE id = $1`

	_, err := dbPool.Exec(context.Background(), query, event_id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateRewardedStatus(event_id int, users_id int, dbPool *pgxpool.Pool) error {
	query := `UPDATE events.many_community_event_has_many_users SET isrewarded = true WHERE community_event_id = $1 AND users_id = $2`

	_, err := dbPool.Exec(context.Background(), query, event_id, users_id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateClaimedRewardStatus(event_id int, users_id int, dbPool *pgxpool.Pool) error {
	query := `UPDATE events.many_community_event_has_many_users SET claimedreward = true WHERE community_event_id = $1 AND users_id = $2`

	_, err := dbPool.Exec(context.Background(), query, event_id, users_id)
	if err != nil {
		return err
	}
	return nil
}

func GetAllUsersFromEventByEventId(event_id int, dbPool *pgxpool.Pool) ([]models.User_Event, error) {
	var users_events []models.User_Event

	query := `SELECT community_event_id, users_id, isrewarded, claimedreward
	FROM events.many_community_event_has_many_users WHERE community_event_id = $1`

	rows, err := dbPool.Query(context.Background(), query, event_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var users_event models.User_Event

		err := rows.Scan(
			&users_event.Event_Id,
			&users_event.User_Id,
			&users_event.IsRewarded,
			&users_event.ClaimedReward,
		)

		if err != nil {
			return nil, err
		}

		users_events = append(users_events, users_event)

	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return users_events, nil
}

func UpdateExpirationDate(date time.Time, event_id int, dbPool *pgxpool.Pool) error {
	query := `UPDATE events.community_event SET expiration_date = $1 WHERE id = $2`

	_, err := dbPool.Exec(context.Background(), query, date, event_id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateEventStatus(status string, event_id int, dbPool *pgxpool.Pool) error {
	query := `UPDATE events.community_event SET status = $1 WHERE id = $2`

	_, err := dbPool.Exec(context.Background(), query, status, event_id)
	if err != nil {
		return err
	}

	return nil
}

func GetAllFinishedEvents(dbPool *pgxpool.Pool) ([]models.Community_Event, error) {
	var events []models.Community_Event

	query := `SELECT id, name, percentage, code, capacity, date_time, manager_id, event_image, duration, local, current_ocupation, status, expiration_date
	FROM events.community_event WHERE status = 'finish'`

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
			&event.Status,
			&event.Expiration_Date,
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

func UpdateEventCurrentOcupation(current_ocupation int, event_id int, dbPool *pgxpool.Pool) error {
	query := `UPDATE events.community_event SET current_ocupation = $1 WHERE id = $2`

	_, err := dbPool.Exec(context.Background(), query, current_ocupation, event_id)
	if err != nil {
		return err
	}

	return nil
}
