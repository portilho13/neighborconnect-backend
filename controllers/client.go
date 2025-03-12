package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	controllers_models "github.com/portilho13/neighborconnect-backend/controllers/models"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
	"github.com/portilho13/neighborconnect-backend/utils"
)

func RegisterClient(w http.ResponseWriter, r* http.Request, dbPool *pgxpool.Pool) {
	var client controllers_models.UserJson
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		http.Error(w, "Invalid JSON Data", http.StatusBadRequest)
		return
	}

	dbClient := models.User {
		Name: client.Name,
		Email: client.Email,
		Password: client.Password,
		Phone: client.Phone,
		Apartment_id: &client.ApartmentID,
	}

	err = repositoryControllers.CreateUser(dbClient, dbPool)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusBadGateway)
		return
	}

	var clientEmail string = dbClient.Email

	dbClient, err = repositoryControllers.GetUserByEmail(clientEmail, dbPool)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusBadGateway)
		return
	}

	userAccount := models.Account {
		Account_number: utils.GenerateRandomHash(),
		Balance: 0,
		Currency: "EUR",
		Users_id: dbClient.Id,
	}

	err = repositoryControllers.CreateAccount(userAccount, dbPool)
	if err != nil {
		http.Error(w, "Error creating account", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}