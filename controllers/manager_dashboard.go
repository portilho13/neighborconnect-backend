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
		apartmentJson := controllers_models.Apartment{
			Id:         *apartment.Id,
			N_Bedrooms: apartment.N_bedrooms,
			Floor:      apartment.Floor,
			Rent:       int(apartment.Rent),
			Manager_Id: apartment.Manager_id,
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
				listingJson := controllers_models.ListingInfo {
					Id: *listing.Id,
					Name: listing.Name,
					Description: listing.Description,
					Buy_Now_Price: listing.Buy_Now_Price,
					Start_Price: listing.Start_Price,
					Created_At: listing.Created_At,
					Expiration_Time: listing.Expiration_Time,
					Status: listing.Status,
					Seller_Id: listing.Seller_Id,
					Category_Id: listing.Category_Id,
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
		Users: usersJson,
		Listings: listingsJson,
		Events: eventsJson,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dashboardInfo); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}
