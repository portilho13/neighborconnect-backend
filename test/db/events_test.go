package test

import (
	"context"
	"testing"
	"time"

	repository "github.com/portilho13/neighborconnect-backend/repository/controlers/events"
	models "github.com/portilho13/neighborconnect-backend/repository/models/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCommunityEvent(t *testing.T) {
	// Connect to the test database
	dbPool, err := GetTestDBConnection()
	require.NoError(t, err, "Failed to connect to test DB")
	defer dbPool.Close()

	// Clean the necessary tables
	CleanDatabase(dbPool, "users.manager, events.community_event")

	// Insert a manager to associate with the event
	var managerId int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.manager (name, email, password, phone) 
		 VALUES ('John Doe', 'john@example.com', 'securepass', '123456789') 
		 RETURNING id`).Scan(&managerId)
	require.NoError(t, err, "Manager insertion should succeed")

	expiration := time.Date(2025, time.June, 1, 0, 0, 0, 0, time.Local).UTC()

	event := models.Community_Event{
		Name:              "test",
		Percentage:        43,
		Code:              "1",
		Capacity:          20,
		Date_Time:         time.Date(2025, time.May, 30, 23, 3, 18, 120304000, time.Local).UTC(),
		Manager_Id:        &managerId,
		Event_Image:       "",
		Duration:          time.Hour * 2,
		Local:             "Barcelos",
		Current_Ocupation: 12,
		Status:            "active",
		Expiration_Date:   &expiration, // Correção aqui
	}

	// Test the repository function
	err = repository.CreateCommunityEvent(event, dbPool)
	require.NoError(t, err, "CreateCommunityEvent should not return an error")

}

func TestAddUserToCommunityEvent(t *testing.T) {
	// Conectar ao banco de testes
	dbPool, err := GetTestDBConnection()
	require.NoError(t, err, "Failed to connect to test DB")
	defer dbPool.Close()

	// Limpar tabelas envolvidas
	CleanDatabase(dbPool, "events.many_community_event_has_many_users, events.community_event, users.users, users.manager")

	// Criar usuário
	var userId int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.users (name, email, password, phone)
		 VALUES ('Jane Doe', 'jane@example.com', 'securepass', '987654321')
		 RETURNING id`).Scan(&userId)
	require.NoError(t, err, "User insertion should succeed")

	// Criar manager
	var managerId int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.manager (name, email, password, phone)
		 VALUES ('Manager One', 'manager1@example.com', 'securepass', '555123456')
		 RETURNING id`).Scan(&managerId)
	require.NoError(t, err, "Manager insertion should succeed")

	// Criar evento comunitário
	eventImage := "image.jpg"
	var eventId int
	eventDate := time.Now().UTC().Add(24 * time.Hour).UTC()
	duration := time.Hour
	status := "active"
	expirationDate := eventDate.Add(7 * 24 * time.Hour).UTC() // 7 dias depois do evento, por exemplo

	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO events.community_event
		(name, percentage, code, capacity, date_time, manager_id, event_image, duration, local, current_ocupation, status, expiration_date)
		VALUES ('Cleanup Day', 30, 'CLN001', 20, $1, $2, $3, $4, 'Main Plaza', 0, $5, $6)
		RETURNING id`,
		eventDate, managerId, eventImage, duration, status, expirationDate,
	).Scan(&eventId)
	require.NoError(t, err, "Community event insertion should succeed")

	userEvent := models.User_Event{
		User_Id:       userId,
		Event_Id:      eventId,
		IsRewarded:    false,
		ClaimedReward: false,
	}
	err = repository.AddUserToCommunityEvent(userEvent, dbPool)
	assert.NoError(t, err, "AddUserToCommunityEvent should not return an error")

	// Verificar se a associação foi feita
	var count int
	err = dbPool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM events.many_community_event_has_many_users
		 WHERE users_id = $1 AND community_event_id = $2`, userId, eventId).Scan(&count)
	require.NoError(t, err, "Query association should succeed")
	assert.Equal(t, 1, count, "User should be associated with the event")

}

func TestGetAllEvents(t *testing.T) {
	// Conectar ao banco de testes
	dbPool, err := GetTestDBConnection()
	require.NoError(t, err, "Failed to connect to test DB")
	defer dbPool.Close()

	// Limpar tabelas relevantes
	CleanDatabase(dbPool, "users.manager,events.community_event")

	// Criar manager
	var managerId int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.manager (name, email, password, phone) 
		 VALUES ('Event Manager', 'manager@example.com', 'securepass', '123456789') 
		 RETURNING id`).Scan(&managerId)
	require.NoError(t, err, "Manager insertion should succeed")

	// Criar eventos de teste
	eventImage := "event_image.jpg"
	events := []models.Community_Event{
		{
			Name:              "Community Meeting",
			Percentage:        50,
			Code:              "MEET001",
			Capacity:          30,
			Date_Time:         time.Now().UTC().Add(24 * time.Hour).UTC(),
			Manager_Id:        &managerId,
			Event_Image:       eventImage,
			Duration:          time.Hour * 2,
			Local:             "Community Hall",
			Current_Ocupation: 15,
		},
		{
			Name:              "Charity Event",
			Percentage:        75,
			Code:              "CHAR001",
			Capacity:          50,
			Date_Time:         time.Now().UTC().Add(48 * time.Hour).UTC(),
			Manager_Id:        &managerId,
			Event_Image:       eventImage,
			Duration:          time.Hour * 3,
			Local:             "City Park",
			Current_Ocupation: 25,
		},
	}

	// Inserir eventos
	for _, event := range events {
		err := repository.CreateCommunityEvent(event, dbPool)
		require.NoError(t, err, "Event creation should succeed")
	}

	// Buscar todos os eventos
	retrievedEvents, err := repository.GetAllEvents(dbPool)
	assert.NoError(t, err, "GetAllEvents should not return an error")
	require.Equal(t, len(events), len(retrievedEvents), "Number of events should match")

	// Comparar cada evento recuperado com o correspondente inserido
	for i, expected := range events {
		actual := retrievedEvents[i]
		assert.Equal(t, expected.Name, actual.Name, "Event name mismatch")
		assert.Equal(t, expected.Percentage, actual.Percentage, "Percentage mismatch")
		assert.Equal(t, expected.Code, actual.Code, "Code mismatch")
		assert.Equal(t, expected.Capacity, actual.Capacity, "Capacity mismatch")
		assert.WithinDuration(t, expected.Date_Time, actual.Date_Time, time.Second, "Date_Time mismatch")
		assert.Equal(t, expected.Event_Image, actual.Event_Image, "Event_Image mismatch")
		assert.Equal(t, expected.Duration, actual.Duration, "Duration mismatch")
		assert.Equal(t, expected.Local, actual.Local, "Local mismatch")
		assert.Equal(t, expected.Current_Ocupation, actual.Current_Ocupation, "Current_Ocupation mismatch")
	}
}

