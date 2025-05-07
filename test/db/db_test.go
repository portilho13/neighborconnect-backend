package test

import (
	"context"
	"testing"

	"github.com/portilho13/neighborconnect-backend/repository"
	"github.com/stretchr/testify/assert"
)

// TestInitDB tests the database initialization
func TestInitDB(t *testing.T) {
	// Use a test database URL (update as needed)
	testDBURL := "postgres://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable"

	// Initialize the DB connection
	dbPool, err := repository.InitDB(testDBURL)

	// Ensure no error occurred
	assert.NoError(t, err, "InitDB should not return an error")
	assert.NotNil(t, dbPool, "Database pool should be initialized")

	// Verify connection by executing a simple query
	var testValue int
	err = dbPool.QueryRow(context.Background(), "SELECT 1").Scan(&testValue)

	assert.NoError(t, err, "Database should be accessible")
	assert.Equal(t, 1, testValue, "Expected to receive 1 from test query")

	// Cleanup: Close the database pool
	repository.CloseDB(dbPool)

}
