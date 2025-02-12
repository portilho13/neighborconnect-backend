package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/portilho13/neighborconnect-backend/middleware"
	"github.com/portilho13/neighborconnect-backend/routes"
)

const IP string = "127.0.0.1:1234"

func InitializeRoutes() http.Handler {
	mux := http.NewServeMux()

	routes.TestApiRoute(mux)

	nextMux := middleware.Logging(middleware.CORS(mux))

	return nextMux
}
func main() {
	mux := InitializeRoutes();
	fmt.Println("Start listening on:", IP)
	if err := http.ListenAndServe(IP, mux); err != nil {
		log.Fatal(err)
	}
	
}