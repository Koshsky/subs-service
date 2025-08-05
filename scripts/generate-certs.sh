#!/bin/bash

# Create certs directory if it doesn't exist
mkdir -p certs

# Generate private key for server
openssl genrsa -out certs/server-key.pem 4096

# Generate certificate signing request for server
openssl req -new -x509 -key certs/server-key.pem -out certs/server-cert.pem -days 365 -subj "/C=US/ST=CA/L=San Francisco/O=MyOrg/OU=MyOrgUnit/CN=localhost"

# Set appropriate permissions
chmod 600 certs/server-key.pem
chmod 644 certs/server-cert.pem

echo "TLS certificates generated successfully!"
echo "Server certificate: certs/server-cert.pem"
echo "Server private key: certs/server-key.pem"
echo ""
echo "Note: These are self-signed certificates for development use only."
echo "For production, use certificates from a trusted Certificate Authority."