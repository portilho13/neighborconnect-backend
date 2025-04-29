package utils

import (
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
	"github.com/robfig/cron/v3"
)

func ConvertMonthToInt(month time.Month) int {
	switch month {
	case time.January:
		return 1
	case time.February:
		return 2
	case time.March:
		return 3
	case time.April:
		return 4
	case time.May:
		return 5
	case time.June:
		return 6
	case time.July:
		return 7
	case time.August:
		return 8
	case time.September:
		return 9
	case time.October:
		return 10
	case time.November:
		return 11
	case time.December:
		return 12
	}
	return -1
}

func AutomateRents(dbPool *pgxpool.Pool) {
	c := cron.New(cron.WithSeconds())
	c.AddFunc("0 0 1 * *", func() { // Every month at 00:00
		err := repositoryControllers.CreateRent(dbPool)
		if err != nil {
			log.Printf("Error creating rent %v", err)
		} else {
			log.Println("Rent Created Sucessfully")
		}
	})

	c.Start()

	defer c.Stop()

	select {}

}
