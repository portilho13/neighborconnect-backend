package test

import (
	"context"
	"testing"

	repository "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	// Connect to the test database
	dbPool, err := GetTestDBConnection()
	require.NoError(t, err, "Failed to connect to test DB")
	defer dbPool.Close()

	// Ensure the database is clean before starting
	CleanDatabase(dbPool, "users.users,users.account")

	// Create test user
	var userId int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.users (name, email, password, phone) 
         VALUES ('Account User', 'account@example.com', 'securepass', '123456789') 
         RETURNING id`).Scan(&userId)
	require.NoError(t, err, "User insertion should succeed")

	// Prepare test account
	account := models.Account{
		Account_number: "ACC123456789",
		Balance:        1000.50,
		Currency:       "USD",
		Users_id:       &userId,
	}

	// Test the function
	err = repository.CreateAccount(account, dbPool)
	require.NoError(t, err, "CreateAccount should not return error")

	// Verify the account was created
	var dbAccount struct {
		Number   string
		Balance  float64
		Currency string
		UserID   int
	}
	err = dbPool.QueryRow(context.Background(),
		`SELECT account_number, balance, currency, users_id 
             FROM users.account 
             WHERE account_number = $1`, account.Account_number).Scan(
		&dbAccount.Number,
		&dbAccount.Balance,
		&dbAccount.Currency,
		&dbAccount.UserID)
	require.NoError(t, err, "Should be able to query created account")

	assert.Equal(t, account.Account_number, dbAccount.Number, "Account number mismatch")
	assert.Equal(t, account.Balance, dbAccount.Balance, "Account balance mismatch")
	assert.Equal(t, account.Currency, dbAccount.Currency, "Account currency mismatch")
	assert.Equal(t, *account.Users_id, dbAccount.UserID, "User ID mismatch")
	CleanDatabase(dbPool, "users.users,users.account")
}
func TestGetAccountByUserId(t *testing.T) {
	// Connect to the test database
	dbPool, err := GetTestDBConnection()
	require.NoError(t, err, "Failed to connect to test DB")
	defer dbPool.Close()

	// Ensure the database is clean before starting
	CleanDatabase(dbPool, "users.account, users.users")

	// Prepare test user
	user := models.User{
		Name:     "Alice Doe",
		Email:    "alice@example.com",
		Password: "securepassword",
		Phone:    "911111111",
	}
	err = repository.CreateUser(user, dbPool)
	require.NoError(t, err, "CreateUser should not return an error")

	// Retrive User ID
	retrievedUser, err := repository.GetUserByEmail(user.Email, dbPool)
	require.NoError(t, err, "GetUserByEmail should not return an error")

	// Prepare test account
	account := models.Account{
		Account_number: "ACC123456",
		Balance:        1000.0,
		Currency:       "BRL",
		Users_id:       &retrievedUser.Id,
	}
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.account (account_number, balance, currency, users_id) 
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		account.Account_number, account.Balance, account.Currency, account.Users_id,
	).Scan(&account.Id)
	require.NoError(t, err, "Failed to insert account")

	// Tests GetAccountByUserId Function
	retrievedAccount, err := repository.GetAccountByUserId(retrievedUser.Id, dbPool)
	require.NoError(t, err, "GetAccountByUserId should not return an error")

	// Verifications
	assert.Equal(t, account.Id, retrievedAccount.Id)
	assert.Equal(t, account.Account_number, retrievedAccount.Account_number)
	assert.Equal(t, account.Balance, retrievedAccount.Balance)
	assert.Equal(t, account.Currency, retrievedAccount.Currency)
	assert.Equal(t, account.Users_id, retrievedAccount.Users_id)
}