func TestGetAllEventsByManagerId(t *testing.T) {
	// Conectar ao banco de testes
	dbPool, err := GetTestDBConnection()
	require.NoError(t, err, "Failed to connect to test DB")
	defer dbPool.Close()

	// Limpar tabelas relevantes
	CleanDatabase(dbPool, "users.manager,events.community_event")

	// Criar manager
	var managerId int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.manager (name, email, password, phone) 
         VALUES ('Event Manager', 'manager@example.com', 'securepass', '123456789') 
         RETURNING id`).Scan(&managerId)
	require.NoError(t, err, "Manager insertion should succeed")

	// Criar eventos de teste
	eventImage := "event_image.jpg"
	events := []models.Community_Event{
		{
			Name:              "Community Meeting",
			Percentage:        50,
			Code:              "MEET001",
			Capacity:          30,
			Date_Time:         time.Now().UTC().Add(24 * time.Hour).UTC(),
			Manager_Id:        &managerId,
			Event_Image:       eventImage,
			Duration:          time.Hour * 2,
			Local:             "Community Hall",
			Current_Ocupation: 15,
		},
		{
			Name:              "Charity Event",
			Percentage:        75,
			Code:              "CHAR001",
			Capacity:          50,
			Date_Time:         time.Now().UTC().Add(48 * time.Hour).UTC(),
			Manager_Id:        &managerId,
			Event_Image:       eventImage,
			Duration:          time.Hour * 3,
			Local:             "City Park",
			Current_Ocupation: 25,
		},
	}

	// Inserir eventos
	for _, event := range events {
		err := repository.CreateCommunityEvent(event, dbPool)
		require.NoError(t, err, "Event creation should succeed")
	}

	t.Run("get events for existing manager", func(t *testing.T) {
		// Obter eventos pelo ID do manager
		retrievedEvents, err := repository.GetAllEventsByManagerId(managerId, dbPool)
		require.NoError(t, err, "Should not return error")
		require.Len(t, retrievedEvents, len(events), "Should return all events for the manager")

		// Verificar se todos os eventos retornados pertencem ao manager correto
		for _, event := range retrievedEvents {
			assert.Equal(t, managerId, *event.Manager_Id, "Event should belong to the correct manager")
		}

		// Verificar se os eventos estão corretos
		eventNames := make([]string, len(retrievedEvents))
		for i, event := range retrievedEvents {
			eventNames[i] = event.Name
		}

		assert.Contains(t, eventNames, events[0].Name, "First event should be present")
		assert.Contains(t, eventNames, events[1].Name, "Second event should be present")
	})

	t.Run("get events for non-existent manager", func(t *testing.T) {
		// Testar com ID de manager que não existe
		nonExistentId := 9999
		emptyEvents, err := repository.GetAllEventsByManagerId(nonExistentId, dbPool)
		require.NoError(t, err, "Should not return error for non-existent manager")
		assert.Empty(t, emptyEvents, "Should return empty slice for non-existent manager")
	})

	t.Run("verify event details", func(t *testing.T) {
		// Obter eventos para verificação detalhada
		retrievedEvents, err := repository.GetAllEventsByManagerId(managerId, dbPool)
		require.NoError(t, err)
		require.Len(t, retrievedEvents, len(events))

		// Verificar detalhes do primeiro evento
		assert.Equal(t, events[0].Name, retrievedEvents[0].Name, "Event name mismatch")
		assert.Equal(t, events[0].Percentage, retrievedEvents[0].Percentage, "Percentage mismatch")
		assert.Equal(t, events[0].Code, retrievedEvents[0].Code, "Code mismatch")
		assert.Equal(t, events[0].Capacity, retrievedEvents[0].Capacity, "Capacity mismatch")
		assert.WithinDuration(t, events[0].Date_Time, retrievedEvents[0].Date_Time, time.Second, "Date_Time mismatch")
		assert.Equal(t, events[0].Event_Image, retrievedEvents[0].Event_Image, "Event_Image mismatch")
		assert.Equal(t, events[0].Duration, retrievedEvents[0].Duration, "Duration mismatch")
		assert.Equal(t, events[0].Local, retrievedEvents[0].Local, "Local mismatch")
		assert.Equal(t, events[0].Current_Ocupation, retrievedEvents[0].Current_Ocupation, "Current_Ocupation mismatch")
	})
}

func TestGetEventById(t *testing.T) {
	// Conectar ao banco de testes
	dbPool, err := GetTestDBConnection()
	require.NoError(t, err, "Failed to connect to test DB")
	defer dbPool.Close()

	// Limpar tabelas relevantes
	CleanDatabase(dbPool, "users.manager,events.community_event")

	// Criar manager
	var managerId int
	err = dbPool.QueryRow(context.Background(),
		`INSERT INTO users.manager (name, email, password, phone) 
         VALUES ('Event Manager', 'manager@example.com', 'securepass', '123456789') 
         RETURNING id`).Scan(&managerId)
	require.NoError(t, err, "Manager insertion should succeed")

	// Criar eventos de teste
	eventImage := "event_image.jpg"
	events := []models.Community_Event{
		{
			Name:              "Community Meeting",
			Percentage:        50,
			Code:              "MEET001",
			Capacity:          30,
			Date_Time:         time.Now().UTC().Add(24 * time.Hour).UTC(),
			Manager_Id:        &managerId,
			Event_Image:       eventImage,
			Duration:          time.Hour * 2,
			Local:             "Community Hall",
			Current_Ocupation: 15,
		},
		{
			Name:              "Charity Event",
			Percentage:        75,
			Code:              "CHAR001",
			Capacity:          50,
			Date_Time:         time.Now().UTC().UTC().Add(48 * time.Hour).UTC(),
			Manager_Id:        &managerId,
			Event_Image:       eventImage,
			Duration:          time.Hour * 3,
			Local:             "City Park",
			Current_Ocupation: 25,
		},
	}

	// Inserir eventos e armazenar os IDs retornados
	var insertedEventIDs []int
	for _, e := range events {
		var id int
		status := "active"
		expiration := e.Date_Time.Add(7 * 24 * time.Hour).UTC() // Exemplo de expiração

		err := dbPool.QueryRow(context.Background(),
			`INSERT INTO events.community_event 
			 (name, percentage, code, capacity, date_time, manager_id, event_image, duration, local, current_ocupation, status, expiration_date) 
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`,
			e.Name, e.Percentage, e.Code, e.Capacity, e.Date_Time, e.Manager_Id,
			e.Event_Image, e.Duration, e.Local, e.Current_Ocupation, status, expiration,
		).Scan(&id)
		require.NoError(t, err, "Failed to insert event")
		insertedEventIDs = append(insertedEventIDs, id)
	}

	// Testar GetEventById para cada evento inserido
	for i, expected := range events {
		actual, err := repository.GetEventById(insertedEventIDs[i], dbPool)
		require.NoError(t, err, "GetEventById should not return an error")

		assert.Equal(t, expected.Name, actual.Name, "Name mismatch")
		assert.Equal(t, expected.Percentage, actual.Percentage, "Percentage mismatch")
		assert.Equal(t, expected.Code, actual.Code, "Code mismatch")
		assert.Equal(t, expected.Capacity, actual.Capacity, "Capacity mismatch")
		assert.WithinDuration(t, expected.Date_Time, actual.Date_Time, time.Second, "Date_Time mismatch")
		assert.Equal(t, *expected.Manager_Id, *actual.Manager_Id, "Manager_Id mismatch")
		assert.Equal(t, expected.Event_Image, actual.Event_Image, "Event_Image mismatch")
		assert.Equal(t, expected.Duration, actual.Duration, "Duration mismatch")
		assert.Equal(t, expected.Local, actual.Local, "Local mismatch")
		assert.Equal(t, expected.Current_Ocupation, actual.Current_Ocupation, "Current_Ocupation mismatch")
	}
}
