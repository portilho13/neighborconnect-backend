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

func TestCreateBidReturningId(t *testing.T) {
	// Connect to the test database
	dbPool, err := GetTestDBConnection()
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}
	defer dbPool.Close()

	// Ensure the database is clean before starting
	CleanDatabase(dbPool, "users.users, marketplace.bid")

	var userId, categoryid int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.users (name, email, password, phone) 
		 VALUES ('John Doe', 'john@example.com', 'securepass', '123456789') 
		 RETURNING id`).Scan(&userId)
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
		Seller_Id:       &userId,
		Category_Id:     &categoryid,
	}

	// Testar a função
	id, err := repository.CreateListingReturningId(listing, dbPool)
	require.NoError(t, err)
	require.NotNil(t, id)

	bid := models.Bid{
		Bid_Ammount: 100,
		Bid_Time:    time.Now(),
		User_Id:     &userId,
		Listing_Id:  id,
	}

	// Testar a função
	id_bid, err := repository.CreateBidReturningId(bid, dbPool)
	require.NoError(t, err)
	require.NotNil(t, id_bid)
	CleanDatabase(dbPool, "users.users, marketplace.category")
}

func TestGetBidsByListingId(t *testing.T) {
	// Conectar ao banco de testes
	dbPool, err := GetTestDBConnection()
	require.NoError(t, err, "Failed to connect to test DB")
	defer dbPool.Close()

	// Limpar tabelas relevantes
	CleanDatabase(dbPool, "users.users, marketplace.category, marketplace.listing, marketplace.bid")

	// Criar usuário
	var userId int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.users (name, email, password, phone) 
         VALUES ('John Doe', 'john@example.com', 'securepass', '123456789') 
         RETURNING id`).Scan(&userId)
	require.NoError(t, err, "User insertion should succeed")

	// Criar categoria
	var categoryId int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO marketplace.category (name, url) 
         VALUES ('test','test') RETURNING id`).Scan(&categoryId)
	require.NoError(t, err, "Category insertion should succeed")

	// Criar listing
	listing := models.Listing{
		Name:            "Test Listing",
		Description:     "Test Description",
		Buy_Now_Price:   250000,
		Start_Price:     200000,
		Created_At:      time.Now(),
		Expiration_Date: time.Now().Add(72 * time.Hour),
		Status:          "active",
		Seller_Id:       &userId,
		Category_Id:     &categoryId,
	}

	listingId, err := repository.CreateListingReturningId(listing, dbPool)
	require.NoError(t, err, "Listing creation should succeed")
	require.NotNil(t, listingId)

	// Criar bids para este listing
	bids := []models.Bid{
		{
			Bid_Ammount: 210000,
			Bid_Time:    time.Date(2025, time.May, 6, 23, 2, 18, 120304000, time.Local).UTC(),
			User_Id:     &userId,
			Listing_Id:  listingId,
		},
		{
			Bid_Ammount: 215000,
			Bid_Time:    time.Date(2025, time.May, 6, 23, 1, 18, 120304000, time.Local).UTC(),
			User_Id:     &userId,
			Listing_Id:  listingId,
		},
		{
			Bid_Ammount: 220000,
			Bid_Time:    time.Date(2025, time.May, 6, 23, 9, 18, 120304000, time.Local).UTC(),
			User_Id:     &userId,
			Listing_Id:  listingId,
		},
	}

	// Inserir bids e armazenar IDs esperados
	expectedBids := make(map[int]models.Bid)
	for _, bid := range bids {
		bidId, err := repository.CreateBidReturningId(bid, dbPool)
		require.NoError(t, err)
		require.NotNil(t, bidId)

		// Atualizar o ID no bid
		bid.Id = &bidId
		expectedBids[bidId] = bid
	}

	// Testar a função GetBidByListingId
	retrievedBids, err := repository.GetBidByListningId(*listingId, dbPool)
	require.NoError(t, err, "GetBidByListingId should not return an error")
	require.Equal(t, len(bids), len(retrievedBids), "Number of bids should match")

	// Verificar cada bid retornado
	for _, actual := range retrievedBids {
		expected, exists := expectedBids[*actual.Id]
		require.True(t, exists, "Unexpected bid ID returned")

		assert.Equal(t, expected.Bid_Ammount, actual.Bid_Ammount, "Bid amount mismatch")
		assert.Equal(t, *expected.User_Id, *actual.User_Id, "User ID mismatch")
		assert.Equal(t, *expected.Listing_Id, *actual.Listing_Id, "Listing ID mismatch")

		// Verificar a hora com margem de 1 segundo para evitar problemas de precisão
		assert.WithinDuration(t, expected.Bid_Time, actual.Bid_Time, time.Second, "Bid time mismatch")
	}

	// Limpar tabelas relevantes
	CleanDatabase(dbPool, "users.users, marketplace.category, marketplace.listing, marketplace.bid")

}
