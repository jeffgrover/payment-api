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

// RefundParams represents the parameters for retrieving a refund
type RefundParams struct {
	ID string `path:"id" description:"Refund ID" example:"ref_123456789"`
}

// ListRefundsParams represents the parameters for listing refunds
type ListRefundsParams struct {
	PaymentID string `query:"payment_id" description:"Filter by payment ID" example:"pay_123456789"`
	Limit     int    `query:"limit" description:"Maximum number of refunds to return" default:"10" example:"10"`
}

// ListRefundsResponse represents the response for listing refunds
type ListRefundsResponse struct {
	Data   []models.Refund `json:"data" description:"List of refunds"`
	Status int             `json:"status" example:"200" description:"HTTP status code"`
}

// RefundResponse wraps a refund with a status field
type RefundResponse struct {
	*models.Refund
	Status int `json:"status" example:"200" description:"HTTP status code"`
}

// registerRefundRoutes registers all refund-related routes
func (a *API) registerRefundRoutes() {
	// Create a refund
	huma.Register(a.API, huma.Operation{
		OperationID: "createRefund",
		Summary:     "Create a new refund",
		Method:      http.MethodPost,
		Path:        "/v1/refunds",
		Tags:        []string{"Refunds"},
	}, a.createRefund)

	// Get a refund by ID
	huma.Register(a.API, huma.Operation{
		OperationID: "getRefund",
		Summary:     "Get a refund by ID",
		Method:      http.MethodGet,
		Path:        "/v1/refunds/{id}",
		Tags:        []string{"Refunds"},
	}, a.getRefund)

	// List refunds
	huma.Register(a.API, huma.Operation{
		OperationID: "listRefunds",
		Summary:     "List refunds",
		Method:      http.MethodGet,
		Path:        "/v1/refunds",
		Tags:        []string{"Refunds"},
	}, a.listRefunds)
}

// createRefund creates a new refund
func (a *API) createRefund(ctx context.Context, req *models.CreateRefundRequest) (*RefundResponse, error) {
	// Verify payment exists
	payment, err := a.DB.GetPayment(req.PaymentID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, huma.Error400BadRequest("Payment not found", err)
		}
		return nil, huma.Error500InternalServerError("Failed to verify payment", err)
	}

	// Verify payment can be refunded
	if payment.Status != "succeeded" {
		return nil, huma.Error400BadRequest("Payment cannot be refunded", nil)
	}

	// Verify refund amount is valid
	if req.Amount > payment.Amount {
		return nil, huma.Error400BadRequest("Refund amount exceeds payment amount", nil)
	}

	// In a real app, you'd process the refund through the payment processor
	// This is a simplified version that always succeeds
	refund := &models.Refund{
		ID:        fmt.Sprintf("ref_%d", time.Now().UnixNano()),
		PaymentID: req.PaymentID,
		Amount:    req.Amount,
		Status:    "succeeded",
		Reason:    req.Reason,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to database
	if err := a.DB.CreateRefund(refund); err != nil {
		return nil, huma.Error500InternalServerError("Failed to create refund", err)
	}

	return &RefundResponse{Refund: refund, Status: 201}, nil
}

// getRefund retrieves a refund by ID
func (a *API) getRefund(ctx context.Context, params *RefundParams) (*RefundResponse, error) {
	// Get refund from database
	refund, err := a.DB.GetRefund(params.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, huma.Error404NotFound("Refund not found", err)
		}
		return nil, huma.Error500InternalServerError("Failed to retrieve refund", err)
	}

	return &RefundResponse{Refund: refund, Status: 200}, nil
}

// listRefunds retrieves a list of refunds
func (a *API) listRefunds(ctx context.Context, params *ListRefundsParams) (*ListRefundsResponse, error) {
	// This would typically filter by payment ID and apply limits
	// For simplicity, we'll just return an empty list for now
	return &ListRefundsResponse{
		Data:   []models.Refund{},
		Status: 200,
	}, nil
}
