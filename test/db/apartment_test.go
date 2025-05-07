package test

import (
	"context"
	"testing"
	// "time"

	repository "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
	models "github.com/portilho13/neighborconnect-backend/repository/models/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateApartment(t *testing.T) {
	// Conectar ao banco de testes
	dbPool, err := GetTestDBConnection()
	require.NoError(t, err, "Failed to connect to test DB")
	defer dbPool.Close()

	// Limpar as tabelas relevantes
	CleanDatabase(dbPool, "users.apartment, users.manager")

	// Inserir um manager necessário para o apartamento
	var managerId int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.manager (name, email, password, phone) 
		 VALUES ('Manager Joe', 'manager@example.com', 'securepass', '123456789') 
		 RETURNING id`).Scan(&managerId)
	require.NoError(t, err, "Manager insertion should succeed")

	// Criar o apartamento a ser testado
	apartment := models.Apartment{
		N_bedrooms:  2,
		Floor:       3,
		Rent:        1500.0,
		Manager_id:  managerId,
		Status:      "available",
	}

	// Testar a criação do apartamento
	err = repository.CreateApartment(apartment, dbPool)
	require.NoError(t, err, "CreateApartment should not return an error")
}
	
func TestUpdateApartmentStatus(t *testing.T) {
	// Conectar ao banco de dados de teste
	dbPool, err := GetTestDBConnection()
	require.NoError(t, err, "Failed to connect to test DB")
	defer dbPool.Close()

	// Limpar as tabelas necessárias
	CleanDatabase(dbPool, "users.apartment, users.manager")

	// Inserir um manager
	var managerId int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.manager (name, email, password, phone)
		 VALUES ('Jane Manager', 'jane@example.com', 'securepass', '987654321')
		 RETURNING id`).Scan(&managerId)
	require.NoError(t, err, "Manager insertion should succeed")

	// Inserir um apartamento com status "available"
	var apartmentId int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.apartment (n_bedrooms, floor, rent, manager_id, status)
		 VALUES (3, 5, 1800, $1, 'available') RETURNING id`, managerId).Scan(&apartmentId)
	require.NoError(t, err, "Apartment insertion should succeed")

	// Executar a função de atualização
	err = repository.UpdateApartmentStatus(apartmentId, dbPool)
	require.NoError(t, err, "UpdateApartmentStatus should not return an error")

	// Verificar se o status foi atualizado
	var updatedStatus string
	err = dbPool.QueryRow(context.Background(),
		`SELECT status FROM users.apartment WHERE id = $1`, apartmentId).Scan(&updatedStatus)
	require.NoError(t, err, "Failed to fetch updated apartment")
	assert.Equal(t, "occupied", updatedStatus, "Apartment status should be updated to 'occupied'")

}


func TestGetAllApartments(t *testing.T) {
	// Conectar ao banco de dados de teste
	dbPool, err := GetTestDBConnection()
	require.NoError(t, err, "Failed to connect to test DB")
	defer dbPool.Close()

	// Limpar as tabelas necessárias
	CleanDatabase(dbPool, "users.apartment, users.manager")

	// Inserir um manager
	var managerId int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.manager (name, email, password, phone)
		 VALUES ('Jane Manager', 'jane@example.com', 'securepass', '987654321')
		 RETURNING id`).Scan(&managerId)
	require.NoError(t, err, "Manager insertion should succeed")

	// Inserir múltiplos apartamentos
	expectedApartments := []models.Apartment{
		{N_bedrooms: 2, Floor: 1, Rent: 1500, Manager_id: managerId, Status: "available"},
		{N_bedrooms: 3, Floor: 2, Rent: 1800, Manager_id: managerId, Status: "occupied"},
		{N_bedrooms: 1, Floor: 3, Rent: 1200, Manager_id: managerId, Status: "available"},
	}

	for i := range expectedApartments {
		err := dbPool.QueryRow(context.Background(),
			`INSERT INTO users.apartment (n_bedrooms, floor, rent, manager_id, status)
			 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			expectedApartments[i].N_bedrooms,
			expectedApartments[i].Floor,
			expectedApartments[i].Rent,
			expectedApartments[i].Manager_id,
			expectedApartments[i].Status,
		).Scan(&expectedApartments[i].Id)
		require.NoError(t, err, "Apartment insertion should succeed")
	}

	// Chamar a função a ser testada
	retrievedApartments, err := repository.GetAllApartments(dbPool)
	require.NoError(t, err, "GetAllApartments should not return an error")
	assert.Equal(t, len(expectedApartments), len(retrievedApartments), "Number of apartments should match")

	// Verificar os dados
	for i, expected := range expectedApartments {
		actual := retrievedApartments[i]
		assert.Equal(t, expected.N_bedrooms, actual.N_bedrooms, "Bedrooms should match")
		assert.Equal(t, expected.Floor, actual.Floor, "Floor should match")
		assert.Equal(t, expected.Rent, actual.Rent, "Rent should match")
		assert.Equal(t, expected.Manager_id, actual.Manager_id, "Manager ID should match")
		assert.Equal(t, expected.Status, actual.Status, "Status should match")
	}
}
func TestGetAllOccupiedApartments(t *testing.T) {
	// Conectar ao banco de dados de teste
	dbPool, err := GetTestDBConnection()
	require.NoError(t, err, "Failed to connect to test DB")
	defer dbPool.Close()

	// Limpar tabelas relevantes
	CleanDatabase(dbPool, "users.apartment, users.manager")

	// Inserir um manager para vincular aos apartamentos
	var managerId int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.manager (name, email, password, phone)
		 VALUES ('Manager Jane', 'jane@example.com', 'securepass', '123123123') RETURNING id`).Scan(&managerId)
	require.NoError(t, err)

	// Inserir apartamentos com status variados
	apartments := []models.Apartment{
		{N_bedrooms: 2, Floor: 1, Rent: 1000, Manager_id: managerId, Status: "occupied"},
		{N_bedrooms: 1, Floor: 2, Rent: 900, Manager_id: managerId, Status: "available"},
		{N_bedrooms: 3, Floor: 3, Rent: 1200, Manager_id: managerId, Status: "occupied"},
	}

	expectedOccupied := []models.Apartment{}

	for i := range apartments {
		err := dbPool.QueryRow(context.Background(),
			`INSERT INTO users.apartment (n_bedrooms, floor, rent, manager_id, status)
			 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			apartments[i].N_bedrooms,
			apartments[i].Floor,
			apartments[i].Rent,
			apartments[i].Manager_id,
			apartments[i].Status,
		).Scan(&apartments[i].Id)
		require.NoError(t, err)

		if apartments[i].Status == "occupied" {
			expectedOccupied = append(expectedOccupied, apartments[i])
		}
	}

	// Executar a função sob teste
	retrieved, err := repository.GetAllOccupiedApartments(dbPool)
	require.NoError(t, err, "GetAllOccupiedApartments should not return an error")
	assert.Equal(t, len(expectedOccupied), len(retrieved), "Number of occupied apartments should match")
}


func TestGetApartmentById(t *testing.T) {
	// Conectar ao banco de dados de teste
	dbPool, err := GetTestDBConnection()
	require.NoError(t, err)
	defer dbPool.Close()

	// Limpar tabelas relevantes
	CleanDatabase(dbPool, "users.apartment, users.manager")

	// Inserir um manager necessário para o relacionamento
	var managerId int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.manager (name, email, password, phone)
		 VALUES ('Manager Jane', 'jane@example.com', 'securepass', '123123123') RETURNING id`).Scan(&managerId)
	require.NoError(t, err)

	// Criar um apartamento para testar
	apartment := models.Apartment{
		N_bedrooms: 2,
		Floor:      5,
		Rent:       1500,
		Manager_id: managerId,
		Status:     "occupied",
	}

	// Inserir o apartamento no banco de dados
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.apartment (n_bedrooms, floor, rent, manager_id, status)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		apartment.N_bedrooms,
		apartment.Floor,
		apartment.Rent,
		apartment.Manager_id,
		apartment.Status,
	).Scan(&apartment.Id)
	require.NoError(t, err)

	// Buscar pelo ID
	retrieved, err := repository.GetApartmentById(*apartment.Id, dbPool)
	require.NoError(t, err, "GetApartmentById should not return an error")

	// Verificar se os dados estão corretos
	assert.Equal(t, apartment.Id, retrieved.Id)
	assert.Equal(t, apartment.N_bedrooms, retrieved.N_bedrooms)
	assert.Equal(t, apartment.Floor, retrieved.Floor)
	assert.Equal(t, apartment.Rent, retrieved.Rent)
	assert.Equal(t, apartment.Manager_id, retrieved.Manager_id)
	assert.Equal(t, apartment.Status, retrieved.Status)
}



