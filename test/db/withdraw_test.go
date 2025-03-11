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
	CleanDatabase(dbPool, "users.withdraw")
	CleanDatabase(dbPool, "users.users")

	// Step 1: Insert a test user (since `withdraw.users_id` is a foreign key)
	user := models.User{
		Name: "t",
		Email: "t",
		Password: "t",
		Phone: "123",
	}

	err = repository.CreateUser(user, dbPool)
	assert.NoError(t, err, "Create user should not give an error");

	accout := models.Account {
		Account_number: "asd",
		Balance: 1000,
		Currency: "CHF",
		Users_id: 1,
	}

	err = repository.CreateAccount(accout, dbPool)
	assert.NoError(t, err, "Create account should not give an error");

	// Step 2: Define test withdraw data
	withdraw := models.Withdraw{
		Ammount:    5000,
		Created_at: time.Now(),
		Account_id:    1, // Ensure valid foreign key reference
	}

	// Step 3: Call the function to insert the withdrawal
	err = repository.CreateWithdraw(withdraw, dbPool)
	assert.NoError(t, err, "CreateWithdraw should not return an error")

	// Step 4: Verify the withdrawal was inserted
	var count int
	err = dbPool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users.withdraw").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count, "Expected one withdraw record in the database")

	// Cleanup after the test
	CleanDatabase(dbPool, "users.withdraw")
	CleanDatabase(dbPool, "users.users")
	CleanDatabase(dbPool, "users.account")
}
