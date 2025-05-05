package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	controllers_models "github.com/portilho13/neighborconnect-backend/models"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/marketplace"
	repositoryControllersUsers "github.com/portilho13/neighborconnect-backend/repository/controlers/users"
	models "github.com/portilho13/neighborconnect-backend/repository/models/marketplace"
	user_model "github.com/portilho13/neighborconnect-backend/repository/models/users"
)

// MockRepository is a mock for repositoryControllers
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateListingReturningId(listing models.Listing, dbPool *pgxpool.Pool) (int, error) {
	args := m.Called(listing, dbPool)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) CreateListingPhotos(photo models.Listing_Photos, dbPool *pgxpool.Pool) error {
	args := m.Called(photo, dbPool)
	return args.Error(0)
}

func (m *MockRepository) GetListingById(id int, dbPool *pgxpool.Pool) (*models.Listing, error) {
	args := m.Called(id, dbPool)
	return args.Get(0).(*models.Listing), args.Error(1)
}

func (m *MockRepository) GetListingPhotosByListingId(id int, dbPool *pgxpool.Pool) ([]models.Listing_Photos, error) {
	args := m.Called(id, dbPool)
	return args.Get(0).([]models.Listing_Photos), args.Error(1)
}

func (m *MockRepository) GetCategoryById(id int, dbPool *pgxpool.Pool) (*models.Category, error) {
	args := m.Called(id, dbPool)
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockRepository) GetBidByListningId(id int, dbPool *pgxpool.Pool) ([]models.Bid, error) {
	args := m.Called(id, dbPool)
	return args.Get(0).([]models.Bid), args.Error(1)
}

func (m *MockRepository) GetAllActiveListings(dbPool *pgxpool.Pool) ([]models.Listing, error) {
	args := m.Called(dbPool)
	return args.Get(0).([]models.Listing), args.Error(1)
}

// MockUsersRepository is a mock for repositoryControllersUsers
type MockUsersRepository struct {
	mock.Mock
}

func (m *MockUsersRepository) GetUsersById(id int, dbPool *pgxpool.Pool) (*user_model.User, error) {
	args := m.Called(id, dbPool)
	return args.Get(0).(*user_model.User), args.Error(1)
}

// MockUtils is a mock for utils
type MockUtils struct {
	mock.Mock
}

func (m *MockUtils) GetApiUrl() string {
	args := m.Called()
	return args.String(0)
}

