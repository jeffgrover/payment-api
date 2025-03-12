package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Models

type Customer struct {
	ID        string    `json:"id" example:"cus_123456789" description:"Unique identifier for the customer"`
	Email     string    `json:"email" example:"user@example.com" description:"Email address of the customer"`
	Name      string    `json:"name" example:"John Doe" description:"Customer's full name"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T12:00:00Z" description:"Time at which the customer was created"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T12:00:00Z" description:"Time at which the customer was last updated"`
}

type PaymentMethod struct {
	ID         string    `json:"id" example:"pm_123456789" description:"Unique identifier for the payment method"`
	CustomerID string    `json:"customer_id" example:"cus_123456789" description:"ID of the customer this payment method belongs to"`
	Type       string    `json:"type" example:"card" description:"Type of payment method (card or bank_account)"`
	Last4      string    `json:"last4" example:"4242" description:"Last 4 digits of the card or bank account"`
	ExpMonth   int       `json:"exp_month,omitempty" example:"12" description:"Expiration month (cards only)"`
	ExpYear    int       `json:"exp_year,omitempty" example:"2025" description:"Expiration year (cards only)"`
	Brand      string    `json:"brand,omitempty" example:"visa" description:"Card brand (cards only)"`
	CreatedAt  time.Time `json:"created_at" example:"2023-01-01T12:00:00Z" description:"Time at which the payment method was created"`
}

type Payment struct {
	ID              string    `json:"id" example:"pay_123456789" description:"Unique identifier for the payment"`
	Amount          int64     `json:"amount" example:"2000" description:"Amount in cents"`
	Currency        string    `json:"currency" example:"usd" description:"Three-letter ISO currency code"`
	CustomerID      string    `json:"customer_id" example:"cus_123456789" description:"ID of the customer making the payment"`
	PaymentMethodID string    `json:"payment_method_id" example:"pm_123456789" description:"ID of the payment method used"`
	Status          string    `json:"status" example:"succeeded" description:"Status of the payment (pending, succeeded, failed)"`
	Description     string    `json:"description,omitempty" example:"Payment for order #1234" description:"Description of what the payment is for"`
	CreatedAt       time.Time `json:"created_at" example:"2023-01-01T12:00:00Z" description:"Time at which the payment was created"`
	UpdatedAt       time.Time `json:"updated_at" example:"2023-01-01T12:00:00Z" description:"Time at which the payment was last updated"`
}

type Refund struct {
	ID        string    `json:"id" example:"ref_123456789" description:"Unique identifier for the refund"`
	PaymentID string    `json:"payment_id" example:"pay_123456789" description:"ID of the payment being refunded"`
	Amount    int64     `json:"amount" example:"2000" description:"Amount to refund in cents"`
	Status    string    `json:"status" example:"succeeded" description:"Status of the refund (pending, succeeded, failed)"`
	Reason    string    `json:"reason,omitempty" example:"requested_by_customer" description:"Reason for the refund"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T12:00:00Z" description:"Time at which the refund was created"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T12:00:00Z" description:"Time at which the refund was last updated"`
}

// Request/Response types

type CreateCustomerRequest struct {
	Email string `json:"email" validate:"required,email" example:"user@example.com" description:"Customer's email address"`
	Name  string `json:"name" validate:"required" example:"John Doe" description:"Customer's full name"`
}

type CreatePaymentMethodRequest struct {
	CustomerID string `json:"customer_id" validate:"required" example:"cus_123456789" description:"ID of the customer"`
	Type       string `json:"type" validate:"required,oneof=card bank_account" example:"card" description:"Type of payment method"`
	// For a real implementation, you'd have additional fields like card number, exp date, etc.
	CardNumber string `json:"card_number,omitempty" validate:"omitempty,len=16" example:"4242424242424242" description:"Credit card number"`
	ExpMonth   int    `json:"exp_month,omitempty" validate:"omitempty,min=1,max=12" example:"12" description:"Expiration month"`
	ExpYear    int    `json:"exp_year,omitempty" validate:"omitempty,min=2023" example:"2025" description:"Expiration year"`
	Cvc        string `json:"cvc,omitempty" validate:"omitempty,len=3" example:"123" description:"Card security code"`
}

