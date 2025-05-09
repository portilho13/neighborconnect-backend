package routes

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/portilho13/neighborconnect-backend/controllers"
)

func GetAccountApi(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("GET /api/v1/account", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetAccount(w, r, dbPool)
	})
}

func GetAccountMovementsApi(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("GET /api/v1/account/movement", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetAccountMovements(w, r, dbPool)
	})
}
