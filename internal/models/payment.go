package models

import (
	"time"
)

// Payment represents a payment transaction in the system
type Payment struct {
	ID              string    `json:"id" gorm:"primaryKey" example:"pay_123456789" description:"Unique identifier for the payment"`
	Amount          int64     `json:"amount" example:"2000" description:"Amount in cents"`
	Currency        string    `json:"currency" example:"usd" description:"Three-letter ISO currency code"`
	CustomerID      string    `json:"customer_id" gorm:"index" example:"cus_123456789" description:"ID of the customer making the payment"`
	PaymentMethodID string    `json:"payment_method_id" example:"pm_123456789" description:"ID of the payment method used"`
	Status          string    `json:"status" example:"succeeded" description:"Status of the payment (pending, succeeded, failed)"`
	Description     string    `json:"description,omitempty" example:"Payment for order #1234" description:"Description of what the payment is for"`
	CreatedAt       time.Time `json:"created_at" example:"2023-01-01T12:00:00Z" description:"Time at which the payment was created"`
	UpdatedAt       time.Time `json:"updated_at" example:"2023-01-01T12:00:00Z" description:"Time at which the payment was last updated"`
}

// CreatePaymentRequest represents the request to create a new payment
type CreatePaymentRequest struct {
	Amount          int64  `json:"amount" validate:"required,min=1" example:"2000" description:"Amount in cents"`
	Currency        string `json:"currency" validate:"required,len=3" example:"usd" description:"Three-letter ISO currency code"`
	CustomerID      string `json:"customer_id" validate:"required" example:"cus_123456789" description:"ID of the customer making the payment"`
	PaymentMethodID string `json:"payment_method_id" validate:"required" example:"pm_123456789" description:"ID of the payment method to use"`
	Description     string `json:"description,omitempty" example:"Payment for order #1234" description:"Description of what the payment is for"`
}

// TableName overrides the table name used by GORM to `payments`
func (Payment) TableName() string {
	return "payments"
}
