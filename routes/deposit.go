package routes

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/portilho13/neighborconnect-backend/controllers"
)

func CreateDepositApi(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("POST /api/v1/deposit", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateDeposit(w, r, dbPool)
	})
}