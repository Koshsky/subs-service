#!/bin/bash

# Docker operations script for subs-service
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check Docker dependencies
check_docker_dependencies() {
    print_status "Checking Docker dependencies..."
    
    local missing_deps=()
    
    if ! command -v docker &> /dev/null; then
        missing_deps+=("docker")
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        missing_deps+=("docker-compose")
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        print_error "Missing required Docker dependencies:"
        for dep in "${missing_deps[@]}"; do
            echo "   - $dep"
        done
        exit 1
    fi
    
    print_success "Docker dependencies installed"
}

# Function to stop and clean up containers
cleanup_containers() {
    print_status "Cleaning up existing containers..."
    
    # Stop containers if they are running
    if docker-compose ps | grep -q "Up"; then
        print_status "Stopping running containers..."
        docker-compose down
    fi
    
    # Remove old images
    print_status "Removing old images..."
    docker-compose down --rmi all --volumes --remove-orphans 2>/dev/null || true
    
    print_success "Containers cleaned up"
}

# Function to build Docker images
build_docker_images() {
    local no_cache=${1:-false}
    
    print_status "Building Docker images..."
    
    # Define cache flag
    local cache_flag=""
    if [ "$no_cache" = true ]; then
        cache_flag="--no-cache"
        print_status "Full rebuild without cache"
    else
        print_status "Using cache for faster build"
    fi
    
    # Build images
    print_status "Building auth-service..."
    docker-compose build $cache_flag auth-service
    
    print_status "Building core-service..."
    docker-compose build $cache_flag core-service
    
    print_status "Building notification-service..."
    docker-compose build $cache_flag notification-service

    print_success "All images built"
}

# Function to check built images
verify_images() {
    print_status "Checking built images..."
    
    # Get list of all project images
    local project_images=$(docker images --format "table {{.Repository}}" | grep "subs-service" || true)
    
    if [ -z "$project_images" ]; then
        print_error "No project images found"
        exit 1
    fi
    
    print_status "Found project images:"
    echo "$project_images" | while read -r image; do
        if [ -n "$image" ]; then
            print_success "   $image"
        fi
    done
    
    # Check for main images
    local expected_images=(
        "subs-service-auth-service"
        "subs-service-core-service"
        "subs-service-notification-service"
    )
    
    local missing_images=0
    for expected_image in "${expected_images[@]}"; do
        if echo "$project_images" | grep -q "$expected_image"; then
            print_success "Image $expected_image found"
        else
            print_error "Image $expected_image not found"
            missing_images=$((missing_images + 1))
        fi
    done
    
    if [ $missing_images -gt 0 ]; then
        print_error "Found $missing_images missing images"
        exit 1
    fi
    
    print_success "All images checked"
}

# Function to start base services
start_base_services() {
    print_status "Starting base services (DB, RabbitMQ)..."
    
    # Start only base services for testing
    docker-compose up -d auth-db core-db notify-db rabbitmq
    
    # Wait a bit for startup
    sleep 5
    
    # Check status
    if docker-compose ps | grep -q "Up"; then
        print_success "Base services started"
    else
        print_error "Error starting base services"
        docker-compose logs
        exit 1
    fi
}

# Function to clean up unused resources
cleanup_all() {
    print_status "Cleaning up unused Docker resources..."
    
    # Remove only stopped containers (not affecting running ones)
    print_status "Removing stopped containers..."
    docker container prune -f
    
    # Remove unused networks
    print_status "Removing unused networks..."
    docker network prune -f
    
    # Remove unused volumes (only those not used by containers)
    print_status "Removing unused volumes..."
    docker volume prune -f
    
    # Remove unused images (dangling images)
    print_status "Removing unused images..."
    docker image prune -f
    
    # Remove unused layers (build cache)
    print_status "Cleaning build cache..."
    docker builder prune -f
    
    print_success "Cleanup completed"
}

# Function to clear volumes (deletes database data)
clear_volumes_data() {
    print_status "Clearing volumes..."
    
    # Stop containers
    print_status "Stopping containers..."
    docker-compose down 2>/dev/null || true
    
    # Remove volumes
    print_status "Removing volumes..."
    docker volume rm auth_postgres_data core_postgres_data notify_postgres_data rabbitmq_data 2>/dev/null || true
    
    # Also remove volumes that may be created automatically
    docker volume ls | grep "subs-service" | awk '{print $2}' | xargs -r docker volume rm 2>/dev/null || true
    
    print_success "Volumes cleared"
}

# Function to show help
show_help() {
    echo "🐳 Docker operations script for subs-service"
    echo "============================================="
    echo
    echo "Usage:"
    echo "  ./scripts/docker.sh [COMMAND] [OPTIONS]"
    echo
    echo "Commands:"
    echo "  build [--no-cache]    Build Docker images"
    echo "  verify                Verify built images"
    echo "  cleanup               Clean up unused resources"
    echo "  clear-volumes         Clear volumes (deletes database data)"
    echo "  start-base            Start base services (DB, RabbitMQ)"
    echo "  stop                  Stop all containers"
    echo "  logs [SERVICE]        Show logs for all services or specific service"
    echo "  status                Show container status"
    echo "  help                  Show this help"
    echo
    echo "Options:"
    echo "  --no-cache            Full rebuild without using cache"
    echo
    echo "Examples:"
    echo "  ./scripts/docker.sh build              # Build images with cache"
    echo "  ./scripts/docker.sh build --no-cache   # Full rebuild"
    echo "  ./scripts/docker.sh verify             # Check built images"
    echo "  ./scripts/docker.sh cleanup            # Clean up unused resources"
    echo "  ./scripts/docker.sh clear-volumes      # Clear database data"
    echo "  ./scripts/docker.sh start-base         # Start base services"
    echo "  ./scripts/docker.sh logs auth-service  # Show auth-service logs"
    echo "  ./scripts/docker.sh status             # Show container status"
}

# Function to show container status
show_status() {
    print_status "Container status:"
    docker-compose ps
}

# Function to show logs
show_logs() {
    local service=${1:-""}
    
    if [ -n "$service" ]; then
        print_status "Showing logs for $service:"
        docker-compose logs "$service"
    else
        print_status "Showing logs for all services:"
        docker-compose logs
    fi
}

# Function to stop containers
stop_containers() {
    print_status "Stopping all containers..."
    docker-compose down
    print_success "All containers stopped"
}

# Main function
main() {
    local command=${1:-"help"}
    local option=${2:-""}
    
    # Check if we are in the root directory of the project
    if [ ! -f "docker-compose.yaml" ]; then
        print_error "docker-compose.yaml not found. Run the script from the root directory of the project."
        exit 1
    fi
    
    # Check Docker dependencies
    check_docker_dependencies
    echo
    
    case $command in
        build)
            local no_cache=false
            if [ "$option" = "--no-cache" ]; then
                no_cache=true
            fi
            build_docker_images $no_cache
            verify_images
            ;;
        verify)
            verify_images
            ;;
        cleanup)
            cleanup_all
            ;;
        clear-volumes)
            clear_volumes_data
            ;;
        start-base)
            start_base_services
            ;;
        stop)
            stop_containers
            ;;
        logs)
            show_logs "$option"
            ;;
        status)
            show_status
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            print_error "Unknown command: $command"
            show_help
            exit 1
            ;;
    esac
    
    echo
    print_success "Docker operation completed!"
}

# Run main function
main "$@"
