package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/portilho13/neighborconnect-backend/middleware"
)

func TestCORS(t *testing.T) {
	// Create a simple handler to pass through the middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap the handler with the CORS middleware
	middleware := middleware.CORS(handler)

	// Create a request
	req, err := http.NewRequest("OPTIONS", "/api/v1/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the CORS middleware
	middleware.ServeHTTP(rr, req)

	// Check if the CORS headers are set correctly
	if allowOrigin := rr.Header().Get("Access-Control-Allow-Origin"); allowOrigin != "http://localhost:3000" {
		t.Errorf("CORS returned wrong Access-Control-Allow-Origin header: got %v want http://localhost:3000", allowOrigin)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("CORS returned wrong status code for OPTIONS request: got %v want %v", status, http.StatusOK)
	}
}
