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

func TestCreateListingReturningId(t *testing.T) {
	// Connect to the test database
	dbPool, err := GetTestDBConnection()
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}
	defer dbPool.Close()

	// Ensure the database is clean before starting
	CleanDatabase(dbPool, "users.users, marketplace.category")

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
		Created_At:      time.Now().UTC(),
		Expiration_Date: time.Now().UTC().Add(72 * time.Hour),
		Status:          "active",
		Seller_Id:       &userId,
		Category_Id:     &categoryid,
	}
	// Testar a função
	id, err := repository.CreateListingReturningId(listing, dbPool)
	require.NoError(t, err)
	require.NotNil(t, id)
	CleanDatabase(dbPool, "users.users, marketplace.category")
}
func TestGetListingById(t *testing.T) {
	// Connect to the test database
	dbPool, err := GetTestDBConnection()
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}
	defer dbPool.Close()

	// Ensure the database is clean before starting
	CleanDatabase(dbPool, "users.users, marketplace.category, marketplace.listing")

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
		Created_At:      time.Date(2025, time.May, 6, 23, 3, 18, 120304000, time.Local).UTC(),
		Expiration_Date: time.Date(2025, time.May, 7, 23, 3, 18, 120304000, time.Local).UTC(),
		Status:          "active",
		Seller_Id:       &userId,
		Category_Id:     &categoryid,
	}
	// Testar a função
	id, err := repository.CreateListingReturningId(listing, dbPool)
	require.NoError(t, err)
	require.NotNil(t, id)

	retrievedlisting, err := repository.GetListingById(*id, dbPool)
	assert.NoError(t, err, "GetListingById should not return an error")
	assert.Equal(t, listing.Name, retrievedlisting.Name, "Retrieved Name should match inserted manager")
	assert.Equal(t, listing.Description, retrievedlisting.Description, "Retrieved Description should match inserted manager")
	assert.Equal(t, listing.Buy_Now_Price, retrievedlisting.Buy_Now_Price, "Retrieved Buy_Now_Price should match inserted manager")
	assert.Equal(t, listing.Start_Price, retrievedlisting.Start_Price, "Retrieved Start_Price should match inserted manager")
	assert.Equal(t, listing.Created_At, retrievedlisting.Created_At, "Retrieved Created_At should match inserted manager")
	assert.Equal(t, listing.Expiration_Date, retrievedlisting.Expiration_Date, "Retrieved Expiration_Date should match inserted manager")
	assert.Equal(t, listing.Status, retrievedlisting.Status, "Retrieved Status should match inserted manager")
	assert.Equal(t, listing.Seller_Id, retrievedlisting.Seller_Id, "Retrieved Seller_Id should match inserted manager")
	assert.Equal(t, listing.Category_Id, retrievedlisting.Category_Id, "Retrieved Category_Id should match inserted manager")

	CleanDatabase(dbPool, "users.users, marketplace.category, marketplace.listing")
}
func TestGetAllListings(t *testing.T) {
	// Connect to the test database
	dbPool, err := GetTestDBConnection()
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}
	defer dbPool.Close()

	// Ensure the database is clean before starting
	CleanDatabase(dbPool, "users.users, marketplace.category, marketplace.listing")

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

	listing := []models.Listing{{
		Name:            "Modern Apartment",
		Description:     "Spacious 2-bedroom apartment",
		Buy_Now_Price:   250000,
		Start_Price:     200000,
		Created_At:      time.Date(2025, time.May, 6, 23, 3, 18, 120304000, time.Local).UTC(),
		Expiration_Date: time.Date(2025, time.May, 7, 23, 3, 18, 120304000, time.Local).UTC(),
		Status:          "active",
		Seller_Id:       &userId,
		Category_Id:     &categoryid,
	},
		{
			Name:            "test2",
			Description:     "test2desc",
			Buy_Now_Price:   222,
			Start_Price:     220,
			Created_At:      time.Date(2025, time.May, 6, 23, 3, 18, 120304000, time.Local).UTC(),
			Expiration_Date: time.Date(2025, time.May, 7, 23, 3, 18, 120304000, time.Local).UTC(),
			Status:          "active",
			Seller_Id:       &userId,
			Category_Id:     &categoryid,
		},
		{
			Name:            "test3",
			Description:     "test3desc",
			Buy_Now_Price:   111,
			Start_Price:     111,
			Created_At:      time.Date(2025, time.May, 6, 23, 3, 18, 120304000, time.Local).UTC(),
			Expiration_Date: time.Date(2025, time.May, 7, 23, 3, 18, 120304000, time.Local).UTC(),
			Status:          "active",
			Seller_Id:       &userId,
			Category_Id:     &categoryid,
		},
	}
	for _, l := range listing {
		// Inserir cada listing individualmente
		id, err := repository.CreateListingReturningId(l, dbPool)
		require.NoError(t, err)
		require.NotNil(t, id)
	}

	// Buscar todos os listings
	retrievedListings, err := repository.GetAllListings(dbPool)
	assert.NoError(t, err, "GetAllListings should not return an error")
	require.Equal(t, len(listing), len(retrievedListings), "Number of listings should match")

	// Comparar cada listing recuperado com o correspondente inserido
	for i, expected := range listing {
		actual := retrievedListings[i]
		assert.Equal(t, expected.Name, actual.Name, "Retrieved Name should match inserted listing")
		assert.Equal(t, expected.Description, actual.Description, "Retrieved Description should match inserted listing")
		assert.Equal(t, expected.Buy_Now_Price, actual.Buy_Now_Price, "Retrieved Buy_Now_Price should match inserted listing")
		assert.Equal(t, expected.Start_Price, actual.Start_Price, "Retrieved Start_Price should match inserted listing")
		assert.Equal(t, expected.Created_At, actual.Created_At, "Retrieved Created_At should match inserted listing")
		assert.Equal(t, expected.Expiration_Date, actual.Expiration_Date, "Retrieved Expiration_Date should match inserted listing")
		assert.Equal(t, expected.Status, actual.Status, "Retrieved Status should match inserted listing")
		assert.Equal(t, *expected.Seller_Id, *actual.Seller_Id, "Retrieved Seller_Id should match inserted listing")
		assert.Equal(t, *expected.Category_Id, *actual.Category_Id, "Retrieved Category_Id should match inserted listing")
	}

}

