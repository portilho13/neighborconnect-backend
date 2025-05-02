package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	controllers_models "github.com/portilho13/neighborconnect-backend/models"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/marketplace"
	models "github.com/portilho13/neighborconnect-backend/repository/models/marketplace"
)

func CreateListing(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	// Extract form values
	name := r.FormValue("name")
	description := r.FormValue("description")

	buyNowPriceStr := strings.Trim(r.FormValue("buy_now_price"), "\"")
	buyNowPrice, err := strconv.Atoi(buyNowPriceStr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid buy_now_price", http.StatusBadRequest)
		return
	}

	// Sanitize string values by trimming surrounding quotes, if present
	startPriceStr := strings.Trim(r.FormValue("start_price"), "\"")
	startPrice, err := strconv.Atoi(startPriceStr)
	fmt.Println("Raw:", r.FormValue("start_price"))
	fmt.Println("Trimmed:", startPriceStr)
	if err != nil {
		http.Error(w, "Invalid start_price", http.StatusBadRequest)
		return
	}

	expirationTimeStr := strings.Trim(r.FormValue("expiration_time"), "\"")
	expirationDate, err := time.Parse(time.RFC3339, expirationTimeStr)
	if err != nil {
		http.Error(w, "Invalid expiration_time format. Use RFC3339", http.StatusBadRequest)
		return
	}

	sellerIDStr := strings.Trim(r.FormValue("seller_id"), "\"")
	sellerID, err := strconv.Atoi(sellerIDStr)
	if err != nil {
		sellerID = 0 // fallback for testing
	}

	categoryIDStr := strings.Trim(r.FormValue("category_id"), "\"")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		categoryID = 0 // fallback for testing
	}

	// Save uploaded files
	files := r.MultipartForm.File["images"]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			log.Println("Error opening uploaded file:", err)
			continue
		}
		defer file.Close()

		// You can change this path to wherever you want to store uploaded files
		savePath := fmt.Sprintf("./uploads/%s", fileHeader.Filename)
		dst, err := os.Create(savePath)
		if err != nil {
			log.Println("Error creating destination file:", err)
			continue
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			log.Println("Error saving uploaded file:", err)
			continue
		}
	}

	// Build and insert listing into DB
	var sellerIDPtr *int
	if sellerID != 0 {
		sellerIDPtr = &sellerID
	}

	listingDB := models.Listing{
		Name:            name,
		Description:     description,
		Buy_Now_Price:   buyNowPrice,
		Start_Price:     startPrice,
		Created_At:      time.Now(),
		Expiration_Date: expirationDate,
		Status:          "active",
		Seller_Id:       sellerIDPtr,
		Category_Id:     &categoryID,
	}

	err = repositoryControllers.CreateListing(listingDB, dbPool)
	if err != nil {
		log.Println("Database error:", err)
		http.Error(w, "Failed to create listing", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Listing Created!"})
}

// func CreateListing(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
// 	var listing controllers_models.ListingCreation
// 	err := json.NewDecoder(r.Body).Decode(&listing)

// 	if err != nil {
// 		http.Error(w, "Invalid JSON Data", http.StatusBadRequest)
// 		return
// 	}

// 	var sellerID *int

// 	if listing.Seller_Id == 0 { // For testing purposes only
// 		sellerID = nil
// 	} else {
// 		sellerID = &listing.Seller_Id
// 	}

// 	listingDB := models.Listing{
// 		Name:            listing.Name,
// 		Description:     listing.Description,
// 		Buy_Now_Price:   listing.Buy_Now_Price,
// 		Start_Price:     listing.Start_Price,
// 		Created_At:      time.Now(),
// 		Expiration_Date: listing.Expiration_Date,
// 		Status:          "active", // When a listing is created status will be active by default
// 		Seller_Id:       sellerID,
// 	}

// 	err = repositoryControllers.CreateListing(listingDB, dbPool)
// 	if err != nil {
// 		log.Fatal(err)
// 		http.Error(w, "Error creating item", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(map[string]string{"message": "Listing Created !"})
// }

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
		Expiration_Date: listing.Expiration_Date,
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

func GetAllListings(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	listings, err := repositoryControllers.GetAllActiveListings(dbPool)
	if err != nil {
		http.Error(w, "Failed fetch listings", http.StatusInternalServerError)
		return
	}

	var listingsJson []controllers_models.ListingInfo

	for _, listing := range listings {
		bids, err := repositoryControllers.GetBidByListningId(*listing.Id, dbPool)
		var highestBid int
		if len(bids) == 0 {
			highestBid = listing.Start_Price
		} else {
			highestBid = bids[0].Bid_Ammount
		}
		if err != nil {
			http.Error(w, "Failed fetch listings", http.StatusInternalServerError)
			return
		}
		listingsJson = append(listingsJson, controllers_models.ListingInfo{
			Name:            listing.Name,
			Description:     listing.Description,
			Buy_Now_Price:   listing.Buy_Now_Price,
			Start_Price:     listing.Start_Price,
			Current_bid:     highestBid,
			Created_At:      listing.Created_At,
			Expiration_Date: listing.Expiration_Date,
			Status:          listing.Status,
			Seller_Id:       listing.Seller_Id,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(listingsJson); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
