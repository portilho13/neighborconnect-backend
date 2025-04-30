package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	controllers_models "github.com/portilho13/neighborconnect-backend/models"
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

	for _, apartment := range apartments {
		apartmentJson := controllers_models.Apartment{
			Id:         *apartment.Id,
			N_Bedrooms: apartment.N_bedrooms,
			Floor:      apartment.Floor,
			Rent:       int(apartment.Rent),
			Manager_Id: apartment.Manager_id,
		}
		apartmentsJson = append(apartmentsJson, apartmentJson)

	}

	dashboardInfo := controllers_models.ManagerDashboardInfo{
		Apartments: apartmentsJson,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dashboardInfo); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}
