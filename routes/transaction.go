package routes

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/portilho13/neighborconnect-backend/controllers"
)

func PayTransactionApi(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("POST /api/v1/transaction", func(w http.ResponseWriter, r *http.Request) {
		controllers.PayTransaction(w, r, dbPool)
	})
}

func GetTransactions(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("GET /api/v1/transaction", func(w http.ResponseWriter, r *http.Request) {
		controllers.GetAllTransactions(w, r, dbPool)
	})
}
