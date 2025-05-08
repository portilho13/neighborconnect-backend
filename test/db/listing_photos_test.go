package test

import (
	"context"
	"testing"
	"time"

	repository "github.com/portilho13/neighborconnect-backend/repository/controlers/marketplace"
	models "github.com/portilho13/neighborconnect-backend/repository/models/marketplace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateListingPhotos(t *testing.T) {
	// Connect to the test database
	dbPool, err := GetTestDBConnection()
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}
	defer dbPool.Close()

	// Ensure the database is clean before starting
	CleanDatabase(dbPool, "users.users, marketplace.category, marketplace.listing_photos, marketplace.listing")

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
	listingid, err := repository.CreateListingReturningId(listing, dbPool)
	require.NoError(t, err)
	require.NotNil(t, listingid)

	listing_photos := models.Listing_Photos{
		Url:        "www.sss.ss",
		Listing_Id: listingid,
	}
	// Testar a função
	err = repository.CreateListingPhotos(listing_photos, dbPool)
	require.NoError(t, err)

	CleanDatabase(dbPool, "users.users, marketplace.category, marketplace.listing_photos, marketplace.listing")
}
func TestGetListingPhotosByListingId(t *testing.T) {
	// Conectar ao banco de testes
	dbPool, err := GetTestDBConnection()
	require.NoError(t, err, "Failed to connect to test DB")
	defer dbPool.Close()

	// Limpar tabelas relevantes
	CleanDatabase(dbPool, "users.users,marketplace.category,marketplace.listing,marketplace.listing_photos")

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
		Name:            "Modern Apartment",
		Description:     "Spacious 2-bedroom apartment",
		Buy_Now_Price:   250000,
		Start_Price:     200000,
		Created_At:      time.Now().UTC(),
		Expiration_Date: time.Now().UTC().Add(24 * time.Hour),
		Status:          "active",
		Seller_Id:       &userId,
		Category_Id:     &categoryId,
	}

	listingId, err := repository.CreateListingReturningId(listing, dbPool)
	require.NoError(t, err, "Listing creation should succeed")
	require.NotNil(t, listingId)

	// Criar fotos para o listing
	photos := []models.Listing_Photos{
		{
			Url:        "http://example.com/photo1.jpg",
			Listing_Id: listingId,
		},
		{
			Url:        "http://example.com/photo2.jpg",
			Listing_Id: listingId,
		},
	}

	// Inserir fotos
	for _, photo := range photos {
		err := repository.CreateListingPhotos(photo, dbPool)
		require.NoError(t, err, "Photo insertion should succeed")
	}

	t.Run("get photos for existing listing", func(t *testing.T) {
		// Obter fotos pelo ID do listing
		retrievedPhotos, err := repository.GetListingPhotosByListingId(*listingId, dbPool)
		require.NoError(t, err, "Should not return error")
		require.Len(t, retrievedPhotos, len(photos), "Should return all photos for the listing")

		// Verificar se todas as fotos retornadas pertencem ao listing correto
		for _, photo := range retrievedPhotos {
			assert.Equal(t, *listingId, *photo.Listing_Id, "Photo should belong to the correct listing")
		}

		// Verificar se as URLs estão corretas
		photoUrls := make([]string, len(retrievedPhotos))
		for i, photo := range retrievedPhotos {
			photoUrls[i] = photo.Url
		}

		assert.Contains(t, photoUrls, photos[0].Url, "First photo URL should be present")
		assert.Contains(t, photoUrls, photos[1].Url, "Second photo URL should be present")
	})

	t.Run("get photos for non-existent listing", func(t *testing.T) {
		// Testar com ID de listing que não existe
		nonExistentId := 9999
		emptyPhotos, err := repository.GetListingPhotosByListingId(nonExistentId, dbPool)
		require.NoError(t, err, "Should not return error for non-existent listing")
		assert.Empty(t, emptyPhotos, "Should return empty slice for non-existent listing")
	})
}
