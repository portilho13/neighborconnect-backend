package utils

import (
	"log"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/events"
	"github.com/robfig/cron/v3"
)

var eventCodeLenght int = 6 // Lengh of event code, 6 by default

var eventCodeArr string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomEventCode() string {

	max := len(eventCodeArr)

	eventCode := make([]string, eventCodeLenght)

	for i := range eventCodeLenght {
		n := rand.IntN(max)
		eventCode[i] = string(eventCodeArr[n])
	}
	return strings.Join(eventCode, "")

}


func deleteExpiredEvents(dbPool *pgxpool.Pool) error {
	events, err := repositoryControllers.GetAllFinishedEvents(dbPool)
	if err != nil {
		return err
	}

	for _, event := range events {
		if time.Now().UTC().After(*event.Expiration_Date) {
			err = repositoryControllers.DeleteEventById(*event.Id, dbPool)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func AutomateEventDeleting(dbPool *pgxpool.Pool) {
	c := cron.New(cron.WithSeconds())
	c.AddFunc("@every 1s", func() { // Maybe change this ???
		err := deleteExpiredEvents(dbPool)
		if err != nil {
			log.Fatal(err)
		}

	})

	c.Start()

	defer c.Stop()

	select {}
}