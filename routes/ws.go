package routes

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/portilho13/neighborconnect-backend/ws"
)

func RegisterWebSocketRoute(mux *http.ServeMux, dbPool *pgxpool.Pool) {
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(w, r, dbPool)
	})
}
