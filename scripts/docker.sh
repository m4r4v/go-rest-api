#!/bin/bash

# Docker management script for Go REST API Framework
# This script handles Docker build, run, and management operations

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

CONTAINER_NAME="go-rest-api"
IMAGE_NAME="go-rest-api:latest"
PORT="8080"

# Function to display usage
usage() {
    echo -e "${GREEN}🐳 Docker Management for Go REST API Framework${NC}"
    echo -e "${GREEN}===============================================${NC}"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  build     Build the Docker image"
    echo "  run       Run the container (builds if needed)"
    echo "  stop      Stop and remove the container"
    echo "  restart   Stop, build, and run the container"
    echo "  logs      Show container logs"
    echo "  shell     Open a shell in the running container"
    echo "  clean     Remove container and image"
    echo "  status    Show container status"
    echo "  test      Run the demo script against the container"
    echo ""
    echo "Examples:"
    echo "  $0 build"
    echo "  $0 run"
    echo "  $0 test"
}

# Function to build Docker image
build_image() {
    echo -e "${BLUE}🏗️  Building Docker image...${NC}"
    if docker build -t $IMAGE_NAME .; then
        echo -e "${GREEN}✅ Docker image built successfully${NC}"
    else
        echo -e "${RED}❌ Docker build failed${NC}"
        exit 1
    fi
}

# Function to run container
run_container() {
    # Check if container is already running
    if docker ps -q -f name=$CONTAINER_NAME | grep -q .; then
        echo -e "${YELLOW}⚠️  Container is already running${NC}"
        echo -e "${YELLOW}Use '$0 stop' to stop it first${NC}"
        return 1
    fi
    
    # Check if image exists
    if ! docker images -q $IMAGE_NAME | grep -q .; then
        echo -e "${YELLOW}📦 Image not found, building...${NC}"
        build_image
    fi
    
    echo -e "${BLUE}🚀 Starting container...${NC}"
    docker run -d \
        --name $CONTAINER_NAME \
        -p $PORT:$PORT \
        -e PORT=$PORT \
        $IMAGE_NAME
    
    echo -e "${GREEN}✅ Container started successfully${NC}"
    echo -e "${GREEN}🌐 Server available at: http://localhost:$PORT${NC}"
    echo -e "${YELLOW}Use '$0 logs' to view logs${NC}"
    echo -e "${YELLOW}Use '$0 test' to run tests${NC}"
}

# Function to stop container
stop_container() {
    if docker ps -q -f name=$CONTAINER_NAME | grep -q .; then
        echo -e "${BLUE}🛑 Stopping container...${NC}"
        docker stop $CONTAINER_NAME
        docker rm $CONTAINER_NAME
        echo -e "${GREEN}✅ Container stopped and removed${NC}"
    else
        echo -e "${YELLOW}ℹ️  Container is not running${NC}"
    fi
}

# Function to show logs
show_logs() {
    if docker ps -q -f name=$CONTAINER_NAME | grep -q .; then
        echo -e "${BLUE}📋 Container logs:${NC}"
        docker logs -f $CONTAINER_NAME
    else
        echo -e "${RED}❌ Container is not running${NC}"
        exit 1
    fi
}

# Function to open shell
open_shell() {
    if docker ps -q -f name=$CONTAINER_NAME | grep -q .; then
        echo -e "${BLUE}🐚 Opening shell in container...${NC}"
        docker exec -it $CONTAINER_NAME /bin/sh
    else
        echo -e "${RED}❌ Container is not running${NC}"
        exit 1
    fi
}

# Function to clean up
clean_up() {
    echo -e "${BLUE}🧹 Cleaning up Docker resources...${NC}"
    
    # Stop and remove container
    if docker ps -q -f name=$CONTAINER_NAME | grep -q .; then
        docker stop $CONTAINER_NAME
        docker rm $CONTAINER_NAME
        echo -e "${GREEN}✅ Container removed${NC}"
    fi
    
    # Remove image
    if docker images -q $IMAGE_NAME | grep -q .; then
        docker rmi $IMAGE_NAME
        echo -e "${GREEN}✅ Image removed${NC}"
    fi
    
    echo -e "${GREEN}🎉 Cleanup completed${NC}"
}

# Function to show status
show_status() {
    echo -e "${BLUE}📊 Container Status:${NC}"
    echo ""
    
    if docker ps -q -f name=$CONTAINER_NAME | grep -q .; then
        echo -e "${GREEN}✅ Container is running${NC}"
        docker ps -f name=$CONTAINER_NAME --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
    else
        echo -e "${YELLOW}⚠️  Container is not running${NC}"
    fi
    
    echo ""
    echo -e "${BLUE}📦 Image Status:${NC}"
    if docker images -q $IMAGE_NAME | grep -q .; then
        echo -e "${GREEN}✅ Image exists${NC}"
        docker images $IMAGE_NAME --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}"
    else
        echo -e "${YELLOW}⚠️  Image not found${NC}"
    fi
}

# Function to run tests
run_tests() {
    if ! docker ps -q -f name=$CONTAINER_NAME | grep -q .; then
        echo -e "${RED}❌ Container is not running${NC}"
        echo -e "${YELLOW}Use '$0 run' to start the container first${NC}"
        exit 1
    fi
    
    echo -e "${BLUE}🧪 Running tests against containerized application...${NC}"
    echo -e "${YELLOW}Waiting 3 seconds for container to be ready...${NC}"
    sleep 3
    
    if [ -f "examples/framework_demo.sh" ]; then
        ./examples/framework_demo.sh
    else
        echo -e "${RED}❌ Test script not found: examples/framework_demo.sh${NC}"
        exit 1
    fi
}

# Main script logic
case "${1:-}" in
    build)
        build_image
        ;;
    run)
        run_container
        ;;
    stop)
        stop_container
        ;;
    restart)
        stop_container
        build_image
        run_container
        ;;
    logs)
        show_logs
        ;;
    shell)
        open_shell
        ;;
    clean)
        clean_up
        ;;
    status)
        show_status
        ;;
    test)
        run_tests
        ;;
    *)
        usage
        exit 1
        ;;
esac
