package routes

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/portilho13/neighborconnect-backend/controllers"
)

func LoginClientApiRoute(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("POST /api/v1/client/login", func(w http.ResponseWriter, r *http.Request) {
		controllers.LoginClient(w, r, dbPool)
	})
}