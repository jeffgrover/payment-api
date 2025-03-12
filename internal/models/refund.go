package models

import (
	"time"
)

// Refund represents a refund transaction in the system
type Refund struct {
	ID        string    `json:"id" gorm:"primaryKey" example:"ref_123456789" description:"Unique identifier for the refund"`
	PaymentID string    `json:"payment_id" gorm:"index" example:"pay_123456789" description:"ID of the payment being refunded"`
	Amount    int64     `json:"amount" example:"2000" description:"Amount to refund in cents"`
	Status    string    `json:"status" example:"succeeded" description:"Status of the refund (pending, succeeded, failed)"`
	Reason    string    `json:"reason,omitempty" example:"requested_by_customer" description:"Reason for the refund"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T12:00:00Z" description:"Time at which the refund was created"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T12:00:00Z" description:"Time at which the refund was last updated"`
}

// CreateRefundRequest represents the request to create a new refund
type CreateRefundRequest struct {
	PaymentID string `json:"payment_id" validate:"required" example:"pay_123456789" description:"ID of the payment to refund"`
	Amount    int64  `json:"amount" validate:"required,min=1" example:"2000" description:"Amount to refund in cents"`
	Reason    string `json:"reason,omitempty" example:"requested_by_customer" description:"Reason for the refund"`
}

// TableName overrides the table name used by GORM to `refunds`
func (Refund) TableName() string {
	return "refunds"
}
