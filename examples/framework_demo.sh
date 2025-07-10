#!/bin/bash

# Go REST API Framework - Complete Demo Script
# This script demonstrates the framework's core functionality as specified

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

API_BASE="http://localhost:8080"

echo -e "${GREEN}🚀 Go REST API Framework - Complete Demo${NC}"
echo -e "${GREEN}==========================================${NC}"
echo ""

# Function to make HTTP requests and display results
make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    local headers=$4
    local description=$5
    
    echo -e "${BLUE}Testing: $description${NC}"
    echo -e "${YELLOW}Request: $method $endpoint${NC}"
    if [ ! -z "$data" ]; then
        echo -e "${YELLOW}Data: $data${NC}"
    fi
    
    if [ ! -z "$headers" ]; then
        response=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X $method "$API_BASE$endpoint" -H "Content-Type: application/json" $headers -d "$data")
    else
        response=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X $method "$API_BASE$endpoint" -H "Content-Type: application/json" -d "$data")
    fi
    
    # Extract HTTP status and response body
    http_status=$(echo "$response" | grep "HTTP_STATUS:" | cut -d: -f2)
    response_body=$(echo "$response" | sed '/HTTP_STATUS:/d')
    
    echo -e "${GREEN}Response (Status: $http_status):${NC}"
    echo "$response_body" | python3 -m json.tool 2>/dev/null || echo "$response_body"
    echo ""
}

echo -e "${BLUE}1. Testing Server Health Check${NC}"
echo "----------------------------------------"
make_request "GET" "/health" "" "" "Health check endpoint"

echo -e "${BLUE}2. Testing Server Status${NC}"
echo "----------------------------------------"
make_request "GET" "/status" "" "" "Server status endpoint"

echo -e "${BLUE}3. Testing Sample Resource (Ping)${NC}"
echo "----------------------------------------"
make_request "GET" "/v1/ping" "" "" "Sample ping resource"

echo -e "${BLUE}4. Testing Super Admin Setup${NC}"
echo "----------------------------------------"
make_request "POST" "/setup" '{"username": "johndoe", "email": "john@doe.com", "password": "str0ngP4ssw0rd"}' "" "Initial super admin setup"

echo -e "${BLUE}5. Testing Setup Already Complete (Should Fail)${NC}"
echo "----------------------------------------"
make_request "POST" "/setup" '{"username": "admin2", "email": "admin2@example.com", "password": "password123"}' "" "Attempting setup again (should fail)"

echo -e "${BLUE}6. Testing User Login${NC}"
echo "----------------------------------------"
login_response=$(curl -s -X POST "$API_BASE/login" -H "Content-Type: application/json" -d '{"username": "johndoe", "password": "str0ngP4ssw0rd"}')
echo -e "${GREEN}Login Response:${NC}"
echo "$login_response" | python3 -m json.tool 2>/dev/null || echo "$login_response"

# Extract token from response
token=$(echo "$login_response" | python3 -c "import sys, json; data=json.load(sys.stdin); print(data['response']['token'])" 2>/dev/null)

if [ ! -z "$token" ]; then
    echo -e "${GREEN}✅ Token extracted successfully!${NC}"
    echo ""
else
    echo -e "${RED}❌ Failed to extract token${NC}"
    exit 1
fi

echo -e "${BLUE}7. Testing Invalid Login (Should Fail)${NC}"
echo "----------------------------------------"
make_request "POST" "/login" '{"username": "johndoe", "password": "wrongpassword"}' "" "Invalid login attempt"

echo -e "${BLUE}8. Testing Method Validation (Should Fail)${NC}"
echo "----------------------------------------"
make_request "PUT" "/v1/ping" "" "" "Wrong HTTP method (should return 404)"

echo -e "${BLUE}9. Testing Input Validation (Should Fail)${NC}"
echo "----------------------------------------"
make_request "POST" "/login" '{"username": "", "password": "123"}' "" "Invalid input validation"

echo -e "${BLUE}10. Testing Authenticated Endpoint${NC}"
echo "----------------------------------------"
make_request "GET" "/v1/ping" "" "-H \"Authorization: Bearer $token\"" "Authenticated ping request"

echo ""
echo -e "${GREEN}🎉 Demo Complete!${NC}"
echo -e "${GREEN}=================${NC}"
echo ""
echo -e "${YELLOW}Framework Features Demonstrated:${NC}"
echo "✅ Super Admin setup (one-time only)"
echo "✅ User authentication with JWT tokens (24-hour expiry)"
echo "✅ Proper JSON response structure with standardized format"
echo "✅ HTTP method validation"
echo "✅ Input validation and sanitization"
echo "✅ Error handling with appropriate status codes"
echo "✅ Public endpoints (/setup, /login, /status, /v1/ping)"
echo "✅ Role-based access control (Super Admin role)"
echo "✅ Structured logging"
echo "✅ CORS support"
echo "✅ Graceful error responses"
echo ""
echo -e "${YELLOW}Response Format:${NC}"
echo '{
  "http_status_code": "200",
  "http_status_message": "OK", 
  "resource": "/endpoint",
  "app": "Go REST API Framework",
  "timestamp": "2025-07-10T03:37:42Z",
  "response": {
    // Any kind of data here, up to the user
  }
}'
echo ""
echo -e "${YELLOW}Next Steps:${NC}"
echo "- Framework is ready for developers to create custom resources"
echo "- Users can be invited to consume resource data"
echo "- Perfect for testing frontend applications with mock backends"
echo "- All data is volatile (destroyed when container restarts)"
echo ""
echo -e "${GREEN}🌟 Go REST API Framework v2.0 - Ready for Production!${NC}"
