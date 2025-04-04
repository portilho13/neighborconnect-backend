package routes

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/portilho13/neighborconnect-backend/controllers"
)

func LoginManagerApiRoute(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("/api/v1/manager/login", func(w http.ResponseWriter, r *http.Request) {
		controllers.LoginManager(w, r, dbPool)
	})
}
