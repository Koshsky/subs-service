#!/bin/bash
set -euo pipefail

# Always run from repo root (script is in scripts/)
cd "$(dirname "$0")/.."

CERTS_DIR="certs"
KEY_FILE="$CERTS_DIR/server-key.pem"
CERT_FILE="$CERTS_DIR/server-cert.pem"

mkdir -p "$CERTS_DIR"

# Clean old certs
rm -f "$KEY_FILE" "$CERT_FILE"

# Create a temporary OpenSSL config with SANs for auth-service and localhost
TMP_CNF="$(mktemp)"
cat > "$TMP_CNF" << 'EOF'
[ req ]
distinguished_name = dn
x509_extensions = v3_req
prompt = no

[ dn ]
C  = US
ST = CA
L  = San Francisco
O  = MyOrg
OU = MyOrgUnit
CN = auth-service

[ v3_req ]
subjectAltName = @alt_names
basicConstraints = CA:false
keyUsage = digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth

[ alt_names ]
DNS.1 = auth-service
DNS.2 = localhost
EOF

# Generate key and self-signed certificate with SANs
openssl req \
  -new -newkey rsa:4096 -nodes \
  -keyout "$KEY_FILE" \
  -x509 -days 365 \
  -out "$CERT_FILE" \
  -config "$TMP_CNF"

# Cleanup temp config
rm -f "$TMP_CNF"

# Permissions
chmod 600 "$KEY_FILE"
chmod 644 "$CERT_FILE"

echo "TLS certificates (with SANs) generated successfully!"
printf "Server certificate: %s\n" "$CERT_FILE"
printf "Server private key: %s\n\n" "$KEY_FILE"
echo "Note: Self-signed certs for development only. For production, use a trusted CA."