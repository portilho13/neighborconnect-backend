package test

import (
	"context"
	"testing"
	"time"

	repository "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
	"github.com/stretchr/testify/assert"
)


func TestCreateWithdraw(t *testing.T) {
	// Connect to the test database
	dbPool, err := GetTestDBConnection()
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}
	defer dbPool.Close()

	// Ensure the database is clean before starting
	CleanDatabase(dbPool, "users.account")
	CleanDatabase(dbPool, "users.users")
	CleanDatabase(dbPool, "users.withdraw")

	// Step 1: Insert a test user
	user := models.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "securepassword",
		Phone:    "123456789",
	}

	// Create the user and capture the ID
	err = repository.CreateUser(user, dbPool)
	assert.NoError(t, err, "Create user should not give an error")

	user, err = repository.GetUserByEmail("test@example.com", dbPool)
	assert.NoError(t, err, "Get user should not give an error")

	// // Step 2: Insert a test account for the user
	account := models.Account{
		Account_number: "123ABC",
		Balance:       1000,
		Currency:      "CHF",
		Users_id:       &user.Id,
	}

	err = repository.CreateAccount(account, dbPool)
	assert.NoError(t, err, "Create account should not give an error")

	account, err = repository.GetAccountByUserId(user.Id, dbPool)
	assert.NoError(t, err, "Get account should not give an error")


	// // Step 3: Insert a withdrawal (Account Movement)
	accountMovement := models.Account_Movement{
		Ammount:    5000,
		Created_at: time.Now(),
		Account_id: &account.Id,
		Type:      "withdrawal",
	}

	err = repository.CreateAccountMovement(accountMovement, dbPool)
	assert.NoError(t, err, "CreateWithdraw should not return an error")

	// // Step 4: Verify the withdrawal was inserted
	var count int
	err = dbPool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users.account_movement").Scan(&count)
	if err != nil {
		t.Fatalf("Query error: %v", err) // Fail if query fails
	}
	assert.Equal(t, 1, count, "Expected one withdraw record in the database")

	//Cleanup after the test
	CleanDatabase(dbPool, "users.withdraw")
	CleanDatabase(dbPool, "users.account")
	CleanDatabase(dbPool, "users.users")
}

