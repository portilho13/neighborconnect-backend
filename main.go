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
	"github.com/portilho13/neighborconnect-backend/utils"
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
	routes.RegisterWebSocketRoute(mux, dbPool)
	routes.CreateEventApiRoute(mux, dbPool)
	routes.CreateCategory(mux, dbPool)
	routes.GetCategory(mux, dbPool)
	routes.GetEventsApiRoute(mux, dbPool)
	routes.AddUserToEventsApi(mux, dbPool)
	routes.GetRentsApi(mux, dbPool)
	routes.GetDashBoardInfoApi(mux, dbPool)
	routes.LoginManagerApiRoute(mux, dbPool)
	routes.CreateApartmentApi(mux, dbPool)
	routes.RegisterManagerApi(mux, dbPool)
	routes.CreateDepositApi(mux, dbPool)
	routes.CreateWithdrawApi(mux, dbPool)
	routes.CreateBuyApiRoute(mux, dbPool)
	routes.PayTransactionApi(mux, dbPool)
	routes.ConcludeEventApiRoute(mux, dbPool)
	routes.RewardEventApiRoute(mux, dbPool)
	routes.PayRentApi(mux, dbPool)
	routes.GetTransactions(mux, dbPool)
	routes.GetUsersEventsApiRoute(mux, dbPool)
	routes.GetAccountApi(mux, dbPool)
	routes.GetAccountMovementsApi(mux, dbPool)
	//routes.LogoutClientApiRoute(mux)
	//routes.LogoutManagerApiRoute(mux)

	routes.ServerFilesApi(mux)

	nextMux := middleware.Logging(middleware.CORS(mux))

	return nextMux
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	err = utils.CreateUploadsFolder()
	if err != nil {
		log.Fatal(err)
	}

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

	go utils.AutomateRents(dbPool)

	go utils.AutomateListingClosing(dbPool)

	go utils.AutomateEventDeleting(dbPool)

	go utils.AutomateTransactionDeleting(dbPool)

	mux := InitializeRoutes(dbPool)
	go ws.Hub.Run()

	fmt.Println("Start listening on:", apiIP)
	if err := http.ListenAndServe(apiIP, mux); err != nil {
		log.Fatal(err)
	}

}
