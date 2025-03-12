package models

import (
	"testing"
	"time"
)

func TestCustomerTableName(t *testing.T) {
	customer := Customer{}
	if customer.TableName() != "customers" {
		t.Errorf("Expected table name to be 'customers', got '%s'", customer.TableName())
	}
}

func TestCustomerFields(t *testing.T) {
	// Create a test customer
	now := time.Now()
	customer := Customer{
		ID:        "cus_test123",
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Test ID field
	if customer.ID != "cus_test123" {
		t.Errorf("Expected ID to be 'cus_test123', got '%s'", customer.ID)
	}

	// Test Email field
	if customer.Email != "test@example.com" {
		t.Errorf("Expected Email to be 'test@example.com', got '%s'", customer.Email)
	}

	// Test Name field
	if customer.Name != "Test User" {
		t.Errorf("Expected Name to be 'Test User', got '%s'", customer.Name)
	}

	// Test CreatedAt field
	if !customer.CreatedAt.Equal(now) {
		t.Errorf("Expected CreatedAt to be '%v', got '%v'", now, customer.CreatedAt)
	}

	// Test UpdatedAt field
	if !customer.UpdatedAt.Equal(now) {
		t.Errorf("Expected UpdatedAt to be '%v', got '%v'", now, customer.UpdatedAt)
	}
}

func TestCreateCustomerRequest(t *testing.T) {
	// Create a test request
	req := CreateCustomerRequest{
		Email: "test@example.com",
		Name:  "Test User",
	}

	// Test Email field
	if req.Email != "test@example.com" {
		t.Errorf("Expected Email to be 'test@example.com', got '%s'", req.Email)
	}

	// Test Name field
	if req.Name != "Test User" {
		t.Errorf("Expected Name to be 'Test User', got '%s'", req.Name)
	}
}
