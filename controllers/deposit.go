package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	controllers_models "github.com/portilho13/neighborconnect-backend/models"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
	"github.com/portilho13/neighborconnect-backend/utils"
)

func CreateDeposit(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	var depositJson controllers_models.Deposit

	err := json.NewDecoder(r.Body).Decode(&depositJson)

	if err != nil {
		http.Error(w, "Invalid JSON Data", http.StatusBadRequest)
		return
	}

	switch depositJson.Type {
	case "credit card":
		// Handle Gatway Logic
		if !utils.ValidateCreditCardDeposit() {
			http.Error(w, "Gatway error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Unsupported Gatway", http.StatusInternalServerError)
		return
	}

	account, err := repositoryControllers.GetAccountByUserId(depositJson.User_id, dbPool)
	if err != nil {
		http.Error(w, "Error Fetching User Account", http.StatusInternalServerError)
		return
	}

	account_movement := models.Account_Movement{
		Ammount:    depositJson.Amount,
		Created_at: time.Now().UTC(),
		Account_id: &account.Id,
		Type:       "deposit",
	}

	err = repositoryControllers.CreateAccountMovement(account_movement, dbPool)
	if err != nil {
		http.Error(w, "Error Creating Account Movement", http.StatusInternalServerError)
		return
	}

	newBalance := account.Balance + depositJson.Amount
	err = repositoryControllers.UpdateAccountBalance(account.Id, newBalance, dbPool)

	if err != nil {
		http.Error(w, "Error Updating Account Balance", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Deposit Sucessfully!"})
}
