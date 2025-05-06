package utils

import (
	"github.com/jackc/pgx/v5/pgxpool"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
)

func ValidateCreditCardDeposit() bool {
	// Implement Gatway ???
	return true
}

func ValidateWalletBalance(user_id int, amount float64, dbPool *pgxpool.Pool) bool {
	user, err := repositoryControllers.GetUsersById(user_id, dbPool)
	if err != nil {
		return false
	}

	account, err := repositoryControllers.GetAccountByUserId(user.Id, dbPool)
	if err != nil {
		return false
	}

	if account.Balance >= amount {
		newBalance := account.Balance - amount
		_ = repositoryControllers.UpdateAccountBalance(account.Id, newBalance, dbPool)
		return true
	}

	return false
}
