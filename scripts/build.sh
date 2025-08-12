#!/bin/bash

# Complex build script for subs-service
set -e

echo "🔨 Complex build script for subs-service"
echo "=========================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function for colored output
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

# Function to check dependencies
check_dependencies() {
    print_status "Checking dependencies..."
    
    # Check for required commands
    local missing_deps=()
    
    if ! command -v docker &> /dev/null; then
        missing_deps+=("docker")
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        missing_deps+=("docker-compose")
    fi
    
    if ! command -v protoc &> /dev/null; then
        missing_deps+=("protoc")
    fi
    
    if ! command -v go &> /dev/null; then
        missing_deps+=("go")
    fi
    
    if ! command -v jq &> /dev/null; then
        missing_deps+=("jq")
    fi
    
    if ! command -v openssl &> /dev/null; then
        missing_deps+=("openssl")
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        print_error "Missing required dependencies:"
        for dep in "${missing_deps[@]}"; do
            echo "   - $dep"
        done
        exit 1
    fi
    
    print_success "All dependencies installed"
}

# Function to validate environment variables
validate_environment() {
    print_status "Validating environment variables..."
    
    if [ -f "scripts/validate-env.sh" ]; then
        if ./scripts/validate-env.sh; then
            print_success "Environment validation passed"
        else
            print_warning "Environment validation failed, but continuing..."
        fi
    else
        print_warning "validate-env.sh not found, skipping environment validation"
    fi
}

# Function to generate TLS certificates
generate_tls_certificates() {
    print_status "Generating TLS certificates..."
    
    if [ -f "scripts/generate-tls.sh" ]; then
        ./scripts/generate-tls.sh
        print_success "TLS certificates generated"
    else
        print_warning "generate-tls.sh not found, skipping TLS certificate generation"
    fi
}

# Function to check and update proto files
check_proto_files() {
    print_status "Checking proto files..."
    
    if [ -f "scripts/generate-proto.sh" ]; then
        ./scripts/generate-proto.sh
        print_success "Code from proto files generated"
    else
        print_warning "generate-proto.sh not found, skipping generation"
    fi
}

# Function to check and update go.mod and go.sum
update_go_dependencies() {
    print_status "Checking and updating Go dependencies..."

    if [ -f "scripts/update-dependencies.sh" ]; then
        ./scripts/update-dependencies.sh
        print_success "Go dependencies updated"
    else
        print_warning "update-dependencies.sh not found, skipping update"
    fi
}

# Function to create and validate .env file
setup_env_file() {
    print_status "Setting up .env file..."
    
    if [ -f "scripts/create-env.sh" ]; then
        ./scripts/create-env.sh
        print_success ".env file created and validated"
    else
        print_error "create-env.sh not found"
        exit 1
    fi
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

# Function for final check
final_check() {
    print_status "Final check..."
    
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

# Function to parse command line arguments
parse_arguments() {
    local docker_build=false
    local env_setup=false
    local cert_generate=false
    local purge_all=false
    local clear_volumes=false
    local proto_generate=false
    local no_cache=false
    local validate_env=false
    
    # Parse arguments
    for arg in "$@"; do
        case $arg in
            --docker)
                docker_build=true
                ;;
            --env)
                env_setup=true
                ;;
            --cert)
                cert_generate=true
                ;;
            --purge|--prune)
                purge_all=true
                ;;
            --clear)
                clear_volumes=true
                ;;
            --proto)
                proto_generate=true
                ;;
            --no-cache)
                no_cache=true
                ;;
            --validate-env)
                validate_env=true
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                print_error "Unknown argument: $arg"
                show_help
                exit 1
                ;;
        esac
    done
    
    # If no flags are specified, use automatic detection
    if [ "$docker_build" = false ] && [ "$env_setup" = false ] && [ "$cert_generate" = false ] && [ "$purge_all" = false ] && [ "$clear_volumes" = false ] && [ "$proto_generate" = false ] && [ "$validate_env" = false ]; then
        auto_detect_operations
        return
    fi
    
    # Execute requested operations
    if [ "$clear_volumes" = true ]; then
        print_status "🗄️  Clearing volumes..."
        clear_volumes_data
        echo
    fi
    
    if [ "$purge_all" = true ]; then
        print_status "🧹 Performing cleanup..."
        cleanup_all
        echo
    fi
    
    if [ "$validate_env" = true ]; then
        print_status "🔍 Validating environment..."
        validate_environment
        echo
    fi
    
    if [ "$cert_generate" = true ]; then
        print_status "🔐 Generating certificates..."
        generate_tls_certificates
        echo
    fi
    
    if [ "$env_setup" = true ]; then
        print_status "⚙️  Setting up .env file..."
        setup_env_file
        echo
    fi
    
    if [ "$proto_generate" = true ]; then
        print_status "📝 Generating proto files..."
        check_proto_files
        echo
    fi
    
    if [ "$docker_build" = true ]; then
        print_status "🐳 Building Docker images..."
        update_go_dependencies
        echo
        build_docker_images
        echo
        verify_images
        echo
    fi
    
    print_success "🎉 Operations completed successfully!"
}

