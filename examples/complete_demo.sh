#!/bin/bash

# Go REST API Framework v2.0 - Complete Demo
# This script demonstrates the complete user journey and all API functionality

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

API_BASE="http://localhost:8080/v1"

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0

# Function to print test results
print_test_result() {
    local test_name="$1"
    local expected_status="$2"
    local actual_status="$3"
    local response="$4"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    if [ "$actual_status" = "$expected_status" ]; then
        echo -e "${GREEN}✅ PASS${NC}: $test_name (Status: $actual_status)"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        if [ -n "$response" ]; then
            echo -e "${BLUE}   Response: $response${NC}"
        fi
    else
        echo -e "${RED}❌ FAIL${NC}: $test_name (Expected: $expected_status, Got: $actual_status)"
        if [ -n "$response" ]; then
            echo -e "${RED}   Response: $response${NC}"
        fi
    fi
    echo
}

# Function to make HTTP requests and capture status
make_request() {
    local method="$1"
    local url="$2"
    local headers="$3"
    local data="$4"
    
    if [ -n "$data" ]; then
        curl -s -w "%{http_code}" -X "$method" "$url" $headers -d "$data"
    else
        curl -s -w "%{http_code}" -X "$method" "$url" $headers
    fi
}

echo -e "${YELLOW}🚀 Go REST API Framework v2.0 - Complete Demo${NC}"
echo -e "${YELLOW}================================================${NC}"
echo

# Test 1: Initial Setup
echo -e "${BLUE}📋 Testing Initial Setup...${NC}"
response=$(make_request "POST" "$API_BASE/setup" '-H "Content-Type: application/json"' '{"username": "admin", "email": "admin@example.com", "password": "password123"}')
status_code="${response: -3}"
response_body="${response%???}"
print_test_result "Initial Admin Setup" "200" "$status_code" "$response_body"

# Test 2: Admin Login
echo -e "${BLUE}🔐 Testing Admin Login...${NC}"
response=$(make_request "POST" "$API_BASE/auth/login" '-H "Content-Type: application/json"' '{"username": "admin", "password": "password123"}')
status_code="${response: -3}"
response_body="${response%???}"
print_test_result "Admin Login" "200" "$status_code" "$response_body"

