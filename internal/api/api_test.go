package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jeffgrover/payment-api/internal/db"
)

// Setup test API
func setupTestAPI(t *testing.T) (*API, func()) {
	// Create an in-memory database for testing
	database, err := db.New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Create API config
	config := Config{
		Title:       "Test API",
		Version:     "1.0.0",
		Description: "Test API for unit tests",
	}

	// Create API
	api := New(database, config)

	// Return the API and a cleanup function
	return api, func() {
		// No cleanup needed for in-memory database
	}
}

func TestNew(t *testing.T) {
	// Create an in-memory database for testing
	database, err := db.New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Create API config
	config := Config{
		Title:       "Test API",
		Version:     "1.0.0",
		Description: "Test API for unit tests",
	}

	// Create API
	api := New(database, config)

	// Check that the API is not nil
	if api == nil {
		t.Fatal("Expected API to not be nil")
	}

	// Check that the router is not nil
	if api.Router == nil {
		t.Fatal("Expected Router to not be nil")
	}

	// Check that the database is not nil
	if api.DB == nil {
		t.Fatal("Expected DB to not be nil")
	}

	// Check that the API is not nil
	if api.API == nil {
		t.Fatal("Expected API to not be nil")
	}
}

func TestHealthCheck(t *testing.T) {
	api, cleanup := setupTestAPI(t)
	defer cleanup()

	// Call the health check handler
	response, err := api.healthCheck(context.Background(), &struct{}{})
	if err != nil {
		t.Fatalf("Health check failed: %v", err)
	}

	// Check the response
	if response.Status != 200 {
		t.Errorf("Expected status to be 200, got %d", response.Status)
	}
	if response.Message != "ok" {
		t.Errorf("Expected message to be 'ok', got '%s'", response.Message)
	}
}

func TestStart(t *testing.T) {
	api, cleanup := setupTestAPI(t)
	defer cleanup()

	// Create a test server
	server := httptest.NewServer(api.Router)
	defer server.Close()

	// Make a request to the health endpoint
	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
