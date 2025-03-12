#!/bin/bash

# Simple script to manage the payment-api application

# Display help message
function show_help {
    echo "Usage: ./run.sh [command]"
    echo ""
    echo "Commands:"
    echo "  server    Run the application without building a binary"
    echo "  build     Build the application"
    echo "  test      Run unit tests"
    echo "  e2e       Run end-to-end tests (handles setup/teardown automatically)"
    echo "  clean     Clean build artifacts"
    echo "  reset-db  Reset the database (delete payments.db)"
    echo "  dev       Run with hot reload using air (for development)"
    echo "  help      Show this help message"
    echo ""
}

# Check if a command was provided
if [ $# -eq 0 ]; then
    show_help
    exit 1
fi

# Process commands
case "$1" in
    server)
        echo "Running the application..."
        go run cmd/server/main.go
        ;;
    run) # Keep 'run' for backward compatibility
        echo "Running the application... (Note: 'run' is deprecated, please use 'server' instead)"
        go run cmd/server/main.go
        ;;
    build)
        echo "Building the application..."
        go build -o payment-api ./cmd/server
        echo "Build complete. Run with ./payment-api"
        ;;
    test)
        echo "Running unit tests..."
        go test ./...
        ;;
    e2e)
        echo "Running end-to-end tests..."
        echo "Note: This will reset the database and restart the server automatically."
        echo ""
        ./e2e_test.sh
        ;;
    clean)
        echo "Cleaning build artifacts..."
        rm -f payment-api
        echo "Clean complete"
        ;;
    reset-db)
        echo "Resetting database..."
        if [ -f payments.db ]; then
            rm payments.db
            echo "Database reset complete"
        else
            echo "No database file found"
        fi
        ;;
    dev)
        echo "Running with hot reload using air..."
        air
        ;;
    help)
        show_help
        ;;
    *)
        echo "Unknown command: $1"
        show_help
        exit 1
        ;;
esac