#!/bin/bash

# Stop script for Go REST API Framework
# This script stops the server gracefully

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🛑 Stopping Go REST API Framework${NC}"
echo -e "${BLUE}=================================${NC}"

# Function to stop processes
stop_processes() {
    local process_name=$1
    local description=$2
    
    if pgrep -f "$process_name" > /dev/null; then
        echo -e "${YELLOW}Stopping $description...${NC}"
        pkill -f "$process_name" || true
        sleep 2
        
        # Force kill if still running
        if pgrep -f "$process_name" > /dev/null; then
            echo -e "${YELLOW}Force stopping $description...${NC}"
            pkill -9 -f "$process_name" || true
        fi
        
        echo -e "${GREEN}✅ $description stopped${NC}"
    else
        echo -e "${YELLOW}ℹ️  $description is not running${NC}"
    fi
}

# Stop Go development server
stop_processes "go run cmd/server/main.go" "Go development server"

# Stop compiled binary
stop_processes "./main" "Go REST API server"

# Stop any process using port 8080
if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo -e "${YELLOW}Stopping processes using port 8080...${NC}"
    lsof -ti:8080 | xargs kill -9 2>/dev/null || true
    echo -e "${GREEN}✅ Port 8080 freed${NC}"
fi

echo -e "${GREEN}🎉 All processes stopped successfully${NC}"
