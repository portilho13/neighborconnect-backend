package test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	repository "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
	"github.com/stretchr/testify/assert"
)

// Cleanup function to delete all data after each test
func cleanDatabase(dbPool *pgxpool.Pool) {
	_, _ = dbPool.Exec(context.Background(), "TRUNCATE users.apartment, users.manager RESTART IDENTITY CASCADE")
}

func TestCreateApartmentIntegration(t *testing.T) {
	// Connect to the test database
	dbPool, err := GetTestDBConnection()
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}
	defer dbPool.Close()

	// Ensure the database is clean before starting
	cleanDatabase(dbPool)

	// Step 1: Insert a manager first
	var managerID int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.manager (name, email, password, phone) 
		 VALUES ('John Doe', 'john@example.com', 'securepass', '123456789') 
		 RETURNING id`).Scan(&managerID)
	assert.NoError(t, err, "Manager insertion should succeed")

	// Step 2: Insert an apartment with the valid manager_id
	apartment := models.Apartment{
		N_bedrooms: 3,
		Floor:      2,
		Rent:       1500,
		Manager_id: managerID,  // Use the valid manager ID
	}

	err = repository.CreateApartment(apartment, dbPool)
	assert.NoError(t, err, "CreateApartment should not return an error")

	// Verify the apartment was inserted
	var count int
	err = dbPool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users.apartment").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count, "Expected one apartment in the database")
	cleanDatabase(dbPool)
}
