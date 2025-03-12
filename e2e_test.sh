#!/bin/bash

# End-to-end test script for the Payment API
# This script tests a typical user flow:
# 1. Create a customer
# 2. Create a payment method
# 3. Create a payment
# 4. Create a refund
# 5. Verify each step with GET requests

# Set variables
API_URL="http://localhost:8080"
CUSTOMER_ID=""
PAYMENT_METHOD_ID=""
PAYMENT_ID=""
REFUND_ID=""

# Generate unique test data with timestamp
TIMESTAMP=$(date +%s)
TEST_EMAIL="e2e_test_${TIMESTAMP}@example.com"
TEST_NAME="E2E Test User ${TIMESTAMP}"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print section headers
print_header() {
    echo -e "\n${BLUE}==== $1 ====${NC}\n"
}

# Setup function - reset the database
setup() {
    print_header "Setup - Resetting Database"
    
    echo "Stopping the server if it's running..."
    pkill -f "go run cmd/server/main.go" || true
    
    echo "Removing the database file..."
    rm -f payments.db
    
    echo "Starting the server again..."
    # Run the server in a separate terminal window for better visibility
    osascript -e 'tell app "Terminal" to do script "cd \"'$(pwd)'\" && ./run.sh server"' > /dev/null 2>&1
    
    # Wait for the server to start
    echo "Waiting for the server to start..."
    sleep 5
    
    # Check if the server is running
    response=$(curl -s -o /dev/null -w "%{http_code}" $API_URL/health)
    if [ "$response" == "200" ]; then
        echo -e "${GREEN}Server is running${NC}"
    else
        echo -e "${RED}Failed to start the server. Response code: $response${NC}"
        echo -e "${RED}Please start the server manually with './run.sh server' in another terminal and try again.${NC}"
        exit 1
    fi
}

# Teardown function - clean up after the test
teardown() {
    print_header "Teardown - Cleaning Up"
    
    echo "Note: The server is still running in a separate terminal."
    echo "You can close it manually when you're done."
    
    echo -e "${GREEN}Cleanup complete${NC}"
}

# Function to create a customer
create_customer() {
    print_header "Creating a customer"
    
    echo "Sending request to create a customer with email '$TEST_EMAIL' and name '$TEST_NAME'"
    
    # Use -i flag to include headers in the output
    response_headers=$(curl -i -s -X POST $API_URL/v1/customers \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$TEST_EMAIL\",\"name\":\"$TEST_NAME\"}")
    
    echo "Response headers: $response_headers"
    
    # Extract customer ID from the Id header
    CUSTOMER_ID=$(echo "$response_headers" | grep -i "^Id:" | awk '{print $2}' | tr -d '\r')
    
    if [ -z "$CUSTOMER_ID" ]; then
        echo -e "${RED}Failed to create customer${NC}"
        echo -e "${RED}Please check that the server is running and accessible at $API_URL${NC}"
        exit 1
    else
        echo -e "${GREEN}Successfully created customer with ID: $CUSTOMER_ID${NC}"
    fi
}

# Function to verify customer creation
verify_customer() {
    print_header "Verifying customer creation"
    
    echo "Sending request to get customer with ID: $CUSTOMER_ID"
    
    response=$(curl -s -X GET $API_URL/v1/customers/$CUSTOMER_ID)
    
    echo "Response: $response"
    
    # Check if the response contains the customer ID
    if echo $response | grep -q "$CUSTOMER_ID"; then
        echo -e "${GREEN}Successfully verified customer${NC}"
    else
        echo -e "${RED}Failed to verify customer${NC}"
        exit 1
    fi
    
    # Also test listing customers
    echo -e "\nListing all customers:"
    
    response=$(curl -s -X GET $API_URL/v1/customers)
    
    echo "Response: $response"
    
    # Check if the response contains the customer ID
    if echo $response | grep -q "$CUSTOMER_ID"; then
        echo -e "${GREEN}Successfully listed customers${NC}"
    else
        echo -e "${RED}Failed to list customers${NC}"
        exit 1
    fi
}

# Function to create a payment method
create_payment_method() {
    print_header "Creating a payment method"
    
    echo "Sending request to create a payment method for customer ID: $CUSTOMER_ID"
    
    # Use -i flag to include headers in the output
    response_headers=$(curl -i -s -X POST $API_URL/v1/payment_methods \
        -H "Content-Type: application/json" \
        -d "{
            \"customer_id\":\"$CUSTOMER_ID\",
            \"type\":\"card\",
            \"card_number\":\"4242424242424242\",
            \"exp_month\":12,
            \"exp_year\":2025,
            \"cvc\":\"123\"
        }")
    
    echo "Response headers: $response_headers"
    
    # Extract payment method ID from the Id header
    PAYMENT_METHOD_ID=$(echo "$response_headers" | grep -i "^Id:" | awk '{print $2}' | tr -d '\r')
    
    if [ -z "$PAYMENT_METHOD_ID" ]; then
        echo -e "${RED}Failed to create payment method${NC}"
        exit 1
    else
        echo -e "${GREEN}Successfully created payment method with ID: $PAYMENT_METHOD_ID${NC}"
    fi
}

