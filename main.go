package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/portilho13/neighborconnect-backend/middleware"
	"github.com/portilho13/neighborconnect-backend/repository"
	"github.com/portilho13/neighborconnect-backend/routes"
	"github.com/portilho13/neighborconnect-backend/ws"
)

func InitializeRoutes(dbPool *pgxpool.Pool) http.Handler {
	mux := http.NewServeMux()

	routes.TestApiRoute(mux)
	routes.RegisterClientApiRoute(mux, dbPool)
	routes.LoginClientApiRoute(mux, dbPool)
	routes.CreateListingApiRoute(mux, dbPool)
	routes.GetListingByIdApiRoute(mux, dbPool)
	routes.GetAllListingApiRoute(mux, dbPool)
	routes.CreateBidApiRoute(mux, dbPool)
	routes.RegisterWebSocketRoute(mux)

	nextMux := middleware.Logging(middleware.CORS(mux))

	return nextMux
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Get database URL
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	apiIP := os.Getenv("API_IP")
	if apiIP == "" {
		log.Fatal("Api Ip not set")
	}

	// Initialize database
	dbPool, err := repository.InitDB(databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer repository.CloseDB(dbPool)

	mux := InitializeRoutes(dbPool)
	go ws.Hub.Run()

	fmt.Println("Start listening on:", apiIP)
	if err := http.ListenAndServe(apiIP, mux); err != nil {
		log.Fatal(err)
	}

}
