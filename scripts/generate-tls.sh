    # Check for certs directory
    if [ ! -d "certs" ]; then
        echo "Creating certs directory..."
        mkdir -p certs
    fi
    
    # Check for generate-certs.sh script
    if [ -f "scripts/generate-certs.sh" ]; then
        echo "Running generate-certs.sh..."
        ./scripts/generate-certs.sh
        echo "TLS certificates generated"
    else
        echo "generate-certs.sh not found, creating basic certificates..."

        # Create self-signed certificates
        local cert_file="certs/server-cert.pem"
        local key_file="certs/server-key.pem"

        # Generate private key
        openssl genrsa -out "$key_file" 2048 2>/dev/null || {
            echo "Failed to generate private key"
            return 1
        }

        # Create configuration file for certificate
        local config_file="certs/cert.conf"
        cat > "$config_file" << 'EOF'
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req
prompt = no

[req_distinguished_name]
C = RU
ST = Moscow
L = Moscow
O = SubsService
OU = Development
CN = localhost

[v3_req]
keyUsage = keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = auth-service
DNS.3 = core-service
DNS.4 = notification-service
DNS.5 = *.localhost
IP.1 = 127.0.0.1
IP.2 = ::1
EOF

        # Generate certificate
        openssl req -new -x509 -key "$key_file" -out "$cert_file" -days 365 -config "$config_file" 2>/dev/null || {
            echo "Failed to generate certificate"
            return 1
        }

        # Set correct permissions
        chmod 600 "$key_file"
        chmod 644 "$cert_file"

        # Remove temporary config
        rm -f "$config_file"

        echo "Basic TLS certificates created"
    fi
    
    # Check if certificates exist
    if [ ! -f "certs/server-cert.pem" ] || [ ! -f "certs/server-key.pem" ]; then
        echo "TLS certificates not found after generation"
        return 1
    fi