# Extract admin token
ADMIN_TOKEN=$(echo "$response_body" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

# Test 3: Get Admin Profile
echo -e "${BLUE}👤 Testing Get Admin Profile...${NC}"
response=$(make_request "GET" "$API_BASE/auth/me" "-H \"Authorization: Bearer $ADMIN_TOKEN\"")
status_code="${response: -3}"
response_body="${response%???}"
print_test_result "Get Admin Profile" "200" "$status_code" "$response_body"

# Test 4: Create Regular User (Admin Only)
echo -e "${BLUE}👥 Testing Create User (Admin)...${NC}"
response=$(make_request "POST" "$API_BASE/admin/users" "-H \"Content-Type: application/json\" -H \"Authorization: Bearer $ADMIN_TOKEN\"" '{"username": "testuser", "email": "test@example.com", "password": "password123", "role": "user"}')
status_code="${response: -3}"
response_body="${response%???}"
print_test_result "Create User by Admin" "200" "$status_code" "$response_body"

# Test 5: List All Users (Admin Only)
echo -e "${BLUE}📋 Testing List Users (Admin)...${NC}"
response=$(make_request "GET" "$API_BASE/admin/users" "-H \"Authorization: Bearer $ADMIN_TOKEN\"")
status_code="${response: -3}"
response_body="${response%???}"
print_test_result "List All Users" "200" "$status_code" "$response_body"

# Test 6: Regular User Login
echo -e "${BLUE}🔐 Testing Regular User Login...${NC}"
response=$(make_request "POST" "$API_BASE/auth/login" '-H "Content-Type: application/json"' '{"username": "testuser", "password": "password123"}')
status_code="${response: -3}"
response_body="${response%???}"
print_test_result "Regular User Login" "200" "$status_code" "$response_body"

# Extract user token
USER_TOKEN=$(echo "$response_body" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

# Test 7: User tries to access admin endpoint (should fail)
echo -e "${BLUE}🚫 Testing User Access to Admin Endpoint (Should Fail)...${NC}"
response=$(make_request "GET" "$API_BASE/admin/users" "-H \"Authorization: Bearer $USER_TOKEN\"")
status_code="${response: -3}"
response_body="${response%???}"
print_test_result "User Access Admin Endpoint (Forbidden)" "403" "$status_code" "$response_body"

# Test 8: Create Resource (User)
echo -e "${BLUE}📦 Testing Create Resource (User)...${NC}"
response=$(make_request "POST" "$API_BASE/resources" "-H \"Content-Type: application/json\" -H \"Authorization: Bearer $USER_TOKEN\"" '{"name": "My API", "description": "A custom API endpoint", "data": {"endpoint": "/api/custom", "method": "GET"}}')
status_code="${response: -3}"
response_body="${response%???}"
print_test_result "Create Resource by User" "200" "$status_code" "$response_body"

# Extract resource ID
RESOURCE_ID=$(echo "$response_body" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)

# Test 9: List Resources
echo -e "${BLUE}📋 Testing List Resources...${NC}"
response=$(make_request "GET" "$API_BASE/resources" "-H \"Authorization: Bearer $USER_TOKEN\"")
status_code="${response: -3}"
response_body="${response%???}"
print_test_result "List Resources" "200" "$status_code" "$response_body"

# Test 10: Get Specific Resource
echo -e "${BLUE}📦 Testing Get Specific Resource...${NC}"
response=$(make_request "GET" "$API_BASE/resources/$RESOURCE_ID" "-H \"Authorization: Bearer $USER_TOKEN\"")
status_code="${response: -3}"
response_body="${response%???}"
print_test_result "Get Specific Resource" "200" "$status_code" "$response_body"

# Test 11: Update Resource (Owner)
echo -e "${BLUE}✏️ Testing Update Resource (Owner)...${NC}"
response=$(make_request "PUT" "$API_BASE/resources/$RESOURCE_ID" "-H \"Content-Type: application/json\" -H \"Authorization: Bearer $USER_TOKEN\"" '{"name": "Updated API", "description": "An updated custom API endpoint"}')
status_code="${response: -3}"
response_body="${response%???}"
print_test_result "Update Resource by Owner" "200" "$status_code" "$response_body"

# Test 12: User Update Own Profile
echo -e "${BLUE}👤 Testing User Update Own Profile...${NC}"
response=$(make_request "PUT" "$API_BASE/users/me" "-H \"Content-Type: application/json\" -H \"Authorization: Bearer $USER_TOKEN\"" '{"email": "newemail@example.com"}')
status_code="${response: -3}"
response_body="${response%???}"
print_test_result "User Update Own Profile" "200" "$status_code" "$response_body"

# Test 13: Server Status
echo -e "${BLUE}📊 Testing Server Status...${NC}"
response=$(make_request "GET" "$API_BASE/status" "-H \"Authorization: Bearer $USER_TOKEN\"")
status_code="${response: -3}"
response_body="${response%???}"
print_test_result "Server Status" "200" "$status_code" "$response_body"

# Test 14: Health Check
echo -e "${BLUE}❤️ Testing Health Check...${NC}"
response=$(make_request "GET" "$API_BASE/health" "-H \"Authorization: Bearer $USER_TOKEN\"")
status_code="${response: -3}"
response_body="${response%???}"
print_test_result "Health Check" "200" "$status_code" "$response_body"

# Test 15: Admin Delete Resource
echo -e "${BLUE}🗑️ Testing Admin Delete Resource...${NC}"
response=$(make_request "DELETE" "$API_BASE/resources/$RESOURCE_ID" "-H \"Authorization: Bearer $ADMIN_TOKEN\"")
status_code="${response: -3}"
response_body="${response%???}"
print_test_result "Admin Delete Resource" "200" "$status_code" "$response_body"

# Test 16: Try Setup Again (Should Fail)
echo -e "${BLUE}🚫 Testing Setup Again (Should Fail)...${NC}"
response=$(make_request "POST" "$API_BASE/setup" '-H "Content-Type: application/json"' '{"username": "admin2", "email": "admin2@example.com", "password": "password123"}')
status_code="${response: -3}"
response_body="${response%???}"
print_test_result "Setup Already Complete (Bad Request)" "400" "$status_code" "$response_body"

# Test 17: Invalid Login
echo -e "${BLUE}🚫 Testing Invalid Login (Should Fail)...${NC}"
response=$(make_request "POST" "$API_BASE/auth/login" '-H "Content-Type: application/json"' '{"username": "admin", "password": "wrongpassword"}')
status_code="${response: -3}"
response_body="${response%???}"
print_test_result "Invalid Login (Unauthorized)" "401" "$status_code" "$response_body"

# Test 18: No Authentication (Should Fail)
echo -e "${BLUE}🚫 Testing No Authentication (Should Fail)...${NC}"
response=$(make_request "GET" "$API_BASE/auth/me" "")
status_code="${response: -3}"
response_body="${response%???}"
print_test_result "No Authentication (Unauthorized)" "401" "$status_code" "$response_body"

# Summary
echo -e "${YELLOW}📊 Test Summary${NC}"
echo -e "${YELLOW}===============${NC}"
echo -e "Total Tests: $TOTAL_TESTS"
echo -e "Passed: ${GREEN}$PASSED_TESTS${NC}"
echo -e "Failed: ${RED}$((TOTAL_TESTS - PASSED_TESTS))${NC}"

if [ $PASSED_TESTS -eq $TOTAL_TESTS ]; then
    echo -e "\n${GREEN}🎉 All tests passed! The Go REST API Framework v2.0 is working perfectly!${NC}"
    exit 0
else
    echo -e "\n${RED}❌ Some tests failed. Please check the output above.${NC}"
    exit 1
fi
