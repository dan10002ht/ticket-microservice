#!/usr/bin/env bash
# Generate self-signed TLS certificates for gRPC mTLS.
# Usage: bash scripts/gen-certs.sh [output_dir]
#
# Creates: ca.{crt,key}, server.{crt,key}, client.{crt,key}
# Default output: ./certs/

set -euo pipefail

CERT_DIR="${1:-$(pwd)/certs}"
DAYS=365
RSA_BITS=2048

# All service DNS names for SAN
SERVICE_NAMES=(
  localhost
  auth-service
  user-service
  ticket-service
  event-service
  booking-service
  payment-service
  checkin-service
  invoice-service
  realtime-service
  booking-worker
  email-worker
  gateway
  "*.ticket-system.svc.cluster.local"
)

# Build SAN string
SAN="DNS:${SERVICE_NAMES[0]}"
for name in "${SERVICE_NAMES[@]:1}"; do
  SAN="${SAN},DNS:${name}"
done
SAN="${SAN},IP:127.0.0.1,IP:0.0.0.0"

echo "==> Generating certificates in ${CERT_DIR}"
mkdir -p "$CERT_DIR"

# 1. Root CA
echo "==> Creating Root CA..."
openssl req -x509 -newkey rsa:${RSA_BITS} -nodes \
  -keyout "${CERT_DIR}/ca.key" \
  -out "${CERT_DIR}/ca.crt" \
  -days ${DAYS} \
  -subj "/CN=ticket-system-ca/O=TicketSystem/C=VN"

# 2. Server certificate
echo "==> Creating Server certificate..."
openssl req -newkey rsa:${RSA_BITS} -nodes \
  -keyout "${CERT_DIR}/server.key" \
  -out "${CERT_DIR}/server.csr" \
  -subj "/CN=ticket-system-server/O=TicketSystem/C=VN"

openssl x509 -req -in "${CERT_DIR}/server.csr" \
  -CA "${CERT_DIR}/ca.crt" -CAkey "${CERT_DIR}/ca.key" -CAcreateserial \
  -out "${CERT_DIR}/server.crt" \
  -days ${DAYS} \
  -extfile <(printf "subjectAltName=${SAN}\nkeyUsage=digitalSignature,keyEncipherment\nextendedKeyUsage=serverAuth,clientAuth")

# 3. Client certificate (for mTLS)
echo "==> Creating Client certificate..."
openssl req -newkey rsa:${RSA_BITS} -nodes \
  -keyout "${CERT_DIR}/client.key" \
  -out "${CERT_DIR}/client.csr" \
  -subj "/CN=ticket-system-client/O=TicketSystem/C=VN"

openssl x509 -req -in "${CERT_DIR}/client.csr" \
  -CA "${CERT_DIR}/ca.crt" -CAkey "${CERT_DIR}/ca.key" -CAcreateserial \
  -out "${CERT_DIR}/client.crt" \
  -days ${DAYS} \
  -extfile <(printf "keyUsage=digitalSignature,keyEncipherment\nextendedKeyUsage=clientAuth")

# Cleanup CSR files
rm -f "${CERT_DIR}"/*.csr "${CERT_DIR}"/*.srl

# Set permissions
chmod 644 "${CERT_DIR}"/*.crt
chmod 600 "${CERT_DIR}"/*.key

echo ""
echo "==> Certificates generated successfully:"
ls -la "${CERT_DIR}"
echo ""
echo "Files:"
echo "  CA:     ${CERT_DIR}/ca.crt, ${CERT_DIR}/ca.key"
echo "  Server: ${CERT_DIR}/server.crt, ${CERT_DIR}/server.key"
echo "  Client: ${CERT_DIR}/client.crt, ${CERT_DIR}/client.key"
echo ""
echo "Test with grpcurl:"
echo "  grpcurl -cacert ${CERT_DIR}/ca.crt localhost:50052 grpc.health.v1.Health/Check"
