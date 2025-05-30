package routes

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/portilho13/neighborconnect-backend/controllers"
	"github.com/portilho13/neighborconnect-backend/middleware"
)

func RegisterClientApiRoute(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("POST /api/v1/client/register", func(w http.ResponseWriter, r *http.Request) {
		controllers.RegisterClient(w, r, dbPool)
	})
}

func LoginClientApiRoute(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.Handle("POST /api/v1/client/login",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controllers.LoginClient(w, r, dbPool)
		}),
	)
}

func LogoutClientApiRoute(mux *http.ServeMux) {
	mux.Handle("/api/v1/client/logout", middleware.RequireAuthentication("user")(http.HandlerFunc(controllers.LogoutHandlerUser)))

}

func GetClientsApiRoute(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.Handle("GET /api/v1/client",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controllers.GetNeighborInfo(w, r, dbPool)
		}),
	)
}

func UploadProfilePictureApi(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.Handle("POST /api/v1/client/upload",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			controllers.UploadProfilePicture(w, r, dbPool)
		}),
	)
}
