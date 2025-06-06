package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	controllers "github.com/portilho13/neighborconnect-backend/controllers"
	controllers_models "github.com/portilho13/neighborconnect-backend/models"
	repository "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
	"github.com/portilho13/neighborconnect-backend/utils"
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

func TestLoginHandler(t *testing.T) {
	dbPool, err := GetTestDBConnection()
	assert.NoError(t, err)
	defer dbPool.Close()

	CleanDatabase(dbPool, "users.users")

	// Ensure the database is clean before starting
	CleanDatabase(dbPool, "users.users")

	encodedPassword, err := utils.GenerateFromPassword("securepassword", utils.DefaultArgon2Params)
	if err != nil {
		t.Fatalf("Error hashing: %v", err)
	}
	var managerId int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.manager (name, email, password, phone) VALUES ('a', 'a', 'a', 'a')
		 RETURNING id`).Scan(&managerId)
	assert.NoError(t, err, "User insertion should succeed")


	var apartmentID int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.apartment (n_bedrooms, floor, rent, manager_id, status) VALUES (1, 1, 1, $1, 'unoccupied')
		 RETURNING id`, managerId).Scan(&apartmentID)
	assert.NoError(t, err, "User insertion should succeed")
	// Step 1: Insert a test user
	user := models.User{
		Name:     "Alice Doe",
		Email:    "alice@example.com",
		Password: encodedPassword,
		Apartment_id: &apartmentID,
		Phone:    "911111111",
	}
	err = repository.CreateUser(user, dbPool)
	assert.NoError(t, err, "CreateUser should not return an error")

	// Test valid login
	credentials := controllers_models.Credentials{Email: "alice@example.com", Password: "securepassword"}
	body, _ := json.Marshal(credentials)
	req, _ := http.NewRequest("POST", "/api/v1/client/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	controllers.LoginClient(w, req, dbPool)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEqual(t, w.Body.String(), "Invalid credentials")

	// Test invalid login
	wrongCredentials := controllers_models.Credentials{Email: "test@example.com", Password: "wrongpassword"}
	body, _ = json.Marshal(wrongCredentials)
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	w = httptest.NewRecorder()

	controllers.LoginClient(w, req, dbPool)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid credentials")

	CleanDatabase(dbPool, "users.users")
	CleanDatabase(dbPool, "users.apartment")
	CleanDatabase(dbPool, "users.manager")


}
