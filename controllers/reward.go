package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	controllers_models "github.com/portilho13/neighborconnect-backend/models"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/events"
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

	for _, user_event := range users_event {
		if user_event.User_Id == reward.User_Id && user_event.Event_Id == *event.Id && user_event.IsRewarded {
			// Apply Discount Latter
		} else {
			http.Error(w, "Error Applying Discount", http.StatusInternalServerError)
			return
		}
	}

}
