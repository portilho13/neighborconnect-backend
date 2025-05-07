package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	controllers_models "github.com/portilho13/neighborconnect-backend/models"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/events"
	repositoryControllersUsers "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
)

func Reward(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	var reward controllers_models.RewardJson

	err := json.NewDecoder(r.Body).Decode(&reward)
	if err != nil {
		http.Error(w, "Invalid JSON Data", http.StatusBadRequest)
		return
	}

	event, err := repositoryControllers.GetEventByCodeReward(reward.Code, dbPool)

	if err != nil {
		http.Error(w, "Error Fetching Event", http.StatusInternalServerError)
		return
	}

	if event == nil {
		http.Error(w, "Invalid Code", http.StatusInternalServerError)
		return
	}

	users_event, err := repositoryControllers.GetAllUsersFromEventByEventId(*event.Id, dbPool)
	if err != nil {
		http.Error(w, "Error Fetching users_events", http.StatusInternalServerError)
		return
	}
	var foundUser bool = false
	for _, user_event := range users_event {

		if user_event.User_Id == reward.User_Id && user_event.Event_Id == *event.Id && user_event.IsRewarded {

			foundUser = true

			if user_event.ClaimedReward {
				http.Error(w, "Reward Already Claimed", http.StatusInternalServerError)
				return
			}
			user, err := repositoryControllersUsers.GetUsersById(user_event.User_Id, dbPool)
			if err != nil {
				http.Error(w, "Error Fetching User", http.StatusInternalServerError)
				return
			}

			rents, err := repositoryControllersUsers.GetRentByApartmentId(*user.Apartment_id, dbPool)
			if err != nil {
				http.Error(w, "Error Fetching Rents", http.StatusInternalServerError)
				return
			}

			var rent_to_discount models.Rent

			for _, rent := range rents { // Find last rent that was not paid to discount
				if rent.Status == "unpaid" {
					rent_to_discount = rent
					break
				}
			}
			reduction := rent_to_discount.Reduction + event.Percentage
			newFinal := rent_to_discount.Base_Amount * (1 - reduction)

			err = repositoryControllersUsers.UpdateRentReductionAndFinalAmount(*rent_to_discount.Id, reduction, newFinal, dbPool)
			if err != nil {
				http.Error(w, "Failed Appling Reduction", http.StatusInternalServerError)
				return
			}

			err = repositoryControllers.UpdateRewardedStatus(*event.Id, user_event.User_Id, dbPool)
			if err != nil {
				http.Error(w, "Error Updating Reward Status", http.StatusInternalServerError)
				return
			}
		}
	}

	if !foundUser {
		http.Error(w, "Error Applying Discount", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Reward Applied Sucessfully"})

}
