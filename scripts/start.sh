#!/bin/bash

# Start script for Go REST API Framework
# This script starts the server and provides helpful information

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 Starting Go REST API Framework v2.0${NC}"
echo -e "${GREEN}=======================================${NC}"

# Check if server is already running
if pgrep -f "go run cmd/server/main.go" > /dev/null || pgrep -f "./main" > /dev/null; then
    echo -e "${YELLOW}⚠️  Server is already running!${NC}"
    echo -e "${YELLOW}Use './scripts/stop.sh' to stop it first${NC}"
    exit 1
fi

# Check if port 8080 is in use
if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo -e "${RED}❌ Port 8080 is already in use${NC}"
    echo -e "${YELLOW}Use './scripts/stop.sh' to stop any running instances${NC}"
    exit 1
fi

# Build the application first
echo -e "${BLUE}📦 Building application...${NC}"
if go build -o main cmd/server/main.go; then
    echo -e "${GREEN}✅ Build successful${NC}"
else
    echo -e "${RED}❌ Build failed${NC}"
    exit 1
fi

# Start the server
echo -e "${BLUE}🌐 Starting server on port 8080...${NC}"
echo -e "${YELLOW}Press Ctrl+C to stop the server${NC}"
echo ""

# Run the server
./main
