package test

import (
	//"context"
	"context"
	"testing"

	repository "github.com/portilho13/neighborconnect-backend/repository/controlers/marketplace"
	models "github.com/portilho13/neighborconnect-backend/repository/models/marketplace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCategory(t *testing.T) {
	// Connect to the test database
	dbPool, err := GetTestDBConnection()
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}
	defer dbPool.Close()

	// Ensure the database is clean before starting
	CleanDatabase(dbPool, "marketplace.bid")
	url := "test"

	category := models.Category{
		Name: "Modern Apartment",
		Url:  &url,
	}

	// Testar a função
	err = repository.CreateCategory(category, dbPool)
	require.NoError(t, err)

	CleanDatabase(dbPool, "marketplace.category")
}

func TestGetAllCategories(t *testing.T) {
	// Connect to the test database
	dbPool, err := GetTestDBConnection()
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}
	defer dbPool.Close()

	// Ensure the database is clean before starting
	CleanDatabase(dbPool, "marketplace.category")
	url := "test"
	category := []models.Category{{
		Name: "cat1",
		Url:  &url,
	},
		{
			Name: "cat2",
			Url:  &url,
		},
		{
			Name: "cat2",
			Url:  &url,
		},
	}
	for _, l := range category {
		// Inserir cada listing individualmente
		err := repository.CreateCategory(l, dbPool)
		require.NoError(t, err)
	}

	// Buscar todos os listings
	retrievedListings, err := repository.GetAllCategories(dbPool)
	assert.NoError(t, err, "GetAllListings should not return an error")
	require.Equal(t, len(category), len(retrievedListings), "Number of listings should match")

	// Comparar cada listing recuperado com o correspondente inserido
	for i, expected := range category {
		actual := retrievedListings[i]
		assert.Equal(t, expected.Name, actual.Name, "Retrieved Name should match inserted listing")
		assert.Equal(t, expected.Url, actual.Url, "Retrieved Description should match inserted listing")
	}

}
func TestGetCategoryById(t *testing.T) {
	dbPool, err := GetTestDBConnection()
    require.NoError(t, err)
    defer dbPool.Close()

    CleanDatabase(dbPool, "marketplace.category")

    url := "test"
    category := models.Category{
        Name: "cat1",
        Url:  &url,
    }

    // Inserir e obter o ID diretamente
    var categoryId int
    err = dbPool.QueryRow(context.Background(),
        `INSERT INTO marketplace.category (name, url) 
         VALUES ($1, $2) 
         RETURNING id`,
        category.Name, 
        category.Url).Scan(&categoryId)
    require.NoError(t, err)

    // Agora você tem o categoryId para usar nos testes
    t.Logf("Created category ID: %d", categoryId)
	// Testar GetCategoryById com o ID obtido
    retrievedCategory, err := repository.GetCategoryById(categoryId, dbPool)
    require.NoError(t, err)

    assert.Equal(t, category.Name, retrievedCategory.Name)
    assert.Equal(t, *category.Url, *retrievedCategory.Url)
}
