package routes

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/portilho13/neighborconnect-backend/controllers"
)

func CreateCategory(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("POST /api/v1/category", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateCategory(w, r, dbPool)
	})
}

func GetCategory(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("GET /api/v1/category", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetCategories(w, r, dbPool)
	})
}