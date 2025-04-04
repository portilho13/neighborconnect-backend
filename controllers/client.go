package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	controllers_models "github.com/portilho13/neighborconnect-backend/models"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
	"github.com/portilho13/neighborconnect-backend/utils"
)

func RegisterClient(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
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
	if client.ApartmentID == 0 { // If apartment id is 0 it means user as not an a
		// partment id set yet so apartment id needs to be converted to pointer to not blow the db
		apartmentID = nil
	} else {
		apartmentID = &client.ApartmentID
	}

	fmt.Println(client.Phone)

	dbClient := models.User{
		Name:         client.Name,
		Email:        client.Email,
		Password:     encodedPassword,
		Phone:        client.Phone,
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

	userAccount := models.Account{
		Account_number: utils.GenerateRandomHash(),
		Balance:        0,
		Currency:       "EUR",
		Users_id:       &dbClient.Id,
	}

	err = repositoryControllers.CreateAccount(userAccount, dbPool)
	if err != nil {
		log.Printf("-%v", err)
		http.Error(w, "Error creating account", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Client Registed !"})
}

func LoginClient(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {

	// var creds controllers_models.UserJson
	var creds controllers_models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Fetch user from the database
	user, err := repositoryControllers.GetUserByEmail(creds.Email, dbPool)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Verify password
	_, err = utils.ComparePasswordAndHash(creds.Password, user.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create a session
	session, _ := utils.Store.Get(r, "session-name")
	session.Values["user_id"] = user.Id
	session.Values["email"] = user.Email
	session.Save(r, w)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}
