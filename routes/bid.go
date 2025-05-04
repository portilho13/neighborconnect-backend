package routes

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/portilho13/neighborconnect-backend/controllers"
)

func CreateBidApiRoute(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("POST /api/v1/bid/", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateBid(w, r, dbPool)
	})
}
