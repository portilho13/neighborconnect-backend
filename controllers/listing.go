package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	controllers_models "github.com/portilho13/neighborconnect-backend/models"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/marketplace"
	models "github.com/portilho13/neighborconnect-backend/repository/models/marketplace"
)

func CreateListing(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	var listing controllers_models.ListingCreation
	err := json.NewDecoder(r.Body).Decode(&listing)

	if err != nil {
		http.Error(w, "Invalid JSON Data", http.StatusBadRequest)
		return
	}

	var sellerID *int

	if listing.Seller_Id == 0 { // For testing purposes only
		sellerID = nil
	} else {
		sellerID = &listing.Seller_Id
	}

	listingDB := models.Listing{
		Name:            listing.Name,
		Description:     listing.Description,
		Buy_Now_Price:   listing.Buy_Now_Price,
		Start_Price:     listing.Start_Price,
		Created_At:      time.Now(),
		Expiration_Time: listing.Expiration_Time,
		Status:          "active", // When a listing is created status will be active by default
		Seller_Id:       sellerID,
	}

	err = repositoryControllers.CreateListing(listingDB, dbPool)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Error creating item", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Listing Created !"})
}

func GetListingById(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	query := r.URL.Query()
	idStr := query.Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	listing, err := repositoryControllers.GetListingById(id, dbPool)
	if err != nil {
		http.Error(w, "Invalid listing", http.StatusBadRequest)
		return
	}

	listingJson := controllers_models.ListingInfo{
		Name:            listing.Name,
		Description:     listing.Description,
		Buy_Now_Price:   listing.Buy_Now_Price,
		Start_Price:     listing.Start_Price,
		Created_At:      listing.Created_At,
		Expiration_Time: listing.Expiration_Time,
		Status:          listing.Status,
		Seller_Id:       listing.Seller_Id,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(listingJson); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
