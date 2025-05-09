package routes

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/portilho13/neighborconnect-backend/controllers"
	"github.com/portilho13/neighborconnect-backend/middleware"
)

func RegisterManagerApi(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("POST /api/v1/manager/register", func(w http.ResponseWriter, r *http.Request) {
		controllers.RegisterManager(w, r, dbPool)
	})
}

func LoginManagerApiRoute(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.Handle("POST /api/v1/manager/login",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controllers.LoginManager(w, r, dbPool)
		}),
	)
}

func LogoutManagerApiRoute(mux *http.ServeMux) {
	mux.Handle("/api/v1/manager/logout", middleware.RequireAuthentication("manager")(http.HandlerFunc(controllers.LogoutHandlerUser)))

}
