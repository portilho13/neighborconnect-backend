package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/marketplace"
)

func CreateListingPhotos(item_photos models.Listing_Photos, dbPool *pgxpool.Pool) error {

	query := `INSERT INTO marketplace.listing_photos (url, listing_id) 
	VALUES ($1, $2)`

	_, err := dbPool.Exec(context.Background(), query,
		item_photos.Url,
		item_photos.Listing_Id,
	)

	if err != nil {
		return err
	}

	return nil
}

func GetListingPhotosByListingId(id int, dbPool *pgxpool.Pool) ([]models.Listing_Photos, error) {
	query := `SELECT id, url, listing_id FROM marketplace.listing_photos WHERE listing_id = $1`
	rows, err := dbPool.Query(context.Background(), query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var item_photos []models.Listing_Photos

	for rows.Next() {
		var item_photo models.Listing_Photos
		err := rows.Scan(
			&item_photo.Id,
			&item_photo.Url,
			&item_photo.Listing_Id,
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
