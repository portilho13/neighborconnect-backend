package test

import (
	"context"
	"testing"

	repository "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
	"github.com/stretchr/testify/assert"
)


func TestCreateUser(t *testing.T) {
	// Connect to the test database
	dbPool, err := GetTestDBConnection()
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}
	defer dbPool.Close()

	// Ensure the database is clean before starting
	CleanDatabase(dbPool, "users.users")

	// Define test user data
	user := models.User{
		Name:     "Alice Doe",
		Email:    "alice@example.com",
		Password: "securepassword",
		Phone:    "987654321",
	}

	// Call the function to insert the user
	err = repository.CreateUser(user, dbPool)
	assert.NoError(t, err, "CreateUser should not return an error")

	// Verify the user was inserted
	var count int
	err = dbPool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users.users").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count, "Expected one user in the database")

	// Cleanup after the test
	CleanDatabase(dbPool, "users.users")
}

func TestGetUserByEmail(t *testing.T) {
	// Connect to the test database
	dbPool, err := GetTestDBConnection()
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}
	defer dbPool.Close()

	// Ensure the database is clean before starting
	CleanDatabase(dbPool, "users.users")

	// Step 1: Insert a test user
	user := models.User{
		Name:     "Alice Doe",
		Email:    "alice@example.com",
		Password: "securepassword",
		Phone:    "911111111",
	}

	err = repository.CreateUser(user, dbPool)
	assert.NoError(t, err, "CreateUser should not return an error")

	// Step 2: Retrieve the user by email
	retrievedUser, err := repository.GetUserByEmail("alice@example.com", dbPool)
	assert.NoError(t, err, "GetUserByEmail should not return an error")
	assert.Equal(t, user.Email, retrievedUser.Email, "Retrieved email should match inserted user")
	assert.Equal(t, user.Name, retrievedUser.Name, "Retrieved name should match inserted user")
	assert.Equal(t, user.Phone, retrievedUser.Phone, "Retrieved phone should match inserted user")

	// Cleanup after the test
	CleanDatabase(dbPool, "users.users")
}