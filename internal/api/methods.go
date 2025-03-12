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

// PaymentMethodParams represents the parameters for retrieving a payment method
type PaymentMethodParams struct {
	ID string `path:"id" description:"Payment Method ID" example:"pm_123456789"`
}

// ListPaymentMethodsParams represents the parameters for listing payment methods
type ListPaymentMethodsParams struct {
	CustomerID string `query:"customer_id" description:"Filter by customer ID" example:"cus_123456789"`
	Limit      int    `query:"limit" description:"Maximum number of payment methods to return" default:"10" example:"10"`
}

// ListPaymentMethodsResponse represents the response for listing payment methods
type ListPaymentMethodsResponse struct {
	Data   []models.PaymentMethod `json:"data" description:"List of payment methods"`
	Status int                    `json:"status" example:"200" description:"HTTP status code"`
}

// PaymentMethodResponse wraps a payment method with a status field
type PaymentMethodResponse struct {
	*models.PaymentMethod
	Status int `json:"status" example:"200" description:"HTTP status code"`
}

// registerPaymentMethodRoutes registers all payment method-related routes
func (a *API) registerPaymentMethodRoutes() {
	// Create a payment method
	huma.Register(a.API, huma.Operation{
		OperationID: "createPaymentMethod",
		Summary:     "Create a new payment method",
		Method:      http.MethodPost,
		Path:        "/v1/payment_methods",
		Tags:        []string{"Payment Methods"},
	}, a.createPaymentMethod)

	// Get a payment method by ID
	huma.Register(a.API, huma.Operation{
		OperationID: "getPaymentMethod",
		Summary:     "Get a payment method by ID",
		Method:      http.MethodGet,
		Path:        "/v1/payment_methods/{id}",
		Tags:        []string{"Payment Methods"},
	}, a.getPaymentMethod)

	// List payment methods
	huma.Register(a.API, huma.Operation{
		OperationID: "listPaymentMethods",
		Summary:     "List payment methods",
		Method:      http.MethodGet,
		Path:        "/v1/payment_methods",
		Tags:        []string{"Payment Methods"},
	}, a.listPaymentMethods)
}

// createPaymentMethod creates a new payment method
func (a *API) createPaymentMethod(ctx context.Context, req *models.CreatePaymentMethodRequest) (*PaymentMethodResponse, error) {
	// Verify customer exists
	_, err := a.DB.GetCustomer(req.CustomerID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, huma.Error400BadRequest("Customer not found", err)
		}
		return nil, huma.Error500InternalServerError("Failed to verify customer", err)
	}

	// In a real app, you'd validate and process card details securely
	// This is a simplified version
	paymentMethod := &models.PaymentMethod{
		ID:         fmt.Sprintf("pm_%d", time.Now().UnixNano()),
		CustomerID: req.CustomerID,
		Type:       req.Type,
		Last4:      req.CardNumber[len(req.CardNumber)-4:],
		ExpMonth:   req.ExpMonth,
		ExpYear:    req.ExpYear,
		Brand:      "visa", // Simplified - would be detected from card number
		CreatedAt:  time.Now(),
	}

	// Save to database
	if err := a.DB.CreatePaymentMethod(paymentMethod); err != nil {
		return nil, huma.Error500InternalServerError("Failed to create payment method", err)
	}

	return &PaymentMethodResponse{PaymentMethod: paymentMethod, Status: 201}, nil
}

// getPaymentMethod retrieves a payment method by ID
func (a *API) getPaymentMethod(ctx context.Context, params *PaymentMethodParams) (*PaymentMethodResponse, error) {
	// Get payment method from database
	method, err := a.DB.GetPaymentMethod(params.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, huma.Error404NotFound("Payment method not found", err)
		}
		return nil, huma.Error500InternalServerError("Failed to retrieve payment method", err)
	}

	return &PaymentMethodResponse{PaymentMethod: method, Status: 200}, nil
}

// listPaymentMethods retrieves a list of payment methods
func (a *API) listPaymentMethods(ctx context.Context, params *ListPaymentMethodsParams) (*ListPaymentMethodsResponse, error) {
	// This would typically filter by customer ID and apply limits
	// For simplicity, we'll just return an empty list for now
	return &ListPaymentMethodsResponse{
		Data:   []models.PaymentMethod{},
		Status: 200,
	}, nil
}