type CreatePaymentRequest struct {
	Amount          int64  `json:"amount" validate:"required,min=1" example:"2000" description:"Amount in cents"`
	Currency        string `json:"currency" validate:"required,len=3" example:"usd" description:"Three-letter ISO currency code"`
	CustomerID      string `json:"customer_id" validate:"required" example:"cus_123456789" description:"ID of the customer making the payment"`
	PaymentMethodID string `json:"payment_method_id" validate:"required" example:"pm_123456789" description:"ID of the payment method to use"`
	Description     string `json:"description,omitempty" example:"Payment for order #1234" description:"Description of what the payment is for"`
}

type CreateRefundRequest struct {
	PaymentID string `json:"payment_id" validate:"required" example:"pay_123456789" description:"ID of the payment to refund"`
	Amount    int64  `json:"amount" validate:"required,min=1" example:"2000" description:"Amount to refund in cents"`
	Reason    string `json:"reason,omitempty" example:"requested_by_customer" description:"Reason for the refund"`
}

// Database setup

func setupDB() (*gorm.DB, error) {
	// Use SQLite instead of PostgreSQL
	return gorm.Open(sqlite.Open("payments.db"), &gorm.Config{})
}

// Database migration function
func migrateDB(db *gorm.DB) error {
	return db.AutoMigrate(&Customer{}, &PaymentMethod{}, &Payment{}, &Refund{})
}

// API handlers

