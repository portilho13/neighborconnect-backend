package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/portilho13/neighborconnect-backend/repository/models/marketplace"
)

func CreateBid(bid models.Bid, dbPool *pgxpool.Pool) error {
	query := `INSERT INTO marketplace.bid (bid_ammount, bid_time, users_id, listing_id) VALUES ($1, $2, $3, $4)`
	_, err := dbPool.Exec(context.Background(), query,
		bid.Bid_Ammount,
		bid.Bid_Time,
		bid.User_Id,
		bid.Listing_Id,
	)

	if err != nil {
		return err
	}

	return nil
}

func GetBidByListningId(id int, dbPool *pgxpool.Pool) ([]models.Bid, error) {
	query := `SELECT id, bid_ammount, bid_time, users_id, listing_id FROM marketplace.bid WHERE listing_id = $1 ORDER BY bid_ammount DESC`

	rows, err := dbPool.Query(context.Background(), query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []models.Bid

	for rows.Next() {
		var bid models.Bid
		err := rows.Scan(
			&bid.Id,
			&bid.Bid_Ammount,
			&bid.Bid_Time,
			&bid.User_Id,
			&bid.Listing_Id,
		)

		if err != nil {
			return nil, err
		}

		bids = append(bids, bid)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return bids, nil

}
