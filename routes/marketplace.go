package routes

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/portilho13/neighborconnect-backend/controllers"
)

func MarketplaceItemsApiRoute(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("POST /api/v1/marketplace", func(w http.ResponseWriter, r *http.Request) {
		controllers.RegisterClient(w, r, dbPool)
	})
}