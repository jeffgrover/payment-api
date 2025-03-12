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

// CustomerParams represents the parameters for retrieving a customer
type CustomerParams struct {
	ID string `path:"id" description:"Customer ID" example:"cus_123456789"`
}

// ListCustomersParams represents the parameters for listing customers
type ListCustomersParams struct {
	Limit int `query:"limit" description:"Maximum number of customers to return" default:"10" example:"10"`
}

// ListCustomersResponse represents the response for listing customers
type ListCustomersResponse struct {
	Data   []models.Customer `json:"data" description:"List of customers"`
	Status int               `json:"status" example:"200" description:"HTTP status code"`
}

// CustomerResponse wraps a customer with a status field
type CustomerResponse struct {
	*models.Customer
	Status int `json:"status" example:"200" description:"HTTP status code"`
}

// registerCustomerRoutes registers all customer-related routes
func (a *API) registerCustomerRoutes() {
	// Create a customer
	huma.Register(a.API, huma.Operation{
		OperationID: "createCustomer",
		Summary:     "Create a new customer",
		Method:      http.MethodPost,
		Path:        "/v1/customers",
		Tags:        []string{"Customers"},
	}, a.createCustomer)

	// Get a customer by ID
	huma.Register(a.API, huma.Operation{
		OperationID: "getCustomer",
		Summary:     "Get a customer by ID",
		Method:      http.MethodGet,
		Path:        "/v1/customers/{id}",
		Tags:        []string{"Customers"},
	}, a.getCustomer)

	// List customers
	huma.Register(a.API, huma.Operation{
		OperationID: "listCustomers",
		Summary:     "List customers",
		Method:      http.MethodGet,
		Path:        "/v1/customers",
		Tags:        []string{"Customers"},
	}, a.listCustomers)
}

// createCustomer creates a new customer
func (a *API) createCustomer(ctx context.Context, req *models.CreateCustomerRequest) (*CustomerResponse, error) {
	// Create a new customer
	customer := &models.Customer{
		ID:        fmt.Sprintf("cus_%d", time.Now().UnixNano()),
		Email:     req.Email,
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to database
	if err := a.DB.CreateCustomer(customer); err != nil {
		return nil, huma.Error400BadRequest("Failed to create customer", err)
	}

	return &CustomerResponse{Customer: customer, Status: 201}, nil
}

// getCustomer retrieves a customer by ID
func (a *API) getCustomer(ctx context.Context, params *CustomerParams) (*CustomerResponse, error) {
	// Get customer from database
	customer, err := a.DB.GetCustomer(params.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, huma.Error404NotFound("Customer not found", err)
		}
		return nil, huma.Error500InternalServerError("Failed to retrieve customer", err)
	}

	return &CustomerResponse{Customer: customer, Status: 200}, nil
}

// listCustomers retrieves a list of customers
func (a *API) listCustomers(ctx context.Context, params *ListCustomersParams) (*ListCustomersResponse, error) {
	// Get customers from database
	customers, err := a.DB.ListCustomers(params.Limit)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to list customers", err)
	}

	return &ListCustomersResponse{
		Data:   customers,
		Status: 200,
	}, nil
}
