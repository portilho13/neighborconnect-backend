package routes

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/portilho13/neighborconnect-backend/controllers"
)

func CreateEventApiRoute(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("POST /api/v1/event", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateEvent(w, r, dbPool)
	})
}

func GetEventsApiRoute(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("GET /api/v1/event", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetEvents(w, r, dbPool)
	})
}