package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	controllers_models "github.com/portilho13/neighborconnect-backend/models"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
	"github.com/portilho13/neighborconnect-backend/utils"
)

func RegisterManager(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	var client controllers_models.ManagerCreationJson
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		http.Error(w, "Invalid JSON Data", http.StatusBadRequest)
		return
	}
	// Verificar duplicação
	existingManager, err := repositoryControllers.GetManagerByEmail(client.Email, dbPool)
	if err == nil && existingManager.Email != "" {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}
	encodedPassword, err := utils.GenerateFromPassword(client.Password, utils.DefaultArgon2Params)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusBadGateway)
		return
	}

	dbManager := models.Manager{
		Name:     client.Name,
		Email:    client.Email,
		Password: encodedPassword,
		Phone:    client.Phone,
	}

	err = repositoryControllers.CreateManager(dbManager, dbPool)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Manager Registed !"})
}

func LoginManager(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	var creds controllers_models.Credentials

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	manager, err := repositoryControllers.GetManagerByEmail(creds.Email, dbPool)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	_, err = utils.ComparePasswordAndHash(creds.Password, manager.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	session, _ := utils.Store.Get(r, "manager-session")
	session.Values["manager_id"] = manager.Id
	session.Values["email"] = manager.Email
	session.Values["role"] = "manager"
	session.Save(r, w)

	managerJson := controllers_models.ManagerInfoJson{
		Id:    manager.Id,
		Name:  manager.Name,
		Email: manager.Email,
		Phone: manager.Phone,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(managerJson); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func LogoutHandlerManager(w http.ResponseWriter, r *http.Request) {
	session, _ := utils.Store.Get(r, "session")

	delete(session.Values, "user_id")
	delete(session.Values, "email")

	session.Options.MaxAge = -1

	session.Save(r, w)

}
