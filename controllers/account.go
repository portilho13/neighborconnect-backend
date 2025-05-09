package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	controllers_models "github.com/portilho13/neighborconnect-backend/models"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
)

func GetAccount(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	query := r.URL.Query()
	idStr := query.Get("user_id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Id", http.StatusBadRequest)
		return
	}

	var accountJson controllers_models.AccountJson

	accout, err := repositoryControllers.GetAccountByUserId(id, dbPool)
	if err != nil {
		http.Error(w, "Error Fetching account", http.StatusInternalServerError)
		return
	}

	acountJson := controllers_models.AccountJson{
		Id:            accout.Id,
		AccountNumber: accountJson.AccountNumber,
		Balance:       accout.Balance,
		Currency:      accout.Currency,
		UsersID:       accout.Users_id,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(acountJson); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
