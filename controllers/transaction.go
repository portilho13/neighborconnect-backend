package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/portilho13/neighborconnect-backend/email"
	controllers_models "github.com/portilho13/neighborconnect-backend/models"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/marketplace"
	repositoryControllersMarketplace "github.com/portilho13/neighborconnect-backend/repository/controlers/marketplace"
	repositoryControllersUsers "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"

	"github.com/portilho13/neighborconnect-backend/utils"
)

const FEES = 0.05 //5%, maybe then put in db ??

func PayTransaction(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	var payJson controllers_models.PayJson

	err := json.NewDecoder(r.Body).Decode(&payJson)
	if err != nil {
		http.Error(w, "Invalid JSON Data", http.StatusBadRequest)
		return
	}

	final_price := 0.0
	var f_price float64

	for _, transaction_id := range payJson.Transaction_Ids {
		transaction, err := repositoryControllers.GetTransactionById(transaction_id, dbPool)
		if err != nil {
			http.Error(w, "Error Fetching Transaction", http.StatusInternalServerError)
			return
		}

		if transaction.Payment_Status != "pending" {
			http.Error(w, "Invalid Transaction", http.StatusInternalServerError)
			return
		}

		if *transaction.Buyer_Id != payJson.User_Id {
			http.Error(w, "Invalid User", http.StatusInternalServerError)
			return
		}
		f_price = transaction.Final_Price
		final_price += transaction.Final_Price * (1 + FEES)

	}

	switch payJson.Type {
	case "wallet":
		if !utils.ValidateWalletBalance(payJson.User_Id, final_price, dbPool) {
			http.Error(w, "Insuficient Balance", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Unsuported Method", http.StatusInternalServerError)
		return
	}

	err = repositoryControllersUsers.UpdateAccountBalance(payJson.User_Id, final_price, dbPool)
	if err != nil {
		http.Error(w, "Error Adding Payout", http.StatusInternalServerError)
		return
	}

	for _, transaction_id := range payJson.Transaction_Ids {

		err = repositoryControllers.UpdateTransactionStatus("paid", transaction_id, dbPool)
		if err != nil {
			http.Error(w, "Error Changing Transaction Status", http.StatusInternalServerError)
			return
		}

		user, err := repositoryControllersUsers.GetUsersById(payJson.User_Id, dbPool)
		if err != nil {
			http.Error(w, "Error Getting User Id", http.StatusInternalServerError)
			return
		}

		transaction, err := repositoryControllers.GetTransactionById(transaction_id, dbPool)
		if err != nil {
			http.Error(w, "Error Fetching Transaction", http.StatusInternalServerError)
			return
		}

		email_struct := email.Email{
			To:      []string{user.Email},
			Subject: "Order Receipt",
		}

		err = email.SendEmail(email_struct, "order receipt", transaction)
		if err != nil {
			http.Error(w, "Error Sending Email", http.StatusInternalServerError)
			return
		}

		err = repositoryControllersUsers.UpdateAccountBalance(*transaction.Seller_Id, transaction.Final_Price, dbPool)
		if err != nil {
			http.Error(w, "Error Updating Seller Balance", http.StatusInternalServerError)
			return
		}

	}

	switch payJson.Type {
	case "wallet":
		if !utils.ValidateWalletBalance(payJson.User_Id, final_price, dbPool) {
			http.Error(w, "Insuficient Balance", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Unsuported Method", http.StatusInternalServerError)
		return
	}

	err = repositoryControllersUsers.UpdateAccountBalance(payJson.User_Id, final_price, dbPool)
	if err != nil {
		http.Error(w, "Error Adding Payout", http.StatusInternalServerError)
		return
	}

	manager_id, err := utils.GetManagerIdByUserId(payJson.User_Id, dbPool)
	if err != nil {
		http.Error(w, "Error fetching manager id", http.StatusInternalServerError)
		return
	}

	manager_transaction := models.Manager_Transaction{
		Type:        "Fee",
		Amount:      f_price * FEES,
		Date:        time.Now().UTC(),
		Description: "Marketplace Fee",
		Users_Id:    &payJson.User_Id,
		Manager_Id:  *manager_id,
	}

	err = repositoryControllersUsers.CreateManagerTransaction(manager_transaction, dbPool)
	if err != nil {
		http.Error(w, "Error creating manager transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Pay Concluded Sucessfully !"})
}

func GetAllTransactions(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	query := r.URL.Query()
	idStr := query.Get("user_id")

	user_id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	transactions, err := repositoryControllers.GetTransactionsByBuyerId(user_id, dbPool)
	if err != nil {
		http.Error(w, "Error Fetching User", http.StatusInternalServerError)
		return
	}

	var transactionsJson []controllers_models.TransactionInfoJson

	for _, transaction := range transactions {

		listing, err := repositoryControllers.GetListingById(*transaction.Listing_Id, dbPool)
		if err != nil {
			http.Error(w, "Error Fetching Transaction", http.StatusInternalServerError)
			return
		}

		category, err := repositoryControllersMarketplace.GetCategoryById(*listing.Category_Id, dbPool)
		if err != nil {
			http.Error(w, "Error Fetching Category", http.StatusInternalServerError)
			return
		}

		categoryJson := controllers_models.CategoryInfo{
			Id:   *category.Id,
			Name: category.Name,
		}

		user, err := repositoryControllersUsers.GetUsersById(*listing.Seller_Id, dbPool)
		if err != nil {
			http.Error(w, "Error Fetching Seller Info", http.StatusInternalServerError)
			return
		}

		userJson := controllers_models.SellerListingInfo{
			Id:   user.Id,
			Name: user.Name,
		}

		bids, err := repositoryControllersMarketplace.GetBidByListningId(*listing.Id, dbPool)
		if err != nil {
			http.Error(w, "Error Fetching Bids", http.StatusInternalServerError)
			return
		}

		var bidJson controllers_models.BidInfo
		if len(bids) == 0 {
			bidJson.Id = nil
			bidJson.Bid_Ammount = listing.Start_Price
			bidJson.Bid_Time = nil
			bidJson.User_Id = nil
			bidJson.Listing_Id = *listing.Id
		} else {
			highestBid := bids[0]

			bidJson.Id = highestBid.Id
			bidJson.Bid_Ammount = highestBid.Bid_Ammount
			bidJson.User_Id = highestBid.User_Id
			bidJson.Bid_Time = &highestBid.Bid_Time
			bidJson.Listing_Id = *highestBid.Listing_Id
		}

		listing_photos, err := repositoryControllersMarketplace.GetListingPhotosByListingId(*listing.Id, dbPool)
		if err != nil {
			http.Error(w, "Error Fetching Listing Photos", http.StatusInternalServerError)
			return
		}

		var listingPhotosJson []controllers_models.Listing_Photos

		for _, photo := range listing_photos {
			listingPhotosJson = append(listingPhotosJson, controllers_models.Listing_Photos{
				Id:  photo.Id,
				Url: photo.Url,
			})
		}

		listingJson := controllers_models.ListingInfo{
			Id:              *listing.Id,
			Name:            listing.Name,
			Description:     listing.Description,
			Buy_Now_Price:   listing.Buy_Now_Price,
			Start_Price:     listing.Start_Price,
			Current_bid:     bidJson,
			Created_At:      listing.Created_At,
			Expiration_Date: listing.Expiration_Date,
			Status:          listing.Status,
			Seller:          userJson,
			Category:        categoryJson,
			Listing_Photos:  listingPhotosJson,
		}

		transactionJson := controllers_models.TransactionInfoJson{
			Id:               transaction.Id,
			Final_Price:      transaction.Final_Price,
			Transaction_Time: transaction.Transaction_Time,
			Transaction_Type: transaction.Transaction_Type,
			Seller_Id:        transaction.Seller_Id,
			Buyer_Id:         transaction.Buyer_Id,
			Listing:          listingJson,
			Payment_Status:   transaction.Payment_Status,
			Payment_Due_time: transaction.Payment_Due_time,
		}

		transactionsJson = append(transactionsJson, transactionJson)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(transactionsJson); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}
