package models

import (
	"time"
)

// PaymentMethod represents a payment method in the system
type PaymentMethod struct {
	ID         string    `json:"id" gorm:"primaryKey" example:"pm_123456789" description:"Unique identifier for the payment method"`
	CustomerID string    `json:"customer_id" gorm:"index" example:"cus_123456789" description:"ID of the customer this payment method belongs to"`
	Type       string    `json:"type" example:"card" description:"Type of payment method (card or bank_account)"`
	Last4      string    `json:"last4" example:"4242" description:"Last 4 digits of the card or bank account"`
	ExpMonth   int       `json:"exp_month,omitempty" example:"12" description:"Expiration month (cards only)"`
	ExpYear    int       `json:"exp_year,omitempty" example:"2025" description:"Expiration year (cards only)"`
	Brand      string    `json:"brand,omitempty" example:"visa" description:"Card brand (cards only)"`
	CreatedAt  time.Time `json:"created_at" example:"2023-01-01T12:00:00Z" description:"Time at which the payment method was created"`
}

// CreatePaymentMethodRequest represents the request to create a new payment method
type CreatePaymentMethodRequest struct {
	CustomerID string `json:"customer_id" validate:"required" example:"cus_123456789" description:"ID of the customer"`
	Type       string `json:"type" validate:"required,oneof=card bank_account" example:"card" description:"Type of payment method"`
	// For a real implementation, you'd have additional fields like card number, exp date, etc.
	CardNumber string `json:"card_number,omitempty" validate:"omitempty,len=16" example:"4242424242424242" description:"Credit card number"`
	ExpMonth   int    `json:"exp_month,omitempty" validate:"omitempty,min=1,max=12" example:"12" description:"Expiration month"`
	ExpYear    int    `json:"exp_year,omitempty" validate:"omitempty,min=2023" example:"2025" description:"Expiration year"`
	Cvc        string `json:"cvc,omitempty" validate:"omitempty,len=3" example:"123" description:"Card security code"`
}

// TableName overrides the table name used by GORM to `payment_methods`
func (PaymentMethod) TableName() string {
	return "payment_methods"
}