# Function to automatically detect operations
auto_detect_operations() {
    print_status "🔍 Automatic detection of required operations..."
    
    local needs_env=false
    local needs_cert=false
    local needs_docker=false
    local needs_proto=false
    
    # Check .env file
    if [ ! -f ".env" ]; then
        print_status "File .env not found - needs setup"
        needs_env=true
    fi
    
    # Check certificates
    if [ ! -f "certs/server-cert.pem" ] || [ ! -f "certs/server-key.pem" ]; then
        print_status "Certificates not found - needs generation"
        needs_cert=true
    fi
    
    # Check Docker images
    local images=(
        "subs-service-auth-service"
        "subs-service-core-service"
        "subs-service-notification-service"
    )
    
    local missing_images=0
    for image in "${images[@]}"; do
        if ! docker images | grep -q "$image"; then
            missing_images=$((missing_images + 1))
        fi
    done
    
    if [ $missing_images -gt 0 ]; then
        print_status "Found $missing_images missing images - needs build"
        needs_docker=true
    fi
    
    # Check proto files
    if [ ! -f "auth-service/internal/authpb/auth_grpc.pb.go" ]; then
        print_status "Proto files not generated - needs generation"
        needs_proto=true
    fi
    
    # Execute required operations
    if [ "$needs_env" = true ]; then
        print_status "⚙️  Setting up .env file..."
        setup_env_file
        echo
    fi
    
    if [ "$needs_cert" = true ]; then
        print_status "🔐 Generating certificates..."
        generate_tls_certificates
        echo
    fi
    
    if [ "$needs_proto" = true ]; then
        print_status "📝 Generating proto files..."
        check_proto_files
        echo
    fi
    
    if [ "$needs_docker" = true ]; then
        print_status "🐳 Building Docker images..."
        update_go_dependencies
        echo
        build_docker_images
        echo
        verify_images
        echo
    fi
    
    if [ "$needs_env" = false ] && [ "$needs_cert" = false ] && [ "$needs_docker" = false ] && [ "$needs_proto" = false ]; then
        print_success "✅ All components are up to date, no additional operations required"
    else
        print_success "🎉 Automatic operations completed successfully!"
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
    echo "🔨 Universal build script for subs-service"
    echo "=================================================="
    echo
    echo "Usage:"
    echo "  ./scripts/build.sh [FLAGS]"
    echo
    echo "Flags:"
    echo "  --docker        Build Docker images"
    echo "  --env           Create/validate .env file"
    echo "  --cert          Generate TLS certificates"
    echo "  --purge/--prune Clean up unused resources (containers, networks, volumes, images, cache)"
    echo "  --clear         Clear volumes (deletes database data)"
    echo "  --proto         Generate proto files"
    echo "  --no-cache      Full rebuild without using cache"
    echo "  --validate-env  Validate environment variables"
    echo "  --help, -h      Show this help"
    echo
    echo "Flag combinations:"
    echo "  --purge --cert --env --docker     # Clean up + full rebuild"
    echo "  --clear --cert --env --docker     # Clear data + full rebuild"
    echo "  --cert --env                      # Only configuration"
    echo "  --docker                          # Only build images"
    echo "  --purge                           # Only clean up unused resources"
    echo "  --clear                           # Only clear volumes (data)"
    echo "  --validate-env                    # Only validate environment"
    echo
    echo "Automatic mode (no flags):"
    echo "  Script automatically detects required operations"
    echo
    echo "Examples:"
    echo "  ./scripts/build.sh                    # Automatic mode"
    echo "  ./scripts/build.sh --docker           # Only build images (with cache)"
    echo "  ./scripts/build.sh --docker --no-cache # Full rebuild of images"
    echo "  ./scripts/build.sh --purge --docker   # Clean up + build"
    echo "  ./scripts/build.sh --clear --docker   # Clear data + build"
    echo "  ./scripts/build.sh --cert --env       # Only configuration"
    echo "  ./scripts/build.sh --clear            # Clear database data"
    echo "  ./scripts/build.sh --validate-env     # Validate environment variables"
}

main() {
    echo
    print_status "Universal build script for subs-service..."
    echo
    
    # Check if we are in the root directory of the project
    if [ ! -f "docker-compose.yaml" ]; then
        print_error "docker-compose.yaml not found. Run the script from the root directory of the project."
        exit 1
    fi
    
    # Check dependencies
    check_dependencies
    echo
    
    # Parse arguments and execute operations
    parse_arguments "$@"
    
    echo
    echo "🚀 To start all services, run:"
    echo "   docker-compose up -d"
    echo
    echo "🔍 To check the status:"
    echo "   docker-compose ps"
    echo "   ./scripts/health-check.sh"
    echo "   ./scripts/test-api.sh"
    echo
    echo "📖 For help:"
    echo "   ./scripts/build.sh --help"
}

# Run main function
main "$@"