package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	controllers_models "github.com/portilho13/neighborconnect-backend/models"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/users"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetRents(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	apartment_id_str := r.URL.Query().Get("apartment_id")

	apartmend_id, err := strconv.Atoi(apartment_id_str)
	if err != nil {
		http.Error(w, "Invalid Apartment ID", http.StatusBadRequest)
		return
	}

	rents, err := repositoryControllers.GetRentByApartmentId(apartmend_id, dbPool)
	if err != nil {
		http.Error(w, "Error Fetching Rents", http.StatusInternalServerError)
		return
	}

	var rentsJson []controllers_models.Rent

	for _, rent := range rents {
		rentJson := controllers_models.Rent{
			Id:           rent.Id,
			Month:        rent.Month,
			Year:         rent.Year,
			Base_Amount:  rent.Base_Amount,
			Reduction:    rent.Reduction,
			Final_Amount: rent.Final_Amount,
			Apartment_Id: rent.Apartment_Id,
			Status:       rent.Status,
			Due_Day:      rent.Due_day,
		}

		rentsJson = append(rentsJson, rentJson)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(rentsJson); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
