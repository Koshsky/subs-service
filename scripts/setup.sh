#!/bin/bash

# Environment setup script for subs-service
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

# Function to check dependencies
check_dependencies() {
    print_status "Checking dependencies..."
    
    local missing_deps=()
    
    if ! command -v go &> /dev/null; then
        missing_deps+=("go")
    fi
    
    if ! command -v protoc &> /dev/null; then
        missing_deps+=("protoc")
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
    
    if [ -f "scripts/generate-certs.sh" ]; then
        ./scripts/generate-certs.sh
        print_success "TLS certificates generated"
    else
        print_warning "generate-certs.sh not found, skipping TLS certificate generation"
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

# Function to check what needs to be set up
check_setup_requirements() {
    print_status "Checking setup requirements..."
    
    local needs_env=false
    local needs_cert=false
    local needs_proto=false
    local needs_deps=false
    
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
    
    # Check proto files
    if [ ! -f "auth-service/internal/authpb/auth_grpc.pb.go" ]; then
        print_status "Proto files not generated - needs generation"
        needs_proto=true
    fi
    
    # Check Go dependencies
    if [ ! -f "auth-service/go.sum" ] || [ ! -f "core-service/go.sum" ] || [ ! -f "notification-service/go.sum" ]; then
        print_status "Go dependencies not initialized - needs setup"
        needs_deps=true
    fi
    
    echo "$needs_env $needs_cert $needs_proto $needs_deps"
}

# Function to show help
show_help() {
    echo "⚙️  Environment setup script for subs-service"
    echo "============================================="
    echo
    echo "Usage:"
    echo "  ./scripts/setup.sh [COMMAND] [OPTIONS]"
    echo
    echo "Commands:"
    echo "  env                   Create/validate .env file"
    echo "  cert                  Generate TLS certificates"
    echo "  proto                 Generate proto files"
    echo "  deps                  Update Go dependencies"
    echo "  validate              Validate environment variables"
    echo "  check                 Check what needs to be set up"
    echo "  all                   Run all setup operations"
    echo "  help                  Show this help"
    echo
    echo "Options:"
    echo "  --force               Force regeneration even if files exist"
    echo
    echo "Examples:"
    echo "  ./scripts/setup.sh all              # Run all setup operations"
    echo "  ./scripts/setup.sh env              # Only create .env file"
    echo "  ./scripts/setup.sh cert             # Only generate certificates"
    echo "  ./scripts/setup.sh proto            # Only generate proto files"
    echo "  ./scripts/setup.sh deps             # Only update dependencies"
    echo "  ./scripts/setup.sh validate         # Only validate environment"
    echo "  ./scripts/setup.sh check            # Check what needs setup"
    echo "  ./scripts/setup.sh env --force      # Force recreate .env file"
}

# Function to run all setup operations
run_all_setup() {
    print_status "Running all setup operations..."
    
    setup_env_file
    echo
    
    generate_tls_certificates
    echo
    
    update_go_dependencies
    echo
    
    check_proto_files
    echo
    
    validate_environment
    echo
    
    print_success "All setup operations completed!"
}

# Main function
main() {
    local command=${1:-"help"}
    local force=${2:-""}
    
    # Check if we are in the root directory of the project
    if [ ! -f "docker-compose.yaml" ]; then
        print_error "docker-compose.yaml not found. Run the script from the root directory of the project."
        exit 1
    fi
    
    # Check dependencies
    check_dependencies
    echo
    
    case $command in
        env)
            if [ "$force" = "--force" ] || [ ! -f ".env" ]; then
                setup_env_file
            else
                print_warning ".env file already exists. Use --force to recreate."
            fi
            ;;
        cert)
            if [ "$force" = "--force" ] || [ ! -f "certs/server-cert.pem" ] || [ ! -f "certs/server-key.pem" ]; then
                generate_tls_certificates
            else
                print_warning "Certificates already exist. Use --force to regenerate."
            fi
            ;;
        proto)
            if [ "$force" = "--force" ] || [ ! -f "auth-service/internal/authpb/auth_grpc.pb.go" ]; then
                check_proto_files
            else
                print_warning "Proto files already generated. Use --force to regenerate."
            fi
            ;;
        deps)
            update_go_dependencies
            ;;
        validate)
            validate_environment
            ;;
        check)
            local requirements=$(check_setup_requirements)
            read -r needs_env needs_cert needs_proto needs_deps <<< "$requirements"
            
            echo "Setup requirements:"
            echo "=================="
            
            if [ "$needs_env" = true ]; then
                echo "❌ .env file needs to be created"
            else
                echo "✅ .env file exists"
            fi
            
            if [ "$needs_cert" = true ]; then
                echo "❌ TLS certificates need to be generated"
            else
                echo "✅ TLS certificates exist"
            fi
            
            if [ "$needs_proto" = true ]; then
                echo "❌ Proto files need to be generated"
            else
                echo "✅ Proto files exist"
            fi
            
            if [ "$needs_deps" = true ]; then
                echo "❌ Go dependencies need to be updated"
            else
                echo "✅ Go dependencies exist"
            fi
            
            if [ "$needs_env" = false ] && [ "$needs_cert" = false ] && [ "$needs_proto" = false ] && [ "$needs_deps" = false ]; then
                echo
                print_success "All setup requirements are satisfied!"
            else
                echo
                print_status "Run './scripts/setup.sh all' to complete setup"
            fi
            ;;
        all)
            run_all_setup
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
    print_success "Setup operation completed!"
}

# Run main function
main "$@"
