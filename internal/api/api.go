package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jeffgrover/payment-api/internal/db"
	"github.com/rs/zerolog/log"
)

// API represents the API server
type API struct {
	Router *chi.Mux
	DB     *db.DB
	API    huma.API
}

// Config represents the API configuration
type Config struct {
	Title       string
	Version     string
	Description string
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  int    `json:"status" example:"200" description:"HTTP status code"`
	Message string `json:"message" example:"ok" description:"Status message"`
}

// New creates a new API server
func New(database *db.DB, config Config) *API {
	// Create router with middleware
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)

	// Create Huma API
	humaConfig := huma.DefaultConfig(config.Title, config.Version)
	humaConfig.Info.Description = config.Description
	humaConfig.Info.TermsOfService = "https://example.com/terms"
	humaConfig.Info.Contact = &huma.Contact{
		Name:  "API Support",
		URL:   "https://example.com/support",
		Email: "support@example.com",
	}
	humaConfig.Info.License = &huma.License{
		Name: "MIT",
		URL:  "https://opensource.org/licenses/MIT",
	}

	// Create Huma adapter
	api := humachi.New(router, humaConfig)

	// Create API server
	server := &API{
		Router: router,
		DB:     database,
		API:    api,
	}

	// Register routes
	server.registerRoutes()

	return server
}

// healthCheck is a simple health check handler
func (a *API) healthCheck(ctx context.Context, input *struct{}) (*HealthResponse, error) {
	return &HealthResponse{
		Status:  200,
		Message: "ok",
	}, nil
}

// registerRoutes registers all API routes
func (a *API) registerRoutes() {
	// Add a health check endpoint
	huma.Register(a.API, huma.Operation{
		OperationID: "healthCheck",
		Summary:     "Health check endpoint",
		Method:      http.MethodGet,
		Path:        "/health",
		Tags:        []string{"System"},
	}, a.healthCheck)

	// Register customer routes
	a.registerCustomerRoutes()

	// Register payment method routes
	a.registerPaymentMethodRoutes()

	// Register payment routes
	a.registerPaymentRoutes()

	// Register refund routes
	a.registerRefundRoutes()
}

// Start starts the API server
func (a *API) Start(addr string) error {
	log.Info().Str("addr", addr).Msg("Starting API server")
	log.Info().Str("docs", "http://localhost"+addr+"/docs").Msg("API documentation available at")
	return http.ListenAndServe(addr, a.Router)
}
