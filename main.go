package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/portilho13/neighborconnect-backend/middleware"
	"github.com/portilho13/neighborconnect-backend/repository"
	"github.com/portilho13/neighborconnect-backend/repository/models"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

		// Get database URL
		databaseURL := os.Getenv("DATABASE_URL")
		if databaseURL == "" {
			log.Fatal("DATABASE_URL is not set")
		}
	
		// Initialize database
		dbPool, err := repository.InitDB(databaseURL)
		if err != nil {
			log.Fatal(err)
		}
		defer repository.CloseDB(dbPool)

		user := models.User{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "securepassword",
			Phone:    "+123456789",
			Apartment_id: 1,
		}
		
		err = repository.CreateUser(user, dbPool)
		if err != nil {
			log.Fatal(err)
		}

	mux := InitializeRoutes();
	fmt.Println("Start listening on:", IP)
	if err := http.ListenAndServe(IP, mux); err != nil {
		log.Fatal(err)
	}
	
}