func createCustomer(db *gorm.DB) func(ctx huma.Context) {
	return func(ctx huma.Context) {
		var req CreateCustomerRequest
		if err := ctx.Bind(&req); err != nil {
			ctx.Error(http.StatusBadRequest, "Invalid request", err)
			return
		}

		customer := Customer{
			ID:        fmt.Sprintf("cus_%d", time.Now().UnixNano()),
			Email:     req.Email,
			Name:      req.Name,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Save to SQLite database
		if err := db.Create(&customer).Error; err != nil {
			ctx.Error(http.StatusInternalServerError, "Failed to create customer", err)
			return
		}

		ctx.JSON(http.StatusCreated, customer)
	}
}

func getCustomer(db *gorm.DB) func(ctx huma.Context) {
	return func(ctx huma.Context) {
		id := ctx.Param("id")

		var customer Customer
		if err := db.First(&customer, "id = ?", id).Error; err != nil {
			ctx.Error(http.StatusNotFound, "Customer not found", err)
			return
		}

		ctx.JSON(http.StatusOK, customer)
	}
}

func listCustomers(db *gorm.DB) func(ctx huma.Context) {
	return func(ctx huma.Context) {
		var customers []Customer

		limit := 10 // Default limit
		if ctx.Query("limit") != "" {
			fmt.Sscanf(ctx.Query("limit"), "%d", &limit)
		}

		if err := db.Limit(limit).Find(&customers).Error; err != nil {
			ctx.Error(http.StatusInternalServerError, "Failed to list customers", err)
			return
		}

		ctx.JSON(http.StatusOK, customers)
	}
}

func createPaymentMethod(db *gorm.DB) func(ctx huma.Context) {
	return func(ctx huma.Context) {
		var req CreatePaymentMethodRequest
		if err := ctx.Bind(&req); err != nil {
			ctx.Error(http.StatusBadRequest, "Invalid request", err)
			return
		}

		// Verify customer exists
		var customer Customer
		if err := db.First(&customer, "id = ?", req.CustomerID).Error; err != nil {
			ctx.Error(http.StatusBadRequest, "Customer not found", err)
			return
		}

		// In a real app, you'd validate and process card details securely
		// This is a simplified version
		paymentMethod := PaymentMethod{
			ID:         fmt.Sprintf("pm_%d", time.Now().UnixNano()),
			CustomerID: req.CustomerID,
			Type:       req.Type,
			Last4:      req.CardNumber[len(req.CardNumber)-4:],
			ExpMonth:   req.ExpMonth,
			ExpYear:    req.ExpYear,
			Brand:      "visa", // Simplified - would be detected from card number
			CreatedAt:  time.Now(),
		}

		if err := db.Create(&paymentMethod).Error; err != nil {
			ctx.Error(http.StatusInternalServerError, "Failed to create payment method", err)
			return
		}

		ctx.JSON(http.StatusCreated, paymentMethod)
	}
}

func createPayment(db *gorm.DB) func(ctx huma.Context) {
	return func(ctx huma.Context) {
		var req CreatePaymentRequest
		if err := ctx.Bind(&req); err != nil {
			ctx.Error(http.StatusBadRequest, "Invalid request", err)
			return
		}

		// Verify customer exists
		var customer Customer
		if err := db.First(&customer, "id = ?", req.CustomerID).Error; err != nil {
			ctx.Error(http.StatusBadRequest, "Customer not found", err)
			return
		}

		// Verify payment method exists and belongs to customer
		var paymentMethod PaymentMethod
		if err := db.First(&paymentMethod, "id = ? AND customer_id = ?", req.PaymentMethodID, req.CustomerID).Error; err != nil {
			ctx.Error(http.StatusBadRequest, "Payment method not found or doesn't belong to customer", err)
			return
		}

		// In a real app, you'd process the payment through a payment processor
		// This is a simplified version that always succeeds
		payment := Payment{
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

		if err := db.Create(&payment).Error; err != nil {
			ctx.Error(http.StatusInternalServerError, "Failed to create payment", err)
			return
		}

		ctx.JSON(http.StatusCreated, payment)
	}
}

func getPayment(db *gorm.DB) func(ctx huma.Context) {
	return func(ctx huma.Context) {
		id := ctx.Param("id")

		var payment Payment
		if err := db.First(&payment, "id = ?", id).Error; err != nil {
			ctx.Error(http.StatusNotFound, "Payment not found", err)
			return
		}

		ctx.JSON(http.StatusOK, payment)
	}
}

func createRefund(db *gorm.DB) func(ctx huma.Context) {
	return func(ctx huma.Context) {
		var req CreateRefundRequest
		if err := ctx.Bind(&req); err != nil {
			ctx.Error(http.StatusBadRequest, "Invalid request", err)
			return
		}

		// Verify payment exists
		var payment Payment
		if err := db.First(&payment, "id = ?", req.PaymentID).Error; err != nil {
			ctx.Error(http.StatusBadRequest, "Payment not found", err)
			return
		}

		// Verify payment can be refunded
		if payment.Status != "succeeded" {
			ctx.Error(http.StatusBadRequest, "Payment cannot be refunded", nil)
			return
		}

		// Verify refund amount is valid
		if req.Amount > payment.Amount {
			ctx.Error(http.StatusBadRequest, "Refund amount exceeds payment amount", nil)
			return
		}

		// In a real app, you'd process the refund through the payment processor
		// This is a simplified version that always succeeds
		refund := Refund{
			ID:        fmt.Sprintf("ref_%d", time.Now().UnixNano()),
			PaymentID: req.PaymentID,
			Amount:    req.Amount,
			Status:    "succeeded",
			Reason:    req.Reason,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := db.Create(&refund).Error; err != nil {
			ctx.Error(http.StatusInternalServerError, "Failed to create refund", err)
			return
		}

		ctx.JSON(http.StatusCreated, refund)
	}
}

func main() {
	// Setup database
	db, err := setupDB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
		os.Exit(1)
	}

	// Run migrations
	if err := migrateDB(db); err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
		os.Exit(1)
	}

	// Setup router with middleware
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// Setup Huma API
	api := humachi.New(r, huma.DefaultConfig("Payments API", "1.0.0"))

	// Register API routes
	api.Register(huma.Operation{
		OperationID: "createCustomer",
		Summary:     "Create a new customer",
		Method:      http.MethodPost,
		Path:        "/v1/customers",
		RequestBody: huma.RequestBody{
			Description: "Customer information",
			Content: map[string]huma.MediaType{
				"application/json": {Schema: huma.Schema{
					Type: "object",
					Properties: map[string]huma.Schema{
						"email": {Type: "string", Format: "email"},
						"name":  {Type: "string"},
					},
					Required: []string{"email", "name"},
				}},
			},
		},
		Responses: map[string]huma.Response{
			"201": {
				Description: "Customer created successfully",
				Content: map[string]huma.MediaType{
					"application/json": {Schema: huma.Schema{Type: "object"}},
				},
			},
		},
	}, createCustomer(db))

	api.Register(huma.Operation{
		OperationID: "getCustomer",
		Summary:     "Get a customer by ID",
		Method:      http.MethodGet,
		Path:        "/v1/customers/{id}",
		Parameters: []huma.Parameter{
			{
				Name:        "id",
				In:          "path",
				Required:    true,
				Description: "Customer ID",
				Schema:      huma.Schema{Type: "string"},
			},
		},
		Responses: map[string]huma.Response{
			"200": {
				Description: "Customer details",
				Content: map[string]huma.MediaType{
					"application/json": {Schema: huma.Schema{Type: "object"}},
				},
			},
		},
	}, getCustomer(db))

	api.Register(huma.Operation{
		OperationID: "listCustomers",
		Summary:     "List customers",
		Method:      http.MethodGet,
		Path:        "/v1/customers",
		Parameters: []huma.Parameter{
			{
				Name:        "limit",
				In:          "query",
				Required:    false,
				Description: "Maximum number of customers to return",
				Schema:      huma.Schema{Type: "integer", Default: 10},
			},
		},
		Responses: map[string]huma.Response{
			"200": {
				Description: "List of customers",
				Content: map[string]huma.MediaType{
					"application/json": {Schema: huma.Schema{Type: "array"}},
				},
			},
		},
	}, listCustomers(db))

	api.Register(huma.Operation{
		OperationID: "createPaymentMethod",
		Summary:     "Create a new payment method",
		Method:      http.MethodPost,
		Path:        "/v1/payment_methods",
		RequestBody: huma.RequestBody{
			Description: "Payment method information",
			Content: map[string]huma.MediaType{
				"application/json": {Schema: huma.Schema{Type: "object"}},
			},
		},
		Responses: map[string]huma.Response{
			"201": {
				Description: "Payment method created successfully",
				Content: map[string]huma.MediaType{
					"application/json": {Schema: huma.Schema{Type: "object"}},
				},
			},
		},
		Responses: map[string]huma.Response{
			"201": {
				Description: "Payment method created successfully",
				Content: map[string]huma.MediaType{
					"application/json": {Schema: huma.Schema{Type: "object"}},
				},
			},
		},
	}, createPaymentMethod(db))

	api.Register(huma.Operation{
		OperationID: "createPayment",
		Summary:     "Create a new payment",
		Method:      http.MethodPost,
		Path:        "/v1/payments",
		RequestBody: huma.RequestBody{
			Description: "Payment information",
			Content: map[string]huma.MediaType{
				"application/json": {Schema: huma.Schema{Type: "object"}},
			},
		},
		Responses: map[string]huma.Response{
			"201": {
				Description: "Payment created successfully",
				Content: map[string]huma.MediaType{
					"application/json": {Schema: huma.Schema{Type: "object"}},
				},
			},
		},
		Responses: map[string]huma.Response{
			"201": {
				Description: "Payment created successfully",
				Content: map[string]huma.MediaType{
					"application/json": {Schema: huma.Schema{Type: "object"}},
				},
			},
		},
	}, createPayment(db))

	api.Register(huma.Operation{
		OperationID: "getPayment",
		Summary:     "Get a payment by ID",
		Method:      http.MethodGet,
		Path:        "/v1/payments/{id}",
		Parameters: []huma.Parameter{
			{
				Name:        "id",
				In:          "path",
				Required:    true,
				Description: "Payment ID",
				Schema:      huma.Schema{Type: "string"},
			},
		},
		Responses: map[string]huma.Response{
			"200": {
				Description: "Payment details",
				Content: map[string]huma.MediaType{
					"application/json": {Schema: huma.Schema{Type: "object"}},
				},
			},
		},
		Responses: map[string]huma.Response{
			"200": {
				Description: "Payment details",
				Content: map[string]huma.MediaType{
					"application/json": {Schema: huma.Schema{Type: "object"}},
				},
			},
		},
	}, getPayment(db))

	api.Register(huma.Operation{
		OperationID: "createRefund",
		Summary:     "Create a new refund",
		Method:      http.MethodPost,
		Path:        "/v1/refunds",
		RequestBody: huma.RequestBody{
			Description: "Refund information",
			Content: map[string]huma.MediaType{
				"application/json": {Schema: huma.Schema{Type: "object"}},
			},
		},
		Responses: map[string]huma.Response{
			"201": {
				Description: "Refund created successfully",
				Content: map[string]huma.MediaType{
					"application/json": {Schema: huma.Schema{Type: huma.TypeObject}},
				},
			},
		},
	}, createRefund(db))

	// Start the server
	addr := ":8080"
	log.Info().Msgf("Starting server on %s", addr)
	log.Info().Msgf("API Documentation available at http://localhost:8080/docs")
	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal().Err(err).Msg("Server failed")
		os.Exit(1)
	}
}
