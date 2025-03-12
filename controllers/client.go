package controllers

import (
	"encoding/json"
	"fmt"
	"log"
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

	encodedPassword, err := utils.GenerateFromPassword(client.Password, utils.DefaultArgon2Params)
	if err != nil {
		log.Printf("Error: %v", err)
		http.Error(w, "Error creating user", http.StatusBadGateway)
		return
	}

	var apartmentID *int
	if client.ApartmentID == 0 {
		apartmentID = nil
	} else {
		apartmentID = &client.ApartmentID
	}

	fmt.Println(client.Phone)

	dbClient := models.User {
		Name: client.Name,
		Email: client.Email,
		Password: encodedPassword,
		Phone: client.Phone,
		Apartment_id: apartmentID,
	}

	err = repositoryControllers.CreateUser(dbClient, dbPool)
	if err != nil {
		log.Printf("-%v", err)
		http.Error(w, "Error creating user", http.StatusBadGateway)
		return
	}

	var clientEmail string = dbClient.Email

	dbClient, err = repositoryControllers.GetUserByEmail(clientEmail, dbPool)
	if err != nil {
		log.Printf("-%v", err)
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
		log.Printf("-%v", err)
		http.Error(w, "Error creating account", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}