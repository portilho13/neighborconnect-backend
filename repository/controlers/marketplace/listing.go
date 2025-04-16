package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/marketplace"
)

func CreateListing(listing models.Listing, dbPool *pgxpool.Pool) error {
	query := `INSERT INTO marketplace.listing 
	(name, description, buy_now_price, start_price, created_at, expiration_time, status, seller_id) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	fmt.Println(listing.Created_At)

	_, err := dbPool.Exec(context.Background(), query,
		listing.Name,
		listing.Description,
		listing.Buy_Now_Price,
		listing.Start_Price,
		listing.Created_At,
		listing.Expiration_Time,
		listing.Status,
		listing.Seller_Id,
	)

	if err != nil {
		return err
	}

	return nil
}

func GetListingById(id int, dbPool *pgxpool.Pool) (models.Listing, error) {
	query := `SELECT 
	id, name, description, buy_now_price, start_price, created_at, expiration_time, status, seller_id
	FROM marketplace.listing
	WHERE id = $1`

	var listing models.Listing

	err := dbPool.QueryRow(context.Background(), query, id).Scan(
		&listing.Id,
		&listing.Name,
		&listing.Description,
		&listing.Buy_Now_Price,
		&listing.Start_Price,
		&listing.Created_At,
		&listing.Expiration_Time,
		&listing.Status,
		&listing.Seller_Id,
	)

	if err != nil {
		return models.Listing{}, err
	}

	return listing, nil
}

func GetAllListings(dbPool *pgxpool.Pool) ([]models.Listing, error) {
	query := `SELECT 
	id, name, description, buy_now_price, start_price, created_at, expiration_time, status, seller_id
	FROM marketplace.listing`

	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var listings []models.Listing

	for rows.Next() {
		var listing models.Listing

		err := rows.Scan(
			&listing.Id,
			&listing.Name,
			&listing.Description,
			&listing.Buy_Now_Price,
			&listing.Start_Price,
			&listing.Created_At,
			&listing.Expiration_Time,
			&listing.Status,
			&listing.Seller_Id,
		)

		if err != nil {
			return nil, err
		}

		listings = append(listings, listing)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return listings, nil
}
