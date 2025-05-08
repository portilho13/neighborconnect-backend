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

func GetUsersEventsApiRoute(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("GET /api/v1/event/users", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetUserListFromEventId(w, r, dbPool)
	})
}

func AddUserToEventsApi(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("POST /api/v1/event/add", func(w http.ResponseWriter, r *http.Request) {
		controllers.AddUserToEvents(w, r, dbPool)
	})
}

func ConcludeEventApiRoute(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("POST /api/v1/event/conclude", func(w http.ResponseWriter, r *http.Request) {
		controllers.ConcludeEvent(w, r, dbPool)
	})
}

func RewardEventApiRoute(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("POST /api/v1/event/reward", func(w http.ResponseWriter, r *http.Request) {
		controllers.Reward(w, r, dbPool)
	})
}