# Function to verify payment method creation
verify_payment_method() {
    print_header "Verifying payment method creation"
    
    echo "Sending request to get payment method with ID: $PAYMENT_METHOD_ID"
    
    response=$(curl -s -X GET $API_URL/v1/payment_methods/$PAYMENT_METHOD_ID)
    
    echo "Response: $response"
    
    # Check if the response contains the payment method ID
    if echo $response | grep -q "$PAYMENT_METHOD_ID"; then
        echo -e "${GREEN}Successfully verified payment method${NC}"
    else
        echo -e "${RED}Failed to verify payment method${NC}"
        exit 1
    fi
}

# Function to create a payment
create_payment() {
    print_header "Creating a payment"
    
    echo "Sending request to create a payment for customer ID: $CUSTOMER_ID with payment method ID: $PAYMENT_METHOD_ID"
    
    # Use -i flag to include headers in the output
    response_headers=$(curl -i -s -X POST $API_URL/v1/payments \
        -H "Content-Type: application/json" \
        -d "{
            \"amount\":2000,
            \"currency\":\"usd\",
            \"customer_id\":\"$CUSTOMER_ID\",
            \"payment_method_id\":\"$PAYMENT_METHOD_ID\",
            \"description\":\"Payment for e2e test\"
        }")
    
    echo "Response headers: $response_headers"
    
    # Extract payment ID from the Id header
    PAYMENT_ID=$(echo "$response_headers" | grep -i "^Id:" | awk '{print $2}' | tr -d '\r')
    
    if [ -z "$PAYMENT_ID" ]; then
        echo -e "${RED}Failed to create payment${NC}"
        exit 1
    else
        echo -e "${GREEN}Successfully created payment with ID: $PAYMENT_ID${NC}"
    fi
}

# Function to verify payment creation
verify_payment() {
    print_header "Verifying payment creation"
    
    echo "Sending request to get payment with ID: $PAYMENT_ID"
    
    response=$(curl -s -X GET $API_URL/v1/payments/$PAYMENT_ID)
    
    echo "Response: $response"
    
    # Check if the response contains the payment ID
    if echo $response | grep -q "$PAYMENT_ID"; then
        echo -e "${GREEN}Successfully verified payment${NC}"
    else
        echo -e "${RED}Failed to verify payment${NC}"
        exit 1
    fi
    
    # Also test listing payments
    echo -e "\nListing all payments:"
    
    response=$(curl -s -X GET $API_URL/v1/payments)
    
    echo "Response: $response"
    
    # Check if the response contains the payment ID
    if echo $response | grep -q "$PAYMENT_ID"; then
        echo -e "${GREEN}Successfully listed payments${NC}"
    else
        echo -e "${RED}Failed to list payments${NC}"
        exit 1
    fi
}

# Function to create a refund
create_refund() {
    print_header "Creating a refund"
    
    echo "Sending request to create a refund for payment ID: $PAYMENT_ID"
    
    # Use -i flag to include headers in the output
    response_headers=$(curl -i -s -X POST $API_URL/v1/refunds \
        -H "Content-Type: application/json" \
        -d "{
            \"payment_id\":\"$PAYMENT_ID\",
            \"amount\":2000,
            \"reason\":\"Customer requested refund\"
        }")
    
    echo "Response headers: $response_headers"
    
    # Extract refund ID from the Id header
    REFUND_ID=$(echo "$response_headers" | grep -i "^Id:" | awk '{print $2}' | tr -d '\r')
    
    if [ -z "$REFUND_ID" ]; then
        echo -e "${RED}Failed to create refund${NC}"
        exit 1
    else
        echo -e "${GREEN}Successfully created refund with ID: $REFUND_ID${NC}"
    fi
}

# Function to verify refund creation
verify_refund() {
    print_header "Verifying refund creation"
    
    echo "Sending request to get refund with ID: $REFUND_ID"
    
    response=$(curl -s -X GET $API_URL/v1/refunds/$REFUND_ID)
    
    echo "Response: $response"
    
    # Check if the response contains the refund ID
    if echo $response | grep -q "$REFUND_ID"; then
        echo -e "${GREEN}Successfully verified refund${NC}"
    else
        echo -e "${RED}Failed to verify refund${NC}"
        exit 1
    fi
    
    # Also test listing refunds
    echo -e "\nListing all refunds:"
    
    response=$(curl -s -X GET $API_URL/v1/refunds)
    
    echo "Response: $response"
    
    # Check if the response contains the refund ID
    if echo $response | grep -q "$REFUND_ID"; then
        echo -e "${GREEN}Successfully listed refunds${NC}"
    else
        echo -e "${RED}Failed to list refunds${NC}"
        exit 1
    fi
}

# Main function to run the tests
run_tests() {
    print_header "Starting E2E Tests"
    
    create_customer
    verify_customer
    create_payment_method
    verify_payment_method
    create_payment
    verify_payment
    create_refund
    verify_refund
    
    print_header "E2E Tests Completed Successfully"
    echo -e "${GREEN}All tests passed!${NC}"
}

# Run the entire test suite with setup and teardown
print_header "E2E Test Suite"
echo -e "This test will reset the database and restart the server."
echo -e "Press Enter to continue or Ctrl+C to cancel..."
read

# Run setup, tests, and teardown
setup
run_tests
teardown