package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/portilho13/neighborconnect-backend/controllers"
)

func TestTestAPI(t *testing.T) {
	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/api/v1/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Call the TestAPI handler
	handler := http.HandlerFunc(controllers.TestAPI)
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("TestAPI returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check if the response content type is JSON
	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("TestAPI returned wrong content type: got %v want application/json", contentType)
	}

	// Check the response body
	var response map[string]string
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal("Failed to decode JSON response:", err)
	}

	if response["message"] != "API Worknig!" {
		t.Errorf("TestAPI returned unexpected message: got %v want API Working!", response["message"])
	}
}

func TestTestAPI_MethodNotAllowed(t *testing.T) {
	// Create a new HTTP request with a method that isn't allowed
	req, err := http.NewRequest("POST", "/api/v1/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Call the TestAPI handler
	handler := http.HandlerFunc(controllers.TestAPI)
	handler.ServeHTTP(rr, req)

	// Check the response status code for Method Not Allowed
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("TestAPI returned wrong status code for POST: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}
