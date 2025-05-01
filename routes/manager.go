package routes

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/portilho13/neighborconnect-backend/controllers"
)

func RegisterManagerApi(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("POST /api/v1/manager/register", func(w http.ResponseWriter, r *http.Request) {
		controllers.RegisterManager(w, r, dbPool)
	})
}

func LoginManagerApiRoute(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("POST /api/v1/manager/login", func(w http.ResponseWriter, r *http.Request) {
		controllers.LoginManager(w, r, dbPool)
	})
}
