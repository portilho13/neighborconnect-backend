package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	controllers_models "github.com/portilho13/neighborconnect-backend/models"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
	"github.com/portilho13/neighborconnect-backend/utils"
)

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

	session, _ := utils.Store.Get(r, "session-name")
	session.Values["manager_id"] = manager.Id
	session.Values["email"] = manager.Email
	session.Save(r, w)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Manager login successful"})
}
