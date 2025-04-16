package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	controllers_models "github.com/portilho13/neighborconnect-backend/models"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/marketplace"
	models "github.com/portilho13/neighborconnect-backend/repository/models/marketplace"
)

func CreateBid(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	var bidJSON controllers_models.BidJson
	err := json.NewDecoder(r.Body).Decode(&bidJSON)

	if err != nil {
		http.Error(w, "Invalid JSON Data", http.StatusBadRequest)
		return
	}

	bid := models.Bid{
		Bid_Ammount: bidJSON.Bid_Ammount,
		Bid_Time:    time.Now(),
		User_Id:     bidJSON.User_Id,
		Listing_Id:  bidJSON.Listing_Id,
	}

	_, err = repositoryControllers.GetListingById(*bidJSON.Listing_Id, dbPool)
	if err != nil {
		http.Error(w, "Invalid Listing", http.StatusBadRequest)
		return
	}

	err = repositoryControllers.CreateBid(bid, dbPool)
	if err != nil {
		http.Error(w, "Error creating bid", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Bid Created !"})

}
