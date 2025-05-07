package utils

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/portilho13/neighborconnect-backend/email"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/marketplace"
	repositoryControllersUsers "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
	models "github.com/portilho13/neighborconnect-backend/repository/models/marketplace"
	"github.com/robfig/cron/v3"
)

func CloseListing(id int, dbPool *pgxpool.Pool) error {
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

	if len(bids) == 0 {
		return nil // Nobody Bidded
	}

	highestBidder := bids[0]

	transaction := models.Transaction{
		Final_Price:      highestBidder.Bid_Ammount,
		Transaction_Time: time.Now(),
		Transaction_Type: "bid",
		Seller_Id:        sellerId,
		Buyer_Id:         highestBidder.User_Id,
		Listing_Id:       listing.Id,
		Payment_Status:   "pending",
		Payment_Due_time: time.Now().UTC().AddDate(0, 0, 5),
	}

	transaction_id, err := repositoryControllers.CreateTransaction(transaction, dbPool)
	if err != nil {
		return err
	}

	user, err := repositoryControllersUsers.GetUsersById(*highestBidder.User_Id, dbPool)
	if err != nil {
		return err
	}

	listing.Id = &transaction_id // Add Transaction Id

	email_struct := email.Email{
		To:      []string{user.Email},
		Subject: "Order Confirmation",
	}

	err = email.SendEmail(email_struct, "order confirmation", transaction)
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
		Payment_Status:   "pending",
		Payment_Due_time: time.Now().UTC().AddDate(0, 0, 5),
	}

	transaction_id, err := repositoryControllers.CreateTransaction(transaction, dbPool)
	if err != nil {
		return err
	}

	user, err := repositoryControllersUsers.GetUsersById(buyerId, dbPool)
	if err != nil {
		return err
	}

	email_struct := email.Email{
		To:      []string{user.Email},
		Subject: "Order Confirmation",
	}

	transaction.Id = &transaction_id // Add Transaction Id

	err = email.SendEmail(email_struct, "order confirmation", transaction)
	if err != nil {
		return err
	}

	return nil
}

func closeExpiredListings(dbPool *pgxpool.Pool) error {
	listings, err := repositoryControllers.GetAllActiveListings(dbPool)
	if err != nil {
		return err
	}

	for _, listing := range listings {
		if time.Now().After(listing.Expiration_Date) {
			err = CloseListing(*listing.Id, dbPool)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func AutomateListingClosing(dbPool *pgxpool.Pool) {
	c := cron.New(cron.WithSeconds())
	c.AddFunc("@every 1s", func() { // Maybe change this ???
		_ = closeExpiredListings(dbPool)

	})

	c.Start()

	defer c.Stop()

	select {}
}
