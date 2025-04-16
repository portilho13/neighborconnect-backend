package routes

import (
	"net/http"

	"github.com/portilho13/neighborconnect-backend/ws"
)

func RegisterWebSocketRoute(mux *http.ServeMux) {
	mux.HandleFunc("/ws", ws.ServeWs)
}
