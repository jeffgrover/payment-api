package models

import (
	"time"
)

// Customer represents a customer in the payment system
type Customer struct {
	ID        string    `json:"id" gorm:"primaryKey" example:"cus_123456789" description:"Unique identifier for the customer"`
	Email     string    `json:"email" example:"user@example.com" description:"Email address of the customer"`
	Name      string    `json:"name" example:"John Doe" description:"Customer's full name"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T12:00:00Z" description:"Time at which the customer was created"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T12:00:00Z" description:"Time at which the customer was last updated"`
}

// CreateCustomerRequest represents the request to create a new customer
type CreateCustomerRequest struct {
	Email string `json:"email" validate:"required,email" example:"user@example.com" description:"Customer's email address"`
	Name  string `json:"name" validate:"required" example:"John Doe" description:"Customer's full name"`
}

// TableName overrides the table name used by GORM to `customers`
func (Customer) TableName() string {
	return "customers"
}
