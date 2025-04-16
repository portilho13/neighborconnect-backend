package routes

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/portilho13/neighborconnect-backend/controllers"
)

func CreateListingApiRoute(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("POST /api/v1/listing/create", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateListing(w, r, dbPool)
	})
}

func GetListingByIdApiRoute(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("GET /api/v1/listing", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetListingById(w, r, dbPool)
	})
}

func GetAllListingApiRoute(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("GET /api/v1/listing/all", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetAllListings(w, r, dbPool)
	})
}
