package db

import (
	"fmt"
	"time"

	"github.com/jeffgrover/payment-api/internal/models"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is a wrapper around gorm.DB
type DB struct {
	*gorm.DB
}

// New creates a new database connection
func New(dbPath string) (*DB, error) {
	// Configure GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// Connect to SQLite database
	db, err := gorm.Open(sqlite.Open(dbPath), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations
	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Info().Str("path", dbPath).Msg("Connected to SQLite database")
	return &DB{db}, nil
}

// migrate runs database migrations
func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Customer{},
		&models.PaymentMethod{},
		&models.Payment{},
		&models.Refund{},
	)
}

// CreateCustomer creates a new customer
func (db *DB) CreateCustomer(customer *models.Customer) error {
	// Set timestamps
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()

	// Create the customer
	return db.Create(customer).Error
}

// GetCustomer retrieves a customer by ID
func (db *DB) GetCustomer(id string) (*models.Customer, error) {
	var customer models.Customer
	if err := db.First(&customer, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

// ListCustomers retrieves a list of customers
func (db *DB) ListCustomers(limit int) ([]models.Customer, error) {
	var customers []models.Customer
	if err := db.Limit(limit).Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

// CreatePaymentMethod creates a new payment method
func (db *DB) CreatePaymentMethod(method *models.PaymentMethod) error {
	method.CreatedAt = time.Now()
	return db.Create(method).Error
}

// GetPaymentMethod retrieves a payment method by ID
func (db *DB) GetPaymentMethod(id string) (*models.PaymentMethod, error) {
	var method models.PaymentMethod
	if err := db.First(&method, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &method, nil
}

// CreatePayment creates a new payment
func (db *DB) CreatePayment(payment *models.Payment) error {
	payment.CreatedAt = time.Now()
	payment.UpdatedAt = time.Now()
	return db.Create(payment).Error
}

// GetPayment retrieves a payment by ID
func (db *DB) GetPayment(id string) (*models.Payment, error) {
	var payment models.Payment
	if err := db.First(&payment, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

// CreateRefund creates a new refund
func (db *DB) CreateRefund(refund *models.Refund) error {
	refund.CreatedAt = time.Now()
	refund.UpdatedAt = time.Now()
	return db.Create(refund).Error
}

// GetRefund retrieves a refund by ID
func (db *DB) GetRefund(id string) (*models.Refund, error) {
	var refund models.Refund
	if err := db.First(&refund, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &refund, nil
}

// GetPaymentMethodByCustomer retrieves a payment method by ID and customer ID
func (db *DB) GetPaymentMethodByCustomer(id string, customerID string) (*models.PaymentMethod, error) {
	var method models.PaymentMethod
	if err := db.First(&method, "id = ? AND customer_id = ?", id, customerID).Error; err != nil {
		return nil, err
	}
	return &method, nil
}