func TestCreateListing(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockRepository)
	mockUsersRepo := new(MockUsersRepository)
	mockUtils := new(MockUtils)

	// Replace actual functions with mocks for testing
	originalRepo := repositoryControllers.CreateListingReturningId
	originalRepoCreatePhotos := repositoryControllers.CreateListingPhotos
	originalUsersRepo := repositoryControllersUsers.GetUsersById
	originalUtils := utils.GetApiUrl
	defer func() {
		repositoryControllers.CreateListingReturningId = originalRepo
		repositoryControllers.CreateListingPhotos = originalRepoCreatePhotos
		repositoryControllersUsers.GetUsersById = originalUsersRepo
		utils.GetApiUrl = originalUtils
	}()

	repositoryControllers.CreateListingReturningId = mockRepo.CreateListingReturningId
	repositoryControllers.CreateListingPhotos = mockRepo.CreateListingPhotos
	repositoryControllersUsers.GetUsersById = mockUsersRepo.GetUsersById
	utils.GetApiUrl = mockUtils.GetApiUrl

	tests := []struct {
		name           string
		listingData    controllers_models.ListingCreation
		files          []string
		setupMocks     func()
		expectedStatus int
		wantErr       bool
	}{
		{
			name: "successful listing creation with images",
			listingData: controllers_models.ListingCreation{
				Name:            "Test Listing",
				Description:     "Test Description",
				Buy_Now_Price:   "100",
				Start_Price:     "50",
				Expiration_Date: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				Seller_Id:       "1",
				Category_Id:     "2",
			},
			files: []string{"test1.jpg", "test2.jpg"},
			setupMocks: func() {
				mockRepo.On("CreateListingReturningId", mock.Anything, mock.Anything).Return(123, nil)
				mockRepo.On("CreateListingPhotos", mock.Anything, mock.Anything).Return(nil).Times(2)
				mockUtils.On("GetApiUrl").Return("localhost:8080")
			},
			expectedStatus: http.StatusOK,
			wantErr:       false,
		},
		{
			name: "invalid multipart form",
			setupMocks: func() {},
			expectedStatus: http.StatusBadRequest,
			wantErr:       true,
		},
		{
			name: "missing listing data",
			setupMocks: func() {},
			expectedStatus: http.StatusBadRequest,
			wantErr:       true,
		},
		{
			name: "invalid buy now price",
			listingData: controllers_models.ListingCreation{
				Buy_Now_Price: "invalid",
			},
			setupMocks: func() {},
			expectedStatus: http.StatusBadRequest,
			wantErr:       true,
		},
		{
			name: "invalid expiration date format",
			listingData: controllers_models.ListingCreation{
				Buy_Now_Price:   "100",
				Start_Price:     "50",
				Expiration_Date: "invalid-date",
			},
			setupMocks: func() {},
			expectedStatus: http.StatusBadRequest,
			wantErr:       true,
		},
		{
			name: "failed to create listing",
			listingData: controllers_models.ListingCreation{
				Name:            "Test Listing",
				Description:     "Test Description",
				Buy_Now_Price:   "100",
				Start_Price:     "50",
				Expiration_Date: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				Seller_Id:       "1",
				Category_Id:     "2",
			},
			setupMocks: func() {
				mockRepo.On("CreateListingReturningId", mock.Anything, mock.Anything).Return(0, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			tt.setupMocks()

			// Create multipart form request
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			// Add listing data if provided
			if tt.listingData.Name != "" {
				listingJSON, _ := json.Marshal(tt.listingData)
				_ = writer.WriteField("listing", string(listingJSON))
			}

			// Add files if provided
			for _, filename := range tt.files {
				part, _ := writer.CreateFormFile("images", filename)
				_, _ = part.Write([]byte("test image content"))
			}

			writer.Close()

			// Create request
			var req *http.Request
			if tt.name == "invalid multipart form" {
				req = httptest.NewRequest("POST", "/listings", bytes.NewBufferString("invalid"))
			} else {
				req = httptest.NewRequest("POST", "/listings", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call the function
			CreateListing(rr, req, &pgxpool.Pool{})

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
			mockUsersRepo.AssertExpectations(t)
			mockUtils.AssertExpectations(t)
		})
	}
}

func TestGetListingById(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockRepository)
	mockUsersRepo := new(MockUsersRepository)

	// Replace actual functions with mocks for testing
	originalRepo := repositoryControllers.GetListingById
	originalRepoPhotos := repositoryControllers.GetListingPhotosByListingId
	originalRepoCategory := repositoryControllers.GetCategoryById
	originalRepoBids := repositoryControllers.GetBidByListningId
	originalUsersRepo := repositoryControllersUsers.GetUsersById
	defer func() {
		repositoryControllers.GetListingById = originalRepo
		repositoryControllers.GetListingPhotosByListingId = originalRepoPhotos
		repositoryControllers.GetCategoryById = originalRepoCategory
		repositoryControllers.GetBidByListningId = originalRepoBids
		repositoryControllersUsers.GetUsersById = originalUsersRepo
	}()

	repositoryControllers.GetListingById = mockRepo.GetListingById
	repositoryControllers.GetListingPhotosByListingId = mockRepo.GetListingPhotosByListingId
	repositoryControllers.GetCategoryById = mockRepo.GetCategoryById
	repositoryControllers.GetBidByListningId = mockRepo.GetBidByListningId
	repositoryControllersUsers.GetUsersById = mockUsersRepo.GetUsersById

	tests := []struct {
		name           string
		id             string
		setupMocks     func()
		expectedStatus int
	}{
		{
			name: "successful get listing",
			id:   "123",
			setupMocks: func() {
				// Mock listing
				id := 123
				categoryId := 1
				sellerId := 1
				mockRepo.On("GetListingById", 123, mock.Anything).Return(&models.Listing{
					Id:              &id,
					Name:            "Test Listing",
					Description:     "Test Description",
					Buy_Now_Price:   100,
					Start_Price:     50,
					Created_At:      time.Now(),
					Expiration_Date: time.Now().Add(24 * time.Hour),
					Status:          "active",
					Category_Id:     &categoryId,
					Seller_Id:       &sellerId,
				}, nil)

				// Mock photos
				mockRepo.On("GetListingPhotosByListingId", 123, mock.Anything).Return([]models.Listing_Photos{
					{Id: 1, Url: "http://example.com/photo1.jpg"},
				}, nil)

				// Mock category
				mockRepo.On("GetCategoryById", 1, mock.Anything).Return(&models.Category{
					Id:   &categoryId,
					Name: "Test Category",
					Url:  stringPtr("http://example.com/category.jpg"),
				}, nil)

				// Mock user
				mockUsersRepo.On("GetUsersById", 1, mock.Anything).Return(&models.User{
					Id:   1,
					Name: "Test User",
				}, nil)

				// Mock bids (empty)
				mockRepo.On("GetBidByListningId", 123, mock.Anything).Return([]models.Bid{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid id format",
			id:   "abc",
			setupMocks: func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "listing not found",
			id:   "123",
			setupMocks: func() {
				mockRepo.On("GetListingById", 123, mock.Anything).Return(&models.Listing{}, assert.AnError)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "error fetching photos",
			id:   "123",
			setupMocks: func() {
				id := 123
				categoryId := 1
				sellerId := 1
				mockRepo.On("GetListingById", 123, mock.Anything).Return(&models.Listing{
					Id:              &id,
					Category_Id:     &categoryId,
					Seller_Id:       &sellerId,
				}, nil)
				mockRepo.On("GetListingPhotosByListingId", 123, mock.Anything).Return([]models.Listing_Photos{}, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			tt.setupMocks()

			// Create request
			req := httptest.NewRequest("GET", "/listings?id="+tt.id, nil)
			
			// Create response recorder
			rr := httptest.NewRecorder()

			// Call the function
			GetListingById(rr, req, &pgxpool.Pool{})

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
			mockUsersRepo.AssertExpectations(t)
		})
	}
}

func TestGetAllListings(t *testing.T) {
	// Setup mocks
	mockRepo := new(MockRepository)
	mockUsersRepo := new(MockUsersRepository)

	// Replace actual functions with mocks for testing
	originalRepo := repositoryControllers.GetAllActiveListings
	originalRepoPhotos := repositoryControllers.GetListingPhotosByListingId
	originalRepoCategory := repositoryControllers.GetCategoryById
	originalRepoBids := repositoryControllers.GetBidByListningId
	originalUsersRepo := repositoryControllersUsers.GetUsersById
	defer func() {
		repositoryControllers.GetAllActiveListings = originalRepo
		repositoryControllers.GetListingPhotosByListingId = originalRepoPhotos
		repositoryControllers.GetCategoryById = originalRepoCategory
		repositoryControllers.GetBidByListningId = originalRepoBids
		repositoryControllersUsers.GetUsersById = originalUsersRepo
	}()

	repositoryControllers.GetAllActiveListings = mockRepo.GetAllActiveListings
	repositoryControllers.GetListingPhotosByListingId = mockRepo.GetListingPhotosByListingId
	repositoryControllers.GetCategoryById = mockRepo.GetCategoryById
	repositoryControllers.GetBidByListningId = mockRepo.GetBidByListningId
	repositoryControllersUsers.GetUsersById = mockUsersRepo.GetUsersById

	tests := []struct {
		name           string
		setupMocks     func()
		expectedStatus int
	}{
		{
			name: "successful get all listings",
			setupMocks: func() {
				// Mock listings
				id1 := 1
				categoryId1 := 1
				sellerId1 := 1
				mockRepo.On("GetAllActiveListings", mock.Anything).Return([]models.Listing{
					{
						Id:              &id1,
						Name:            "Listing 1",
						Description:     "Description 1",
						Buy_Now_Price:   100,
						Start_Price:     50,
						Created_At:      time.Now(),
						Expiration_Date: time.Now().Add(24 * time.Hour),
						Status:          "active",
						Category_Id:     &categoryId1,
						Seller_Id:       &sellerId1,
					},
				}, nil)

				// Mock photos for listing 1
				mockRepo.On("GetListingPhotosByListingId", 1, mock.Anything).Return([]models.Listing_Photos{
					{Id: 1, Url: "http://example.com/photo1.jpg"},
				}, nil)

				// Mock category for listing 1
				mockRepo.On("GetCategoryById", 1, mock.Anything).Return(&models.Category{
					Id:   &categoryId1,
					Name: "Category 1",
					Url:  stringPtr("http://example.com/category1.jpg"),
				}, nil)

				// Mock user for listing 1
				mockUsersRepo.On("GetUsersById", 1, mock.Anything).Return(&models.User{
					Id:   1,
					Name: "User 1",
				}, nil)

				// Mock bids for listing 1 (empty)
				mockRepo.On("GetBidByListningId", 1, mock.Anything).Return([]models.Bid{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "error fetching listings",
			setupMocks: func() {
				mockRepo.On("GetAllActiveListings", mock.Anything).Return([]models.Listing{}, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "error fetching photos for a listing",
			setupMocks: func() {
				id1 := 1
				categoryId1 := 1
				sellerId1 := 1
				mockRepo.On("GetAllActiveListings", mock.Anything).Return([]models.Listing{
					{
						Id:              &id1,
						Category_Id:     &categoryId1,
						Seller_Id:       &sellerId1,
					},
				}, nil)
				mockRepo.On("GetListingPhotosByListingId", 1, mock.Anything).Return([]models.Listing_Photos{}, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			tt.setupMocks()

			// Create request
			req := httptest.NewRequest("GET", "/listings", nil)
			
			// Create response recorder
			rr := httptest.NewRecorder()

			// Call the function
			GetAllListings(rr, req, &pgxpool.Pool{})

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code)

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
			mockUsersRepo.AssertExpectations(t)
		})
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}