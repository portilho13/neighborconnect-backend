package routes

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/portilho13/neighborconnect-backend/controllers"
)

func GetRentsApi(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("GET /api/v1/rent", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetRents(w, r, dbPool)
	})
}

func PayRentApi(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("POST /api/v1/rent", func(w http.ResponseWriter, r *http.Request) {
		controllers.PayRent(w, r, dbPool)
	})
}
