package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	controllers_models "github.com/portilho13/neighborconnect-backend/models"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/marketplace"
	"github.com/portilho13/neighborconnect-backend/utils"
)

func CreateBuyOrder(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	var buyJson controllers_models.BuyJson
	err := json.NewDecoder(r.Body).Decode(&buyJson)

	if err != nil {
		http.Error(w, "Invalid JSON Data", http.StatusBadRequest)
		return
	}

	listing, err := repositoryControllers.GetListingById(buyJson.Listing_Id, dbPool)
	if err != nil {
		http.Error(w, "Invalid Listing Data", http.StatusBadRequest)
		return
	}

	timeNow := time.Now()

	if timeNow.After(listing.Expiration_Date) {
		if listing.Status == "active" {
			err = utils.CloseListingBuy(*listing.Id, buyJson.User_Id, dbPool)
			if err != nil {
				http.Error(w, "Invalid Listing", http.StatusBadRequest)
				return
			}
		}
		http.Error(w, "Listing is Closed", http.StatusBadRequest)
		return
	}

	err = utils.CloseListingBuy(*listing.Id, buyJson.User_Id, dbPool)
	if err != nil {
		http.Error(w, "Invalid Listing", http.StatusBadRequest)
		return
	}
}
