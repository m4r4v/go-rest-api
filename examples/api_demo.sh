#!/bin/bash

# Go REST API Framework Demo Script
# This script demonstrates the main features of the API

echo "=========================================="
echo "  Go REST API Framework Demo"
echo "=========================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

API_BASE="http://localhost:8080"

echo -e "${BLUE}Starting API Demo...${NC}"
echo ""

# Function to make HTTP requests and display results
make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    local headers=$4
    
    echo -e "${YELLOW}Request: $method $endpoint${NC}"
    if [ ! -z "$data" ]; then
        echo -e "${YELLOW}Data: $data${NC}"
    fi
    
    if [ ! -z "$headers" ]; then
        response=$(curl -s -X $method "$API_BASE$endpoint" -H "Content-Type: application/json" $headers -d "$data")
    else
        response=$(curl -s -X $method "$API_BASE$endpoint" -H "Content-Type: application/json" -d "$data")
    fi
    
    echo -e "${GREEN}Response:${NC}"
    echo "$response" | python3 -m json.tool 2>/dev/null || echo "$response"
    echo ""
}

echo -e "${BLUE}1. Testing Initial Setup${NC}"
echo "----------------------------------------"
make_request "POST" "/setup" '{"username": "admin", "email": "admin@example.com", "password": "password123"}'

echo -e "${BLUE}2. Testing User Login${NC}"
echo "----------------------------------------"
login_response=$(curl -s -X POST "$API_BASE/login" -H "Content-Type: application/json" -d '{"username": "admin", "password": "password123"}')
echo -e "${GREEN}Login Response:${NC}"
echo "$login_response" | python3 -m json.tool 2>/dev/null || echo "$login_response"

# Extract token from response
token=$(echo "$login_response" | python3 -c "import sys, json; print(json.load(sys.stdin)['response']['token'])" 2>/dev/null)

if [ ! -z "$token" ]; then
    echo -e "${GREEN}Token extracted successfully!${NC}"
    echo ""
    
    echo -e "${BLUE}3. Testing Protected Endpoint (Health Check)${NC}"
    echo "----------------------------------------"
    make_request "GET" "/health" "" "-H \"Authorization: Bearer $token\""
    
    echo -e "${BLUE}4. Testing User Profile${NC}"
    echo "----------------------------------------"
    make_request "GET" "/v1/users/me" "" "-H \"Authorization: Bearer $token\""
    
    echo -e "${BLUE}5. Testing Ping Endpoint${NC}"
    echo "----------------------------------------"
    make_request "GET" "/v1/ping" "" "-H \"Authorization: Bearer $token\""
    
    echo -e "${BLUE}6. Testing Admin Endpoint${NC}"
    echo "----------------------------------------"
    make_request "GET" "/v1/admin/users" "" "-H \"Authorization: Bearer $token\""
    
    echo -e "${BLUE}7. Testing User Update${NC}"
    echo "----------------------------------------"
    make_request "PUT" "/v1/users/me" '{"email": "admin_updated@example.com"}' "-H \"Authorization: Bearer $token\""
    
else
    echo -e "${RED}Failed to extract token from login response${NC}"
fi

echo -e "${BLUE}8. Testing Error Handling (No Auth)${NC}"
echo "----------------------------------------"
make_request "GET" "/health" ""

echo -e "${BLUE}9. Testing Validation Error${NC}"
echo "----------------------------------------"
make_request "POST" "/login" '{"username": "", "password": "123"}'

echo -e "${BLUE}10. Testing Invalid Endpoint${NC}"
echo "----------------------------------------"
make_request "GET" "/nonexistent" ""

echo ""
echo -e "${GREEN}=========================================="
echo "  Demo Complete!"
echo "==========================================${NC}"
echo ""
echo -e "${YELLOW}Key Features Demonstrated:${NC}"
echo "- User registration and authentication"
echo "- JWT token generation and validation"
echo "- Role-based access control"
echo "- Input validation with detailed errors"
echo "- Structured error responses"
echo "- Request logging and middleware"
echo "- Graceful error handling"
echo ""
echo -e "${YELLOW}To run this demo:${NC}"
echo "1. Start the server: go run cmd/server/main.go"
echo "2. Run this script: bash examples/api_demo.sh"
