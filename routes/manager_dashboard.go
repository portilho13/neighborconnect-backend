package routes

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/portilho13/neighborconnect-backend/controllers"
)

func GetDashBoardInfoApi(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("GET /api/v1/manager/dashboard/info", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetDashBoardInfo(w, r, dbPool)
	})
}
