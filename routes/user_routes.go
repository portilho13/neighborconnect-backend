package routes

import (
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/portilho13/neighborconnect-backend/controllers"
)

func RegisterUserRoutes(router *mux.Router, dbPool *pgxpool.Pool) {
	// Atualiza a password através do email (PUT /users/password-by-email)
	router.HandleFunc("/users/password-by-email", controllers.UpdatePasswordByEmail(dbPool)).Methods("PUT")
}