func TestGetAllApartmentsByManagerId(t *testing.T) {
	// Conectar ao banco de dados de teste
	dbPool, err := GetTestDBConnection()
	require.NoError(t, err)
	defer dbPool.Close()

	// Limpar tabelas necessárias
	CleanDatabase(dbPool, "users.apartment, users.manager")

	// Inserir dois managers
	var manager1Id, manager2Id int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.manager (name, email, password, phone) 
		 VALUES ('Manager One', 'manager1@example.com', 'pass1', '111111111') 
		 RETURNING id`).Scan(&manager1Id)
	require.NoError(t, err)

	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.manager (name, email, password, phone) 
		 VALUES ('Manager Two', 'manager2@example.com', 'pass2', '222222222') 
		 RETURNING id`).Scan(&manager2Id)
	require.NoError(t, err)

	// Inserir apartamentos para ambos os managers
	apartments := []models.Apartment{
		{N_bedrooms: 2, Floor: 1, Rent: 1000, Manager_id: manager1Id, Status: "available"},
		{N_bedrooms: 3, Floor: 2, Rent: 1200, Manager_id: manager1Id, Status: "occupied"},
		{N_bedrooms: 1, Floor: 3, Rent: 900, Manager_id: manager2Id, Status: "available"},
	}

	expectedByManager := make(map[int][]models.Apartment)

	for i := range apartments {
		ap := &apartments[i]
		err := dbPool.QueryRow(context.Background(),
			`INSERT INTO users.apartment (n_bedrooms, floor, rent, manager_id, status) 
			 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			ap.N_bedrooms, ap.Floor, ap.Rent, ap.Manager_id, ap.Status,
		).Scan(&ap.Id)
		require.NoError(t, err)

		expectedByManager[ap.Manager_id] = append(expectedByManager[ap.Manager_id], *ap)
	}

	t.Run("get listings for manager 1", func(t *testing.T) {
		retrieved, err := repository.GetAllApartmentsByManagerId(manager1Id, dbPool)
		require.NoError(t, err)
		require.Len(t, retrieved, 2, "Should return 1 listings for seller 2")

		// Verificar se todos os listings retornados pertencem ao seller1Id
		for _, apartment := range retrieved {
			assert.Equal(t, manager1Id, *&apartment.Manager_id)
		}

		// Verificar detalhes dos listings
		for _, expected := range expectedByManager[manager1Id] {
			found := false
			for _, actual := range retrieved {
				if *actual.Id == *expected.Id {
					found = true
					assert.Equal(t, expected.N_bedrooms, actual.N_bedrooms)
					assert.Equal(t, expected.Floor, actual.Floor)
					assert.Equal(t, expected.Rent, actual.Rent)
					assert.Equal(t, expected.Manager_id, actual.Manager_id)
					assert.Equal(t, expected.Status, actual.Status)
					break
				}
			}
			assert.True(t, found, "Listing with ID %d not found in results", expected.Id)
		}
	 })
	

	t.Run("get apartments for manager 2", func(t *testing.T) {
		retrieved, err := repository.GetAllApartmentsByManagerId(manager2Id, dbPool)
		require.NoError(t, err)
		assert.Len(t, retrieved, 1)

		expected := expectedByManager[manager2Id][0]
		actual := retrieved[0]

		assert.Equal(t, expected.Id, actual.Id)
		assert.Equal(t, expected.N_bedrooms, actual.N_bedrooms)
		assert.Equal(t, expected.Floor, actual.Floor)
		assert.Equal(t, expected.Rent, actual.Rent)
		assert.Equal(t, expected.Manager_id, actual.Manager_id)
		assert.Equal(t, expected.Status, actual.Status)
	})
}