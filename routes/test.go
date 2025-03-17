package routes

import (
	"net/http"

	"github.com/portilho13/neighborconnect-backend/controllers"
)

func TestApiRoute(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/test", controllers.TestAPI)
}