#!/bin/bash

# Main build script for subs-service - orchestrates specialized scripts
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

# Function to show help
show_help() {
    echo "🔨 Main build script for subs-service"
    echo "====================================="
    echo
    echo "This script orchestrates specialized scripts for different operations."
    echo
    echo "Usage:"
    echo "  ./scripts/build.sh [MODE] [OPTIONS]"
    echo
    echo "Modes:"
    echo "  setup                 Run environment setup (env, cert, proto, deps)"
    echo "  build                 Build Docker images"
    echo "  full                  Run setup + build (recommended for first time)"
    echo "  clean                 Clean up Docker resources"
    echo "  reset                 Clear data and rebuild (destructive)"
    echo "  help                  Show this help"
    echo
    echo "Options:"
    echo "  --no-cache            Full rebuild without using cache"
    echo "  --force               Force regeneration of setup files"
    echo
    echo "Examples:"
    echo "  ./scripts/build.sh full              # Complete setup and build"
    echo "  ./scripts/build.sh setup             # Only environment setup"
    echo "  ./scripts/build.sh build             # Only Docker build"
    echo "  ./scripts/build.sh build --no-cache  # Full rebuild"
    echo "  ./scripts/build.sh clean             # Clean up resources"
    echo "  ./scripts/build.sh reset             # Clear data and rebuild"
    echo
    echo "Specialized scripts:"
    echo "  ./scripts/setup.sh [COMMAND]         # Environment setup operations"
    echo "  ./scripts/docker.sh [COMMAND]        # Docker operations"
    echo "  ./scripts/health-check.sh            # Health checks"
    echo "  ./scripts/test-api.sh                # API testing"
    echo "  ./scripts/test-notification.sh       # Notification testing"
}

# Function to run setup
run_setup() {
    local force=${1:-""}
    
    print_status "Running environment setup..."
    
    if [ "$force" = "--force" ]; then
        ./scripts/setup.sh all --force
    else
        ./scripts/setup.sh all
    fi
    
    print_success "Environment setup completed"
}

# Function to run build
run_build() {
    local no_cache=${1:-""}
    
    print_status "Building Docker images..."
    
    if [ "$no_cache" = "--no-cache" ]; then
        ./scripts/docker.sh build --no-cache
    else
        ./scripts/docker.sh build
    fi
    
    print_success "Docker build completed"
}

# Function to run full build
run_full_build() {
    local no_cache=${1:-""}
    local force=${2:-""}
    
    print_status "Running full build process..."
    echo
    
    # Run setup
    run_setup "$force"
    echo
    
    # Run build
    run_build "$no_cache"
    echo
    
    print_success "Full build process completed!"
}

# Function to clean up
run_cleanup() {
    print_status "Cleaning up Docker resources..."
    ./scripts/docker.sh cleanup
    print_success "Cleanup completed"
}

# Function to reset (clear data and rebuild)
run_reset() {
    print_warning "This will clear all data and rebuild from scratch!"
    print_warning "Are you sure? (y/N)"
    read -r response
    
    if [[ "$response" =~ ^[Yy]$ ]]; then
        print_status "Running reset process..."
        
        # Clear volumes
        ./scripts/docker.sh clear-volumes
        
        # Clean up resources
        ./scripts/docker.sh cleanup
        
        # Run full build
        run_full_build "--no-cache" "--force"
        
        print_success "Reset completed!"
    else
        print_status "Reset cancelled"
    fi
}

# Function to check if we're in the right directory
check_directory() {
    if [ ! -f "docker-compose.yaml" ]; then
        print_error "docker-compose.yaml not found. Run the script from the root directory of the project."
        exit 1
    fi
}

# Function to show final status
show_final_status() {
    echo
    echo "📊 Build status:"
    echo "=================="
    
    # Check images
    echo "🐳 Docker images:"
    docker images | grep "subs-service" || echo "   Images not found"
    
    # Check containers
    echo
    echo "📦 Containers:"
    docker-compose ps
    
    # Check .env file
    echo
    echo "⚙️  Configuration:"
    if [ -f ".env" ]; then
        echo "   ✅ .env file exists"
        echo "   📋 Variables: $(grep -c '^[^#]' .env)"
    else
        echo "   ❌ .env file not found"
    fi
    
    echo
    print_success "Build completed successfully!"
}

# Main function
main() {
    local mode=${1:-"help"}
    local option1=${2:-""}
    local option2=${3:-""}
    
    echo "🔨 Main build script for subs-service"
    echo "====================================="
    echo
    
    # Check if we are in the root directory of the project
    check_directory
    
    case $mode in
        setup)
            run_setup "$option1"
            ;;
        build)
            run_build "$option1"
            ;;
        full)
            run_full_build "$option1" "$option2"
            ;;
        clean)
            run_cleanup
            ;;
        reset)
            run_reset
            ;;
        help|--help|-h)
            show_help
            exit 0
            ;;
        *)
            print_error "Unknown mode: $mode"
            show_help
            exit 1
            ;;
    esac
    
    # Show final status for build operations
    if [ "$mode" = "build" ] || [ "$mode" = "full" ]; then
        show_final_status
    fi
    
    echo
    echo "🚀 Next steps:"
    echo "   docker-compose up -d                    # Start all services"
    echo "   ./scripts/health-check.sh              # Check service health"
    echo "   ./scripts/test-api.sh                  # Test API endpoints"
    echo "   ./scripts/test-notification.sh         # Test notification integration"
    echo
    echo "📖 For more options:"
    echo "   ./scripts/build.sh help                # Show this help"
    echo "   ./scripts/setup.sh help                # Setup script help"
    echo "   ./scripts/docker.sh help               # Docker script help"
}

# Run main function
main "$@"