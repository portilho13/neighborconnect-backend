package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	controllers_models "github.com/portilho13/neighborconnect-backend/models"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/marketplace"
	models "github.com/portilho13/neighborconnect-backend/repository/models/marketplace"
	"github.com/portilho13/neighborconnect-backend/utils"
	"github.com/portilho13/neighborconnect-backend/ws"
)

func CreateBid(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	var bidJSON controllers_models.BidJson
	err := json.NewDecoder(r.Body).Decode(&bidJSON)

	if err != nil {
		http.Error(w, "Invalid JSON Data", http.StatusBadRequest)
		return
	}

	//Dont need to fetch all active listings because frontend will only send active listings
	listing, err := repositoryControllers.GetListingById(*bidJSON.Listing_Id, dbPool) // Check if listing is valid
	if err != nil {
		http.Error(w, "Invalid Listing", http.StatusBadRequest)
		return
	}

	// Double check if status is active due to someone crawling api
	if listing.Status != "active" {
		http.Error(w, "Invalid Listing", http.StatusBadRequest)
		return
	}

	//Lazy check is listing is over
	timeNow := time.Now()
	if timeNow.After(listing.Expiration_Time) {
		if listing.Status == "active" {
			err = utils.CloseListingBid(*listing.Id, dbPool)
			if err != nil {
				http.Error(w, "Invalid Listing", http.StatusBadRequest)
				return
			}
		}
		http.Error(w, "Listing is Closed", http.StatusBadRequest)
		return

	}

	bids, err := repositoryControllers.GetBidByListningId(*bidJSON.Listing_Id, dbPool) // Get current bids for listings
	if err != nil {
		http.Error(w, "Error getting bids", http.StatusBadRequest)
		return
	}

	var highestBid int

	if len(bids) > 0 {
		highestBid = bids[0].Bid_Ammount
	} else {
		highestBid = 0
	}

	if bidJSON.Bid_Ammount > highestBid || len(bids) == 0 { // Only accept if bid ammount is bigger than highest bid
		bid := models.Bid{
			Bid_Ammount: bidJSON.Bid_Ammount,
			Bid_Time:    time.Now(),
			User_Id:     bidJSON.User_Id,
			Listing_Id:  bidJSON.Listing_Id,
		}

		err = repositoryControllers.CreateBid(bid, dbPool)
		if err != nil {
			http.Error(w, "Error creating bid", http.StatusBadRequest)
			return
		}

		highestBidJSON, err := json.Marshal(bidJSON.Bid_Ammount)
		if err != nil {
			http.Error(w, "Error converting", http.StatusBadRequest)
			return
		}

		ws.Hub.Broadcast <- highestBidJSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Bid Created !"})
	} else {
		http.Error(w, "Biding lower than highest", http.StatusBadRequest)
		return
	}

}
