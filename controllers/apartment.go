package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	controllers_models "github.com/portilho13/neighborconnect-backend/models"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
)

func CreateApartment(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	var apartmentJson controllers_models.ApartmentCreation

	err := json.NewDecoder(r.Body).Decode(&apartmentJson)

	if err != nil {
		http.Error(w, "Invalid JSON Data", http.StatusBadRequest)
		return
	}

	apartment := models.Apartment{
		N_bedrooms: apartmentJson.N_Bedrooms,
		Floor:      apartmentJson.Floor,
		Rent:       apartmentJson.Rent,
		Manager_id: apartmentJson.Manager_Id,
		Status:     "unoccupied",
	}

	err = repositoryControllers.CreateApartment(apartment, dbPool)

	if err != nil {
		http.Error(w, "Error Creating Apartment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Apartment Created !"})

}
