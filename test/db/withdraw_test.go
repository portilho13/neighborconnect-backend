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
	var userID int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.users (name, email, password, phone) 
		 VALUES ('Alice Doe', 'alice@example.com', 'securepassword', '987654321') 
		 RETURNING id`).Scan(&userID)
	assert.NoError(t, err, "User insertion should succeed")

	// Step 2: Define test withdraw data
	withdraw := models.Withdraw{
		Ammount:    5000,
		Created_at: time.Now(),
		User_Id:    userID, // Ensure valid foreign key reference
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
}
