package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/portilho13/neighborconnect-backend/middleware"
)

func TestLoggingMiddleware(t *testing.T) {
	// Create a simple handler to pass through the middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap the handler with the Logging middleware
	middleware := middleware.Logging(handler)

	// Create a request
	req, err := http.NewRequest("GET", "/api/v1/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the middleware
	middleware.ServeHTTP(rr, req)

	// Check if the status code is correct
	if rr.Code != http.StatusOK {
		t.Errorf("Logging middleware returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}
}
