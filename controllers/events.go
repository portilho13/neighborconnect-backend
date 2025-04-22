package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	controllers_models "github.com/portilho13/neighborconnect-backend/models"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/events"
	models "github.com/portilho13/neighborconnect-backend/repository/models/events"
	"github.com/portilho13/neighborconnect-backend/utils"
)

func CreateEvent(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	var eventJson controllers_models.EventCreation
	err := json.NewDecoder(r.Body).Decode(&eventJson)

	if err != nil {
		http.Error(w, "Invalid JSON Data", http.StatusBadRequest)
		return
	}

	var managerId *int
	if eventJson.Manager_Id == 0 {
		managerId = nil
	} else {
		managerId = &eventJson.Manager_Id
	}

	event := models.Community_Event{
		Name:        eventJson.Name,
		Percentage:  eventJson.Percentage,
		Code:        utils.GenerateRandomEventCode(),
		Capacity:    eventJson.Capacity,
		Date_Time:   eventJson.Date_time,
		Manager_Id:  managerId,
		Event_Image: &eventJson.Event_Image,
		Duration:    eventJson.Duration,
	}

	err = repositoryControllers.CreateCommunityEvent(event, dbPool)
	if err != nil {
		http.Error(w, "Error Creating Event", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Event Created Sucessfully"})
}
