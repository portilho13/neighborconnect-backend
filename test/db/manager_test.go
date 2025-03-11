package test

import (
	"context"
	"testing"

	repository "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
	"github.com/stretchr/testify/assert"
)


func TestCreateManager(t *testing.T) {
	// Connect to the test database
	dbPool, err := GetTestDBConnection()
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}
	defer dbPool.Close()

	// Ensure the database is clean before starting
	CleanDatabase(dbPool, "users.manager")

	// Define test manager data
	manager := models.Manager{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "securepass",
		Phone:    "123456789",
	}

	// Call the function to insert the manager
	err = repository.CreateManager(manager, dbPool)
	assert.NoError(t, err, "CreateManager should not return an error")

	// Verify the manager was inserted
	var count int
	err = dbPool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users.manager").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count, "Expected one manager in the database")

	// Cleanup after the test
	CleanDatabase(dbPool, "users.manager")
}

func TestGetManagerByEmail(t *testing.T) {
	// Connect to the test database
	dbPool, err := GetTestDBConnection()
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}
	defer dbPool.Close()

	// Ensure the database is clean before starting
	CleanDatabase(dbPool, "users.manager")

	// Insert a test manager
	manager := models.Manager{
		Name:     "Alice Doe",
		Email:    "alice@example.com",
		Password: "securepassword",
		Phone:    "911111111",
	}

	err = repository.CreateManager(manager, dbPool)
	assert.NoError(t, err, "CreateManager should not return an error")

	retrievedManager, err := repository.GetManagerByEmail("alice@example.com", dbPool)
	assert.NoError(t, err, "GetManagerByEmail should not return an error")
	assert.Equal(t, manager.Email, retrievedManager.Email, "Retrieved email should match inserted manager")
	assert.Equal(t, manager.Name, retrievedManager.Name, "Retrieved name should match inserted manager")
	assert.Equal(t, manager.Phone, retrievedManager.Phone, "Retrieved phone should match inserted manager")

	// Cleanup
	CleanDatabase(dbPool, "users.manager")
}


