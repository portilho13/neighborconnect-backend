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

func CreateWithdraw(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	var withdrawJson controllers_models.Withdraw

	err := json.NewDecoder(r.Body).Decode(&withdrawJson)

	if err != nil {
		http.Error(w, "Invalid JSON Data", http.StatusBadRequest)
		return
	}

	switch withdrawJson.Type {
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

	account, err := repositoryControllers.GetAccountByUserId(withdrawJson.User_id, dbPool)
	if err != nil {
		http.Error(w, "Error Fetching User Account", http.StatusInternalServerError)
		return
	}

	if account.Balance-withdrawJson.Amount < 0 {
		http.Error(w, "User balance too low", http.StatusInternalServerError)
		return
	}

	account_movement := models.Account_Movement{
		Ammount:    0 - withdrawJson.Amount,
		Created_at: time.Now(),
		Account_id: &account.Id,
		Type:       "withdraw",
	}

	err = repositoryControllers.CreateAccountMovement(account_movement, dbPool)
	if err != nil {
		http.Error(w, "Error Creating Account Movement", http.StatusInternalServerError)
		return
	}

	newBalance := account.Balance - withdrawJson.Amount

	err = repositoryControllers.UpdateAccountBalance(account.Id, newBalance, dbPool)

	if err != nil {
		http.Error(w, "Error Updating Account Balance", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Withdraw Created !"})
}
