package controllers

import (
	"encoding/json"
	"net/http"

	users "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
	"github.com/portilho13/neighborconnect-backend/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UpdatePasswordRequest struct {
	Email       string `json:"email"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func UpdatePasswordByEmail(dbPool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdatePasswordRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		//Obter dados do utilizador através do email
		user, err := users.GetUserByEmail(req.Email, dbPool)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Verificar se a password antiga está correta
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
			http.Error(w, "Incorrect old password", http.StatusUnauthorized)
			return
		}

		// Hashear a nova password
		hashedPassword, err := utils.HashPassword(req.NewPassword)
		if err != nil {
			http.Error(w, "Could not hash password", http.StatusInternalServerError)
			return
		}

		//Atualizar password na base de dados
		err = users.UpdateUserPasswordByEmail(req.Email, hashedPassword, dbPool)
		if err != nil {
			http.Error(w, "Could not update password", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Password updated successfully"))
	}
}