func TestGetAllActiveListings(t *testing.T) {
	// Connect to the test database
	dbPool, err := GetTestDBConnection()
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}
	defer dbPool.Close()

	// Ensure the database is clean before starting
	CleanDatabase(dbPool, "users.users, marketplace.category, marketplace.listing")

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

	listing := []models.Listing{{
		Name:            "Modern Apartment",
		Description:     "Spacious 2-bedroom apartment",
		Buy_Now_Price:   250000,
		Start_Price:     200000,
		Created_At:      time.Date(2025, time.May, 6, 23, 3, 18, 120304000, time.Local).UTC(),
		Expiration_Date: time.Date(2025, time.May, 7, 23, 3, 18, 120304000, time.Local).UTC(),
		Status:          "active",
		Seller_Id:       &userId,
		Category_Id:     &categoryid,
	},
		{
			Name:            "test2",
			Description:     "test2desc",
			Buy_Now_Price:   222,
			Start_Price:     220,
			Created_At:      time.Date(2025, time.May, 6, 23, 3, 18, 120304000, time.Local).UTC(),
			Expiration_Date: time.Date(2025, time.May, 7, 23, 3, 18, 120304000, time.Local).UTC(),
			Status:          "inactive",
			Seller_Id:       &userId,
			Category_Id:     &categoryid,
		},
		{
			Name:            "test3",
			Description:     "test3desc",
			Buy_Now_Price:   111,
			Start_Price:     111,
			Created_At:      time.Date(2025, time.May, 6, 23, 3, 18, 120304000, time.Local).UTC(),
			Expiration_Date: time.Date(2025, time.May, 7, 23, 3, 18, 120304000, time.Local).UTC(),
			Status:          "active",
			Seller_Id:       &userId,
			Category_Id:     &categoryid,
		},
	}
	expectedListings := make(map[int]models.Listing)

	for _, l := range listing {
		id, err := repository.CreateListingReturningId(l, dbPool)
		require.NoError(t, err)
		require.NotNil(t, id)

		// Mapear apenas os que são "active"
		if l.Status == "active" {
			expectedListings[*id] = l
		}
	}

	// Buscar todos os listings ativos
	retrievedListings, err := repository.GetAllActiveListings(dbPool)
	assert.NoError(t, err, "GetAllActiveListings should not return an error")

	// Verificar que cada listing ativo retornado bate com o esperado
	for _, actual := range retrievedListings {
		expected, exists := expectedListings[*actual.Id]
		assert.True(t, exists, "Unexpected listing with ID %d returned", actual.Id)

		assert.Equal(t, expected.Name, actual.Name, "Name mismatch for ID %d", actual.Id)
		assert.Equal(t, expected.Description, actual.Description, "Description mismatch for ID %d", actual.Id)
		assert.Equal(t, expected.Buy_Now_Price, actual.Buy_Now_Price, "Buy_Now_Price mismatch for ID %d", actual.Id)
		assert.Equal(t, expected.Start_Price, actual.Start_Price, "Start_Price mismatch for ID %d", actual.Id)
		assert.Equal(t, expected.Created_At, actual.Created_At, "Created_At mismatch for ID %d", actual.Id)
		assert.Equal(t, expected.Expiration_Date, actual.Expiration_Date, "Expiration_Date mismatch for ID %d", actual.Id)
		assert.Equal(t, "active", actual.Status, "Status mismatch for ID %d", actual.Id)
		assert.Equal(t, *expected.Seller_Id, *actual.Seller_Id, "Seller_Id mismatch for ID %d", actual.Id)
		assert.Equal(t, *expected.Category_Id, *actual.Category_Id, "Category_Id mismatch for ID %d", actual.Id)
	}
}

