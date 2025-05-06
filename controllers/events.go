package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	controllers_models "github.com/portilho13/neighborconnect-backend/models"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/events"
	models "github.com/portilho13/neighborconnect-backend/repository/models/events"
	"github.com/portilho13/neighborconnect-backend/utils"
)

func CreateEvent(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {

	err := r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	eventJson := r.FormValue("event")
	if eventJson == "" {
		http.Error(w, "Missing event data", http.StatusBadRequest)
		return
	}

	var eventData controllers_models.EventCreation
	if err := json.Unmarshal([]byte(eventJson), &eventData); err != nil {
		http.Error(w, "Invalid event data format", http.StatusBadRequest)
		return
	}

	var api_path string

	files := r.MultipartForm.File["images"]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Failed to create event", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		filename := fileHeader.Filename
		ext := filepath.Ext(filename)
		newFilename := uuid.New().String() + ext
		savePath := fmt.Sprintf("./uploads/events/%s", newFilename)

		dst, err := os.Create(savePath)
		if err != nil {
			http.Error(w, "Failed to create event", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Failed to create listing", http.StatusInternalServerError)
			return
		}

		api_url := utils.GetApiUrl()

		api_path = fmt.Sprintf("http://%s/api/v1/uploads/events/%s", api_url, newFilename)

	}
	event := models.Community_Event{
		Name:              eventData.Name,
		Percentage:        eventData.Percentage,
		Code:              utils.GenerateRandomEventCode(),
		Capacity:          eventData.Capacity,
		Date_Time:         eventData.Date_time,
		Manager_Id:        &eventData.Manager_Id,
		Event_Image:       api_path,
		Duration:          eventData.Duration,
		Local:             eventData.Local,
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
				Event_Image:       event.Event_Image,
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
				Event_Image:       event.Event_Image,
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

	user_event := models.User_Event{
		User_Id:    joinEventJson.User_Id,
		Event_Id:   joinEventJson.Community_Event_Id,
		IsRewarded: false,
	}

	err = repositoryControllers.AddUserToCommunityEvent(user_event, dbPool)
	if err != nil {
		http.Error(w, "Error Adding User to Event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User Added Sucessfully"})
}

func ConcludeEvent(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	var concludeEvent controllers_models.ConcludeEventJson

	err := json.NewDecoder(r.Body).Decode(&concludeEvent)
	if err != nil {
		http.Error(w, "Invalid JSON Data", http.StatusBadRequest)
		return
	}

	for _, user_id := range concludeEvent.Awarded_Users_Ids {
		err = repositoryControllers.UpdateRewardedStatus(concludeEvent.Event_Id, user_id, dbPool)
		if err != nil {
			http.Error(w, "Error Updating Is Rewarded Status", http.StatusBadRequest)
			return
		}

	}

}
