package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

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
		Name:              eventJson.Name,
		Percentage:        eventJson.Percentage,
		Code:              utils.GenerateRandomEventCode(),
		Capacity:          eventJson.Capacity,
		Date_Time:         eventJson.Date_time,
		Manager_Id:        managerId,
		Event_Image:       &eventJson.Event_Image,
		Duration:          eventJson.Duration,
		Local:             eventJson.Local,
		Current_Ocupation: 0,
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

func GetEvents(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	var eventJsonList []controllers_models.EventInfo

	userID_str := r.URL.Query().Get("user_id")

	if userID_str != "" {
		userId, err := strconv.Atoi(userID_str)
		if err != nil {
			http.Error(w, "Error Converting User Id", http.StatusInternalServerError)
			return
		}

		events, err := repositoryControllers.GetEventsByUserId(userId, dbPool)
		if err != nil {
			http.Error(w, "Error Fetching Events", http.StatusBadRequest)
			return
		}

		for _, event := range events {
			eventJsonList = append(eventJsonList, controllers_models.EventInfo{
				Id:                *event.Id,
				Name:              event.Name,
				Percentage:        event.Percentage,
				Capacity:          event.Capacity,
				Date_time:         event.Date_Time,
				Manager_Id:        *event.Manager_Id,
				Event_Image:       *event.Event_Image,
				Duration:          event.Duration,
				Local:             event.Local,
				Current_Ocupation: event.Current_Ocupation,
			})

		}
	} else {
		events, err := repositoryControllers.GetAllEvents(dbPool)
		if err != nil {
			http.Error(w, "Error Fetching Events", http.StatusBadRequest)
			return
		}

		for _, event := range events {
			eventJsonList = append(eventJsonList, controllers_models.EventInfo{
				Id:                *event.Id,
				Name:              event.Name,
				Percentage:        event.Percentage,
				Capacity:          event.Capacity,
				Date_time:         event.Date_Time,
				Manager_Id:        *event.Manager_Id,
				Event_Image:       *event.Event_Image,
				Duration:          event.Duration,
				Local:             event.Local,
				Current_Ocupation: event.Current_Ocupation,
			})

		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(eventJsonList); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func AddUserToEvents(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	var joinEventJson controllers_models.JoinEvent

	err := json.NewDecoder(r.Body).Decode(&joinEventJson)
	if err != nil {
		http.Error(w, "Invalid JSON Data", http.StatusBadRequest)
		return
	}

	err = repositoryControllers.AddUserToCommunityEvent(joinEventJson.User_Id, joinEventJson.Community_Event_Id, dbPool)
	if err != nil {
		http.Error(w, "Error Adding User to Event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User Added Sucessfully"})
}
