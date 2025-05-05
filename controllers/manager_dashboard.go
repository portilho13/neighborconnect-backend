package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	controllers_models "github.com/portilho13/neighborconnect-backend/models"
	repositoryControllersEvents "github.com/portilho13/neighborconnect-backend/repository/controlers/events"
	repositoryControllersMarketplace "github.com/portilho13/neighborconnect-backend/repository/controlers/marketplace"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
)

func GetDashBoardInfo(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	manager_id_str := r.URL.Query().Get("manager_id")

	manager_id, err := strconv.Atoi(manager_id_str)
	if err != nil {
		http.Error(w, "Error Converting Manager Id", http.StatusInternalServerError)
		return
	}

	apartments, err := repositoryControllers.GetAllApartmentsByManagerId(manager_id, dbPool)
	if err != nil {
		http.Error(w, "Error Fetching apartments", http.StatusInternalServerError)
		return
	}

	var apartmentsJson []controllers_models.Apartment
	var usersJson []controllers_models.UserLogin
	var listingsJson []controllers_models.ListingInfo
	var eventsJson []controllers_models.EventInfo

	for _, apartment := range apartments {

		rents, err := repositoryControllers.GetRentByApartmentId(*apartment.Id, dbPool)
		if err != nil {
			http.Error(w, "Error Fetching Apartment Rent", http.StatusInternalServerError)
			return
		}

		last_rent := rents[0]

		last_rent_json := controllers_models.Rent{
			Id:           last_rent.Id,
			Month:        last_rent.Month,
			Year:         last_rent.Year,
			Base_Amount:  last_rent.Base_Amount,
			Reduction:    last_rent.Reduction,
			Final_Amount: last_rent.Final_Amount,
			Apartment_Id: last_rent.Apartment_Id,
			Status:       last_rent.Status,
			Due_Day:      last_rent.Due_day,
		}

		apartmentJson := controllers_models.Apartment{
			Id:         *apartment.Id,
			N_Bedrooms: apartment.N_bedrooms,
			Floor:      apartment.Floor,
			Rent:       int(apartment.Rent),
			Manager_Id: apartment.Manager_id,
			Status:     apartment.Status,
			Last_Rent:  last_rent_json,
		}
		apartmentsJson = append(apartmentsJson, apartmentJson)

		users, err := repositoryControllers.GetUsersByApartmentId(apartmentJson.Id, dbPool)

		if err != nil {
			http.Error(w, "Error Fetching Users", http.StatusInternalServerError)
			return
		}

		for _, user := range users {
			avatar := ""
			if user.Profile_Picture != nil {
				avatar = *user.Profile_Picture
			}

			userJson := controllers_models.UserLogin{
				Id:          user.Id,
				Name:        user.Name,
				Email:       user.Email,
				Phone:       user.Phone,
				ApartmentID: *user.Apartment_id,
				Avatar:      avatar,
			}

			usersJson = append(usersJson, userJson)

			listings, err := repositoryControllersMarketplace.GetListingsBySellerId(user.Id, dbPool)
			if err != nil {
				http.Error(w, "Error Fetching Listings", http.StatusInternalServerError)
				return
			}

			for _, listing := range listings {
				category, err := repositoryControllersMarketplace.GetCategoryById(*listing.Category_Id, dbPool)
				if err != nil {
					http.Error(w, "Error Fetching Category", http.StatusInternalServerError)
					return
				}

				categoryJson := controllers_models.CategoryInfo{
					Id:   *category.Id,
					Name: category.Name,
					Url:  *category.Url,
				}

				user, err := repositoryControllers.GetUsersById(*listing.Seller_Id, dbPool)
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

				listingsJson = append(listingsJson, listingJson)
			}

		}

	}

	events, err := repositoryControllersEvents.GetAllEventsByManagerId(manager_id, dbPool)
	if err != nil {
		http.Error(w, "Error Fetching Events", http.StatusInternalServerError)
		return
	}

	for _, event := range events {
		eventImage := ""
		if event.Event_Image != nil {
			eventImage = *event.Event_Image
		}
		eventsJson = append(eventsJson, controllers_models.EventInfo{
			Id:                *event.Id,
			Name:              event.Name,
			Percentage:        event.Percentage,
			Capacity:          event.Capacity,
			Date_time:         event.Date_Time,
			Manager_Id:        *event.Manager_Id,
			Event_Image:       eventImage,
			Duration:          event.Duration,
			Local:             event.Local,
			Current_Ocupation: event.Current_Ocupation,
		})
	}

	dashboardInfo := controllers_models.ManagerDashboardInfo{
		Apartments: apartmentsJson,
		Users:      usersJson,
		Listings:   listingsJson,
		Events:     eventsJson,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dashboardInfo); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}
