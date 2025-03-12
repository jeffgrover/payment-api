package db

import (
	"os"
	"testing"
	"time"

	"github.com/jeffgrover/payment-api/internal/models"
)

// Setup test database
func setupTestDB(t *testing.T) (*DB, func()) {
	// Use an in-memory SQLite database for testing
	db, err := New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Return the database and a cleanup function
	return db, func() {
		// No cleanup needed for in-memory database
	}
}

func TestNew(t *testing.T) {
	// Create a temporary file for the database
	tmpFile, err := os.CreateTemp("", "test-payments-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create a new database
	db, err := New(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	// Check that the database is not nil
	if db == nil {
		t.Fatal("Expected database to not be nil")
	}

	// Check that the database connection is not nil
	if db.DB == nil {
		t.Fatal("Expected database connection to not be nil")
	}
}

func TestCustomerOperations(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a test customer
	customer := &models.Customer{
		ID:        "cus_test123",
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test CreateCustomer
	err := db.CreateCustomer(customer)
	if err != nil {
		t.Fatalf("Failed to create customer: %v", err)
	}

	// Test GetCustomer
	retrievedCustomer, err := db.GetCustomer(customer.ID)
	if err != nil {
		t.Fatalf("Failed to get customer: %v", err)
	}

	// Check that the retrieved customer matches the original
	if retrievedCustomer.ID != customer.ID {
		t.Errorf("Expected ID to be '%s', got '%s'", customer.ID, retrievedCustomer.ID)
	}
	if retrievedCustomer.Email != customer.Email {
		t.Errorf("Expected Email to be '%s', got '%s'", customer.Email, retrievedCustomer.Email)
	}
	if retrievedCustomer.Name != customer.Name {
		t.Errorf("Expected Name to be '%s', got '%s'", customer.Name, retrievedCustomer.Name)
	}

	// Test ListCustomers
	customers, err := db.ListCustomers(10)
	if err != nil {
		t.Fatalf("Failed to list customers: %v", err)
	}

	// Check that the list contains the customer
	if len(customers) != 1 {
		t.Errorf("Expected 1 customer, got %d", len(customers))
	}
	if customers[0].ID != customer.ID {
		t.Errorf("Expected ID to be '%s', got '%s'", customer.ID, customers[0].ID)
	}
}

func TestPaymentMethodOperations(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a test customer first
	customer := &models.Customer{
		ID:        "cus_test123",
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := db.CreateCustomer(customer)
	if err != nil {
		t.Fatalf("Failed to create customer: %v", err)
	}

	// Create a test payment method
	paymentMethod := &models.PaymentMethod{
		ID:         "pm_test123",
		CustomerID: customer.ID,
		Type:       "card",
		Last4:      "4242",
		ExpMonth:   12,
		ExpYear:    2025,
		Brand:      "visa",
		CreatedAt:  time.Now(),
	}

	// Test CreatePaymentMethod
	err = db.CreatePaymentMethod(paymentMethod)
	if err != nil {
		t.Fatalf("Failed to create payment method: %v", err)
	}

	// Test GetPaymentMethod
	retrievedMethod, err := db.GetPaymentMethod(paymentMethod.ID)
	if err != nil {
		t.Fatalf("Failed to get payment method: %v", err)
	}

	// Check that the retrieved payment method matches the original
	if retrievedMethod.ID != paymentMethod.ID {
		t.Errorf("Expected ID to be '%s', got '%s'", paymentMethod.ID, retrievedMethod.ID)
	}
	if retrievedMethod.CustomerID != paymentMethod.CustomerID {
		t.Errorf("Expected CustomerID to be '%s', got '%s'", paymentMethod.CustomerID, retrievedMethod.CustomerID)
	}
	if retrievedMethod.Type != paymentMethod.Type {
		t.Errorf("Expected Type to be '%s', got '%s'", paymentMethod.Type, retrievedMethod.Type)
	}
	if retrievedMethod.Last4 != paymentMethod.Last4 {
		t.Errorf("Expected Last4 to be '%s', got '%s'", paymentMethod.Last4, retrievedMethod.Last4)
	}

	// Test GetPaymentMethodByCustomer
	retrievedMethodByCustomer, err := db.GetPaymentMethodByCustomer(paymentMethod.ID, customer.ID)
	if err != nil {
		t.Fatalf("Failed to get payment method by customer: %v", err)
	}

	// Check that the retrieved payment method matches the original
	if retrievedMethodByCustomer.ID != paymentMethod.ID {
		t.Errorf("Expected ID to be '%s', got '%s'", paymentMethod.ID, retrievedMethodByCustomer.ID)
	}
}
