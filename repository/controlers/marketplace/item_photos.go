package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/marketplace"
)

func CreateItemPhotos(item_photos models.Item_Photos, dbPool *pgxpool.Pool) error {

	query := `INSERT INTO marketplace.item_photos (url, item_id) 
	VALUES ($1, $2)`

	_, err := dbPool.Exec(context.Background(), query,
		item_photos.Url,
		item_photos.Item_Id,
	)

	if err != nil {
		return nil
	}

	return nil
}

func GetItemPhotosById(id int, dbPool *pgxpool.Pool) ([]models.Item_Photos, error) {
	query := `SELECT (id, url, item_id) FROM marketplace.item_photos WHERE item_id = $1`
	rows, err := dbPool.Query(context.Background(), query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var item_photos []models.Item_Photos

	for rows.Next() {
		var item_photo models.Item_Photos
		err := rows.Scan(
			&item_photo.Id,
			&item_photo.Url,
			&item_photo.Item_Id,
		)

		if err != nil {
			return nil, err
		}

		item_photos = append(item_photos, item_photo)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return item_photos, nil

}
