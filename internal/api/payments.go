package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jeffgrover/payment-api/internal/models"
	"gorm.io/gorm"
)

// PaymentParams represents the parameters for retrieving a payment
type PaymentParams struct {
	ID string `path:"id" description:"Payment ID" example:"pay_123456789"`
}

// ListPaymentsParams represents the parameters for listing payments
type ListPaymentsParams struct {
	CustomerID string `query:"customer_id" description:"Filter by customer ID" example:"cus_123456789"`
	Limit      int    `query:"limit" description:"Maximum number of payments to return" default:"10" example:"10"`
}

// ListPaymentsResponse represents the response for listing payments
type ListPaymentsResponse struct {
	Data   []models.Payment `json:"data" description:"List of payments"`
	Status int              `json:"status" example:"200" description:"HTTP status code"`
}

// registerPaymentRoutes registers all payment-related routes
func (a *API) registerPaymentRoutes() {
	// Create a payment
	huma.Register(a.API, huma.Operation{
		OperationID: "createPayment",
		Summary:     "Create a new payment",
		Method:      http.MethodPost,
		Path:        "/v1/payments",
		Tags:        []string{"Payments"},
	}, a.createPayment)

	// Get a payment by ID
	huma.Register(a.API, huma.Operation{
		OperationID: "getPayment",
		Summary:     "Get a payment by ID",
		Method:      http.MethodGet,
		Path:        "/v1/payments/{id}",
		Tags:        []string{"Payments"},
	}, a.getPayment)

	// List payments
	huma.Register(a.API, huma.Operation{
		OperationID: "listPayments",
		Summary:     "List payments",
		Method:      http.MethodGet,
		Path:        "/v1/payments",
		Tags:        []string{"Payments"},
	}, a.listPayments)
}

// PaymentResponse wraps a payment with a status field
type PaymentResponse struct {
	*models.Payment
	Status int `json:"status" example:"200" description:"HTTP status code"`
}

// createPayment creates a new payment
func (a *API) createPayment(ctx context.Context, req *models.CreatePaymentRequest) (*PaymentResponse, error) {
	// Verify customer exists
	_, err := a.DB.GetCustomer(req.CustomerID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, huma.Error400BadRequest("Customer not found", err)
		}
		return nil, huma.Error500InternalServerError("Failed to verify customer", err)
	}

	// Verify payment method exists and belongs to customer
	_, err = a.DB.GetPaymentMethodByCustomer(req.PaymentMethodID, req.CustomerID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, huma.Error400BadRequest("Payment method not found or doesn't belong to customer", err)
		}
		return nil, huma.Error500InternalServerError("Failed to verify payment method", err)
	}

	// In a real app, you'd process the payment through a payment processor
	// This is a simplified version that always succeeds
	payment := &models.Payment{
		ID:              fmt.Sprintf("pay_%d", time.Now().UnixNano()),
		Amount:          req.Amount,
		Currency:        req.Currency,
		CustomerID:      req.CustomerID,
		PaymentMethodID: req.PaymentMethodID,
		Status:          "succeeded",
		Description:     req.Description,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Save to database
	if err := a.DB.CreatePayment(payment); err != nil {
		return nil, huma.Error500InternalServerError("Failed to create payment", err)
	}

	return &PaymentResponse{Payment: payment, Status: 201}, nil
}

// getPayment retrieves a payment by ID
func (a *API) getPayment(ctx context.Context, params *PaymentParams) (*PaymentResponse, error) {
	// Get payment from database
	payment, err := a.DB.GetPayment(params.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, huma.Error404NotFound("Payment not found", err)
		}
		return nil, huma.Error500InternalServerError("Failed to retrieve payment", err)
	}

	return &PaymentResponse{Payment: payment, Status: 200}, nil
}

// listPayments retrieves a list of payments
func (a *API) listPayments(ctx context.Context, params *ListPaymentsParams) (*ListPaymentsResponse, error) {
	// This would typically filter by customer ID and apply limits
	// For simplicity, we'll just return an empty list for now
	return &ListPaymentsResponse{
		Data:   []models.Payment{},
		Status: 200,
	}, nil
}
