# Simple Payments API

A lightweight payment processing API built with Go, Huma, and SQLite, inspired by Stripe and Square but with a more focused feature set.

![Go version](https://img.shields.io/badge/Go-1.24+-00ADD8.svg)
![License](https://img.shields.io/badge/license-MIT-blue.svg)

## Features

- **Modern REST API**: Clean, RESTful endpoints following industry conventions
- **Interactive Documentation**: SwaggerUI-powered API explorer and documentation
- **SQLite Database**: Simple, file-based database with zero configuration
- **Core Payment Functionality**:
  - Customer management
  - Payment method storage
  - Payment processing
  - Refund handling
- **Developer Experience**: Easy to set up and extend

## Architecture

The API is designed with a clean, modular architecture:

```mermaid
graph TD
    A[Client] -->|HTTP Requests| B[API Server]
    B -->|Responses| A
    B -->|API Documentation| C[SwaggerUI]
    B -->|Database Operations| D[SQLite Database]
    
    subgraph "API Server Components"
        E[Huma API Framework]
        F[Request Handlers]
        G[Data Models]
        H[Validation]
    end
    
    B --- E
    E --- F
    F --- G
    F --- H
```

## API Workflow

The typical payment workflow looks like this:

```mermaid
sequenceDiagram
    participant C as Client
    participant A as API
    participant D as Database
    
    C->>A: Create Customer
    A->>D: Store Customer
    D-->>A: Confirm
    A-->>C: Customer Object
    
    C->>A: Create Payment Method
    A->>D: Store Payment Method
    D-->>A: Confirm
    A-->>C: Payment Method Object
    
    C->>A: Create Payment
    A->>A: Validate Request
    A->>D: Store Payment
    D-->>A: Confirm
    A-->>C: Payment Object
    
    Note over C,A: Optional
    C->>A: Create Refund
    A->>A: Validate Payment
    A->>D: Store Refund
    D-->>A: Confirm
    A-->>C: Refund Object
```

## Project Structure

```
payment-api/
├── cmd/
│   └── server/
│       └── main.go         # Application entry point
├── internal/
│   ├── api/
│   │   ├── api.go          # API setup and configuration
│   │   ├── customers.go    # Customer endpoints
│   │   ├── payments.go     # Payment endpoints
│   │   ├── methods.go      # Payment method endpoints
│   │   └── refunds.go      # Refund endpoints
│   ├── models/
│   │   ├── customer.go     # Customer model
│   │   ├── payment.go      # Payment model
│   │   ├── method.go       # Payment method model
│   │   └── refund.go       # Refund model
│   └── db/
│       └── db.go           # Database setup and operations
├── payments.db             # SQLite database file (created at runtime)
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
└── README.md               # This documentation
```

## Prerequisites

- Go 1.24 or higher
- Git (for version control)

## Setup Instructions

1. **Clone the Repository**

```bash
git clone https://github.com/jeffgrover/payment-api.git
cd payment-api
```

2. **Install Dependencies**

```bash
go mod tidy
```

3. **Run the Application**

```bash
go run cmd/server/main.go
```

The server will start on port 8080, and you'll see output like:

```
INF Starting Payments API server
INF Starting API server addr=:8080
INF API documentation available at docs=http://localhost:8080/docs
```

4. **Browse API Documentation**

Open your browser and navigate to http://localhost:8080/docs to explore the API using the SwaggerUI interface.

## API Endpoints

### Customers
- `POST /v1/customers` - Create a customer
- `GET /v1/customers/{id}` - Retrieve a customer
- `GET /v1/customers` - List customers

### Payment Methods
- `POST /v1/payment_methods` - Create a payment method
- `GET /v1/payment_methods/{id}` - Retrieve a payment method
- `GET /v1/payment_methods` - List payment methods

### Payments
- `POST /v1/payments` - Create a payment
- `GET /v1/payments/{id}` - Retrieve a payment
- `GET /v1/payments` - List payments

### Refunds
- `POST /v1/refunds` - Create a refund
- `GET /v1/refunds/{id}` - Retrieve a refund
- `GET /v1/refunds` - List refunds

## Example Usage

### Create a Customer

```bash
curl -X POST http://localhost:8080/v1/customers \
  -H "Content-Type: application/json" \
  -d '{"email":"customer@example.com","name":"John Doe"}'
```

### Create a Payment Method

```bash
curl -X POST http://localhost:8080/v1/payment_methods \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id":"cus_1234567890",
    "type":"card",
    "card_number":"4242424242424242",
    "exp_month":12,
    "exp_year":2025,
    "cvc":"123"
  }'
```

### Create a Payment

```bash
curl -X POST http://localhost:8080/v1/payments \
  -H "Content-Type: application/json" \
  -d '{
    "amount":2000,
    "currency":"usd",
    "customer_id":"cus_1234567890",
    "payment_method_id":"pm_1234567890",
    "description":"Payment for order #1234"
  }'
```

## Development

### Auto-Reloading During Development

For a better development experience, you can use [Air](https://github.com/cosmtrek/air) for hot reloading:

```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run with Air
air
```

### Building for Production

```bash
go build -o payment-api ./cmd/server
```

### Running Tests

```bash
go test ./...
```

## Future Enhancements

- Add authentication (JWT or API Keys)
- Implement rate limiting
- Add webhook notifications
- Support for additional payment methods
- Dockerize the application

## License

This project is licensed under the MIT License - see the LICENSE file for details.
