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

func GetAccountMovements(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	query := r.URL.Query()
	idStr := query.Get("user_id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Id", http.StatusBadRequest)
		return
	}

	accout, err := repositoryControllers.GetAccountByUserId(id, dbPool)
	if err != nil {
		http.Error(w, "Error Fetching account", http.StatusInternalServerError)
		return
	}

	account_movements, err := repositoryControllers.GetAccountMovementsByAccountId(accout.Id, dbPool)
	if err != nil {
		http.Error(w, "Error Fetching Account Movements", http.StatusInternalServerError)
		return
	}

	var account_movements_json []controllers_models.AccountMovementJson

	for _, account_movement := range account_movements {
		account_movement_json := controllers_models.AccountMovementJson{
			Id:         account_movement.Id,
			Amount:     account_movement.Ammount,
			Created_At: account_movement.Created_at,
			Account_Id: account_movement.Account_id,
			Type:       account_movement.Type,
		}

		account_movements_json = append(account_movements_json, account_movement_json)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(account_movements_json); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
