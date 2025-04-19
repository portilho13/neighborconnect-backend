package utils

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/marketplace"
	models "github.com/portilho13/neighborconnect-backend/repository/models/marketplace"
)

func CloseListingBid(id int, dbPool *pgxpool.Pool) error {
	err := repositoryControllers.UpdateListingStatus("closed", id, dbPool)
	if err != nil {
		return err
	}

	listing, err := repositoryControllers.GetListingById(id, dbPool)
	if err != nil {
		return err
	}

	sellerId := listing.Seller_Id

	bids, err := repositoryControllers.GetBidByListningId(id, dbPool)
	if err != nil {
		return err
	}

	highestBidder := bids[0]

	transaction := models.Transaction{
		Final_Price:      highestBidder.Bid_Ammount,
		Transaction_Time: time.Now(),
		Transaction_Type: "bid",
		Seller_Id:        sellerId,
		Buyer_Id:         highestBidder.User_Id,
		Listing_Id:       listing.Id,
	}

	err = repositoryControllers.CreateTransaction(transaction, dbPool)
	if err != nil {
		return err
	}

	return nil
}

func CloseListingBuy(listingId int, buyerId int, dbPool *pgxpool.Pool) error {
	err := repositoryControllers.UpdateListingStatus("closed", listingId, dbPool)
	if err != nil {
		return err
	}

	listing, err := repositoryControllers.GetListingById(listingId, dbPool)
	if err != nil {
		return err
	}

	sellerId := listing.Seller_Id

	transaction := models.Transaction{
		Final_Price:      listing.Buy_Now_Price,
		Transaction_Time: time.Now(),
		Transaction_Type: "buy",
		Seller_Id:        sellerId,
		Buyer_Id:         &buyerId,
		Listing_Id:       listing.Id,
	}

	err = repositoryControllers.CreateTransaction(transaction, dbPool)
	if err != nil {
		return err
	}

	return nil
}
