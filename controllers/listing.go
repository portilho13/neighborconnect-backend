package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	controllers_models "github.com/portilho13/neighborconnect-backend/models"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/marketplace"
	repositoryControllersUsers "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
	models "github.com/portilho13/neighborconnect-backend/repository/models/marketplace"
	"github.com/portilho13/neighborconnect-backend/utils"
)

func CreateListing(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	err := r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	listingJSON := r.FormValue("listing")
	if listingJSON == "" {
		http.Error(w, "Missing listing data", http.StatusBadRequest)
		return
	}

	var listingData controllers_models.ListingCreation
	if err := json.Unmarshal([]byte(listingJSON), &listingData); err != nil {
		http.Error(w, "Invalid listing data format", http.StatusBadRequest)
		return
	}

	// Convert string values to appropriate types
	buyNowPrice, err := strconv.Atoi(listingData.Buy_Now_Price)
	if err != nil {
		http.Error(w, "Invalid buy_now_price", http.StatusBadRequest)
		return
	}

	startPrice, err := strconv.Atoi(listingData.Start_Price)
	if err != nil {
		http.Error(w, "Invalid start_price", http.StatusBadRequest)
		return
	}

	// Parse the expiration time from string to time.Time
	var expirationDate time.Time
	if listingData.Expiration_Date != "" {
		expirationDate, err = time.Parse(time.RFC3339, listingData.Expiration_Date)
		if err != nil {
			http.Error(w, "Invalid expiration_time format. Use RFC3339", http.StatusBadRequest)
			return
		}
	}

	sellerID, err := strconv.Atoi(listingData.Seller_Id)
	if err != nil {
		http.Error(w, "Invalid seller id", http.StatusBadRequest)
		return
	}

	categoryID, err := strconv.Atoi(listingData.Category_Id)
	if err != nil {
		http.Error(w, "Invalid category id", http.StatusBadRequest)
		return
	}

	var sellerIDPtr *int
	if sellerID != 0 {
		sellerIDPtr = &sellerID
	}

	listingDB := models.Listing{
		Name:            listingData.Name,
		Description:     listingData.Description,
		Buy_Now_Price:   buyNowPrice,
		Start_Price:     startPrice,
		Created_At:      time.Now(),
		Expiration_Date: expirationDate,
		Status:          "active",
		Seller_Id:       sellerIDPtr,
		Category_Id:     &categoryID,
	}

	id, err := repositoryControllers.CreateListingReturningId(listingDB, dbPool)
	if err != nil {
		http.Error(w, "Failed to create listing", http.StatusInternalServerError)
		return
	}

	files := r.MultipartForm.File["images"]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Failed to create listing", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		filename := fileHeader.Filename
		ext := filepath.Ext(filename)
		newFilename := uuid.New().String() + ext
		savePath := fmt.Sprintf("./uploads/listing/%s", newFilename)

		dst, err := os.Create(savePath)
		if err != nil {
			http.Error(w, "Failed to create listing", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Failed to create listing", http.StatusInternalServerError)
			return
		}

		api_url := utils.GetApiUrl()

		api_path := fmt.Sprintf("http://%s/api/v1/uploads/listing/%s", api_url, newFilename)

		listing_photo := models.Listing_Photos{
			Url:        api_path,
			Listing_Id: id,
		}

		err = repositoryControllers.CreateListingPhotos(listing_photo, dbPool)
		if err != nil {
			http.Error(w, "Failed to create listing photos", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Listing Created!"})
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

	photos, err := repositoryControllers.GetListingPhotosByListingId(*listing.Id, dbPool)
	if err != nil {
		http.Error(w, "Failed fetching photo listings", http.StatusInternalServerError)
		return
	}

	var listing_photos_json []controllers_models.Listing_Photos

	for _, photo := range photos {
		photo_json := controllers_models.Listing_Photos{
			Id:  photo.Id,
			Url: photo.Url,
		}

		listing_photos_json = append(listing_photos_json, photo_json)
	}

	category, err := repositoryControllers.GetCategoryById(*listing.Category_Id, dbPool)
	if err != nil {
		http.Error(w, "Error Fetching Category", http.StatusInternalServerError)
		return
	}

	categoryJson := controllers_models.CategoryInfo{
		Id:   *category.Id,
		Name: category.Name,
		Url:  *category.Url,
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

	bids, err := repositoryControllers.GetBidByListningId(*listing.Id, dbPool)
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
		bidJson.Bid_Time = &highestBid.Bid_Time
		bidJson.User_Id = highestBid.User_Id
		bidJson.Listing_Id = *highestBid.Listing_Id
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
		Listing_Photos:  listing_photos_json,
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

		if err != nil {
			http.Error(w, "Failed fetching bids", http.StatusInternalServerError)
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
			bidJson.Bid_Time = &highestBid.Bid_Time
			bidJson.User_Id = highestBid.User_Id
			bidJson.Listing_Id = *highestBid.Listing_Id
		}

		photos, err := repositoryControllers.GetListingPhotosByListingId(*listing.Id, dbPool)
		if err != nil {
			http.Error(w, "Failed fetching photo listings", http.StatusInternalServerError)
			return
		}

		var listing_photos_json []controllers_models.Listing_Photos

		for _, photo := range photos {
			photo_json := controllers_models.Listing_Photos{
				Id:  photo.Id,
				Url: photo.Url,
			}

			listing_photos_json = append(listing_photos_json, photo_json)
		}

		category, err := repositoryControllers.GetCategoryById(*listing.Category_Id, dbPool)
		if err != nil {
			http.Error(w, "Error Fetching Category", http.StatusInternalServerError)
			return
		}

		categoryJson := controllers_models.CategoryInfo{
			Id:   *category.Id,
			Name: category.Name,
			Url:  *category.Url,
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

		listingsJson = append(listingsJson, controllers_models.ListingInfo{
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
			Listing_Photos:  listing_photos_json,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(listingsJson); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
