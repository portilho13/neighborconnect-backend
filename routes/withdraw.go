package routes

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/portilho13/neighborconnect-backend/controllers"
)

func CreateWithdrawApi(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("POST /api/v1/withdraw", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateWithdraw(w, r, dbPool)
	})
}