func TestUpdateListingStatus(t *testing.T) {
	// Connect to the test database
	dbPool, err := GetTestDBConnection()
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}
	defer dbPool.Close()

	// Ensure the database is clean before starting
	CleanDatabase(dbPool, "users.users, marketplace.category, marketplace.listing")

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
		Created_At:      time.Date(2025, time.May, 6, 23, 3, 18, 120304000, time.Local).UTC(),
		Expiration_Date: time.Date(2025, time.May, 7, 23, 3, 18, 120304000, time.Local).UTC(),
		Status:          "active",
		Seller_Id:       &userId,
		Category_Id:     &categoryid,
	}
	// Testar a função
	id, err := repository.CreateListingReturningId(listing, dbPool)
	require.NoError(t, err)
	require.NotNil(t, id)

	err = repository.UpdateListingStatus("closed", *id, dbPool)
	assert.NoError(t, err, "UpdateListingStatus should not return an error")

	retrievedlisting, err := repository.GetListingById(*id, dbPool)
	assert.NoError(t, err, "GetListingById should not return an error")
	assert.Equal(t, listing.Name, retrievedlisting.Name, "Retrieved Name should match inserted manager")
	assert.Equal(t, listing.Description, retrievedlisting.Description, "Retrieved Description should match inserted manager")
	assert.Equal(t, listing.Buy_Now_Price, retrievedlisting.Buy_Now_Price, "Retrieved Buy_Now_Price should match inserted manager")
	assert.Equal(t, listing.Start_Price, retrievedlisting.Start_Price, "Retrieved Start_Price should match inserted manager")
	assert.Equal(t, listing.Created_At, retrievedlisting.Created_At, "Retrieved Created_At should match inserted manager")
	assert.Equal(t, listing.Expiration_Date, retrievedlisting.Expiration_Date, "Retrieved Expiration_Date should match inserted manager")
	assert.Equal(t, "closed", retrievedlisting.Status, "Retrieved Status should match inserted manager")
	assert.Equal(t, listing.Seller_Id, retrievedlisting.Seller_Id, "Retrieved Seller_Id should match inserted manager")
	assert.Equal(t, listing.Category_Id, retrievedlisting.Category_Id, "Retrieved Category_Id should match inserted manager")

	CleanDatabase(dbPool, "users.users, marketplace.category, marketplace.listing")
}

