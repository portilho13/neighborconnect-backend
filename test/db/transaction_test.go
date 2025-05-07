package test

import (
	"context"
	"testing"
	"time"

	// "github.com/pashagolub/pgxmock"
	repository "github.com/portilho13/neighborconnect-backend/repository/controlers/marketplace"
	// repository_users "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
	models "github.com/portilho13/neighborconnect-backend/repository/models/marketplace"
	// models_users "github.com/portilho13/neighborconnect-backend/repository/models/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTransaction(t *testing.T) {
	// Connect to the test database
	dbPool, err := GetTestDBConnection()
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}
	defer dbPool.Close()
	// Ensure the database is clean before starting
	CleanDatabase(dbPool, "users.users, marketplace.category, marketplace.listing, marketplace.transaction")

	var sellerId, buyerId, categoryid int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.users (name, email, password, phone) 
		 VALUES ('John Doe', 'john@example.com', 'securepass', '123456789') 
		 RETURNING id`).Scan(&sellerId)
	assert.NoError(t, err, "User insertion should succeed")
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.users (name, email, password, phone) 
		 VALUES ('test', 'test@example.com', 'securepass', '123456789') 
		 RETURNING id`).Scan(&buyerId)
	assert.NoError(t, err, "User insertion should succeed")

	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO marketplace.category (name, url) 
		VALUES ('test','test') RETURNING id`).Scan(&categoryid)
	assert.NoError(t, err, "Category insertion should succeed")

	listing := models.Listing{
		Name:            "Modern Apartment",
		Description:     "Spacious 2-bedroom apartment",
		Buy_Now_Price:   250000,
		Start_Price:     200000,
		Created_At:      time.Now(),
		Expiration_Date: time.Now().Add(72 * time.Hour),
		Status:          "active",
		Seller_Id:       &sellerId,
		Category_Id:     &categoryid,
	}
	// Testing function CreateListingReturningId
	listingId, err := repository.CreateListingReturningId(listing, dbPool)
	require.NoError(t, err)
	require.NotNil(t, listingId)

	transaction := models.Transaction{
		Final_Price:      100,
		Transaction_Time: time.Date(2025, time.May, 6, 23, 3, 18, 120304000, time.Local).UTC(),
		Transaction_Type: "mb",
		Listing_Id:       listingId,
		Buyer_Id:         &buyerId,
		Payment_Status:   "paid",
		Payment_Due_time: time.Date(2025, time.May, 7, 23, 3, 18, 120304000, time.Local).UTC(),
		Seller_Id:        &sellerId,
	}
	// Testing function CreateTransaction
	err = repository.CreateTransaction(transaction, dbPool)
	require.NoError(t, err)

	CleanDatabase(dbPool, "users.users, marketplace.category, marketplace.listing, marketplace.transaction")

}
func TestGetTransactionById(t *testing.T) {
	// Connect to the test database
	dbPool, err := GetTestDBConnection()
	require.NoError(t, err, "Failed to connect to test DB")
	defer dbPool.Close()

	// Ensure the database is clean before starting
	CleanDatabase(dbPool, "users.users,marketplace.category,marketplace.listing,marketplace.transaction")

	// Create test seller
	var sellerId int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.users (name, email, password, phone) 
         VALUES ('John Doe', 'john@example.com', 'securepass', '123456789') 
         RETURNING id`).Scan(&sellerId)
	require.NoError(t, err, "Seller insertion should succeed")

	// Create test buyer
	var buyerId int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.users (name, email, password, phone) 
         VALUES ('Buyer', 'buyer@example.com', 'securepass', '987654321') 
         RETURNING id`).Scan(&buyerId)
	require.NoError(t, err, "Buyer insertion should succeed")

	// Create test category
	var categoryId int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO marketplace.category (name, url) 
         VALUES ('Apartments', 'apartments') 
         RETURNING id`).Scan(&categoryId)
	require.NoError(t, err, "Category insertion should succeed")

	// Create test listing
	listing := models.Listing{
		Name:            "Modern Apartment",
		Description:     "Spacious 2-bedroom apartment",
		Buy_Now_Price:   250000,
		Start_Price:     200000,
		Created_At:      time.Now(),
		Expiration_Date: time.Now().Add(72 * time.Hour),
		Status:          "active",
		Seller_Id:       &sellerId,
		Category_Id:     &categoryId,
	}

	listingId, err := repository.CreateListingReturningId(listing, dbPool)
	require.NoError(t, err, "Listing creation should succeed")
	require.NotNil(t, listingId)

	// Create test transação
	transaction := models.Transaction{
		Final_Price:      220000,
		Transaction_Time: time.Date(2025, time.May, 5, 23, 3, 18, 120304000, time.Local).UTC(),
		Transaction_Type: "mbway",
		Listing_Id:       listingId,
		Buyer_Id:         &buyerId,
		Payment_Status:   "paid",
		Payment_Due_time: time.Date(2025, time.May, 10, 23, 3, 18, 120304000, time.Local).UTC(),
		Seller_Id:        &sellerId,
	}

	// insert transaction to receive Id
	var transactionId int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO marketplace.transaction 
         (final_price, transaction_time, transaction_type, seller_id, buyer_id, listing_id, payment_status, payment_due_time) 
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
         RETURNING id`,
		transaction.Final_Price,
		transaction.Transaction_Time,
		transaction.Transaction_Type,
		transaction.Seller_Id,
		transaction.Buyer_Id,
		transaction.Listing_Id,
		transaction.Payment_Status,
		transaction.Payment_Due_time).Scan(&transactionId)
	require.NoError(t, err, "Transaction insertion should succeed")

	// Gets Transaction by Id
	retrievedTransaction, err := repository.GetTransactionById(transactionId, dbPool)
	require.NoError(t, err, "Should not return error for existing transaction")

	// Verifications
	assert.Equal(t, transaction.Final_Price, retrievedTransaction.Final_Price, "Final price mismatch")
	assert.WithinDuration(t, transaction.Transaction_Time, retrievedTransaction.Transaction_Time, time.Second, "Transaction time mismatch")
	assert.Equal(t, transaction.Transaction_Type, retrievedTransaction.Transaction_Type, "Transaction type mismatch")
	assert.Equal(t, *transaction.Seller_Id, *retrievedTransaction.Seller_Id, "Seller ID mismatch")
	assert.Equal(t, *transaction.Buyer_Id, *retrievedTransaction.Buyer_Id, "Buyer ID mismatch")
	assert.Equal(t, *transaction.Listing_Id, *retrievedTransaction.Listing_Id, "Listing ID mismatch")
	assert.Equal(t, transaction.Payment_Status, retrievedTransaction.Payment_Status, "Payment status mismatch")
	assert.WithinDuration(t, transaction.Payment_Due_time, retrievedTransaction.Payment_Due_time, time.Second, "Payment due time mismatch")

}
