package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/events"
)

func CreateCommunityEvent(event models.Community_Event, dbPool *pgxpool.Pool) error {
	query := `INSERT INTO events.community_event (name, percentage, code, capacity, date_time, manager_id, event_image, duration) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := dbPool.Exec(context.Background(), query,
		event.Name,
		event.Percentage,
		event.Code,
		event.Capacity,
		event.Date_Time,
		event.Manager_Id,
		event.Event_Image,
		event.Duration,
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
