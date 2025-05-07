package utils

import (
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/marketplace"
	"github.com/robfig/cron/v3"
)

func deleteExpiredTransactions(dbPool *pgxpool.Pool) error {
	transactions, err := repositoryControllers.GetAllTransactions(dbPool)
	if err != nil {
		return err
	}

	for _, transaction := range transactions {
		if time.Now().UTC().After(transaction.Payment_Due_time) {
			err = repositoryControllers.DeleteTransactionById(*transaction.Id, dbPool)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func AutomateTransactionDeleting(dbPool *pgxpool.Pool) {
	c := cron.New(cron.WithSeconds())
	c.AddFunc("@every 1s", func() { // Maybe change this ???
		err := deleteExpiredTransactions(dbPool)
		if err != nil {
			log.Fatal(err)
		}

	})

	c.Start()

	defer c.Stop()

	select {}
}