func TestGetListingsBySellerId(t *testing.T) {
	// Conectar ao banco de testes
	dbPool, err := GetTestDBConnection()
	require.NoError(t, err, "Failed to connect to test DB")
	defer dbPool.Close()

	// Limpar tabelas relevantes
	CleanDatabase(dbPool, "users.users, marketplace.category, marketplace.listing")

	// Criar 2 usuários diferentes
	var seller1Id, seller2Id int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.users (name, email, password, phone) 
	 VALUES ('Seller One', 'seller1@example.com', 'pass1', '111111111') 
	 RETURNING id`).Scan(&seller1Id)
	require.NoError(t, err, "First user insertion should succeed")

	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.users (name, email, password, phone) 
	 VALUES ('Seller Two', 'seller2@example.com', 'pass2', '222222222') 
	 RETURNING id`).Scan(&seller2Id)
	require.NoError(t, err, "Second user insertion should succeed")

	// Criar categoria
	var categoryId int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO marketplace.category (name, url) 
		VALUES ('test','test') RETURNING id`).Scan(&categoryId)
	assert.NoError(t, err, "Category insertion should succeed")

	listings := []models.Listing{{
		Name:            "Modern Apartment",
		Description:     "Spacious 2-bedroom apartment",
		Buy_Now_Price:   250000,
		Start_Price:     200000,
		Created_At:      time.Date(2025, time.May, 6, 23, 3, 18, 120304000, time.Local).UTC(),
		Expiration_Date: time.Date(2025, time.May, 7, 23, 3, 18, 120304000, time.Local).UTC(),
		Status:          "active",
		Seller_Id:       &seller1Id,
		Category_Id:     &categoryId,
	},
		{
			Name:            "test2",
			Description:     "test2desc",
			Buy_Now_Price:   222,
			Start_Price:     220,
			Created_At:      time.Date(2025, time.May, 6, 23, 3, 18, 120304000, time.Local).UTC(),
			Expiration_Date: time.Date(2025, time.May, 7, 23, 3, 18, 120304000, time.Local).UTC(),
			Status:          "inactive",
			Seller_Id:       &seller2Id,
			Category_Id:     &categoryId,
		},
		{
			Name:            "test3",
			Description:     "test3desc",
			Buy_Now_Price:   111,
			Start_Price:     111,
			Created_At:      time.Date(2025, time.May, 6, 23, 3, 18, 120304000, time.Local).UTC(),
			Expiration_Date: time.Date(2025, time.May, 7, 23, 3, 18, 120304000, time.Local).UTC(),
			Status:          "active",
			Seller_Id:       &seller2Id,
			Category_Id:     &categoryId,
		},
	}

	// Mapa para armazenar os listings esperados por seller
	expectedListings := make(map[int][]models.Listing)

	for _, l := range listings {
		id, err := repository.CreateListingReturningId(l, dbPool)
		require.NoError(t, err)
		require.NotNil(t, id)
	}

	t.Run("get listings for seller 1", func(t *testing.T) {
		retrieved, err := repository.GetListingsBySellerId(seller2Id, dbPool)
		require.NoError(t, err)
		require.Len(t, retrieved, 2, "Should return 1 listings for seller 2")

		// Verificar se todos os listings retornados pertencem ao seller1Id
		for _, listing := range retrieved {
			assert.Equal(t, seller2Id, *listing.Seller_Id)
		}

		// Verificar detalhes dos listings
		for _, expected := range expectedListings[seller2Id] {
			found := false
			for _, actual := range retrieved {
				if *actual.Id == *expected.Id {
					found = true
					assert.Equal(t, expected.Name, actual.Name)
					assert.Equal(t, expected.Description, actual.Description)
					assert.Equal(t, expected.Buy_Now_Price, actual.Buy_Now_Price)
					assert.Equal(t, expected.Start_Price, actual.Start_Price)
					assert.Equal(t, expected.Status, actual.Status)
					break
				}
			}
			assert.True(t, found, "Listing with ID %d not found in results", expected.Id)
		}
	})
}
