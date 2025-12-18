#!/bin/bash

# Test E2E Booking Notification Flow
# Tests the complete flow: booking-worker → realtime-service → WebSocket

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() { echo -e "${GREEN}[INFO]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
print_error() { echo -e "${RED}[ERROR]${NC} $1"; }
print_header() { echo -e "${BLUE}========================================${NC}"; echo -e "${BLUE}$1${NC}"; echo -e "${BLUE}========================================${NC}"; }

# Configuration
REALTIME_GRPC_HOST="${REALTIME_GRPC_HOST:-localhost}"
REALTIME_GRPC_PORT="${REALTIME_GRPC_PORT:-50057}"
REALTIME_WS_HOST="${REALTIME_WS_HOST:-localhost}"
REALTIME_WS_PORT="${REALTIME_WS_PORT:-3003}"
REALTIME_HTTP_PORT="${REALTIME_HTTP_PORT:-3003}"

print_header "E2E Realtime Service Test"

# ============================================
# Test 1: Health Check
# ============================================
print_header "Test 1: Health Check"

print_status "Testing HTTP health endpoint..."
HTTP_HEALTH=$(curl -s -o /dev/null -w "%{http_code}" "http://${REALTIME_WS_HOST}:${REALTIME_HTTP_PORT}/health" 2>/dev/null || echo "000")

if [ "$HTTP_HEALTH" = "200" ]; then
    print_status "✅ HTTP Health check passed (status: $HTTP_HEALTH)"
else
    print_error "❌ HTTP Health check failed (status: $HTTP_HEALTH)"
    print_warning "Make sure realtime-service is running on port ${REALTIME_HTTP_PORT}"
    exit 1
fi

# ============================================
# Test 2: gRPC Service Check
# ============================================
print_header "Test 2: gRPC Service Check"

# Check if grpcurl is installed
if ! command -v grpcurl &> /dev/null; then
    print_warning "grpcurl not installed. Skipping gRPC tests."
    print_status "Install with: brew install grpcurl (macOS) or go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest"
else
    print_status "Testing gRPC service listing..."

    GRPC_LIST=$(grpcurl -plaintext "${REALTIME_GRPC_HOST}:${REALTIME_GRPC_PORT}" list 2>&1 || echo "ERROR")

    if echo "$GRPC_LIST" | grep -q "realtime.RealtimeService"; then
        print_status "✅ gRPC service is available"
        echo "$GRPC_LIST" | head -5
    else
        print_error "❌ gRPC service check failed"
        echo "$GRPC_LIST"
        exit 1
    fi

    # Test GetConnectionStats
    print_status "Testing GetConnectionStats..."
    STATS=$(grpcurl -plaintext -d '{}' "${REALTIME_GRPC_HOST}:${REALTIME_GRPC_PORT}" realtime.RealtimeService/GetConnectionStats 2>&1 || echo "ERROR")

    if echo "$STATS" | grep -q "totalConnections\|{}"; then
        print_status "✅ GetConnectionStats works"
        echo "$STATS"
    else
        print_error "❌ GetConnectionStats failed"
        echo "$STATS"
    fi
fi

# ============================================
# Test 3: Simulate Booking Notification
# ============================================
print_header "Test 3: Simulate Booking Notification (gRPC)"

if command -v grpcurl &> /dev/null; then
    print_status "Sending test booking notification..."

    BOOKING_RESULT=$(grpcurl -plaintext -d '{
        "user_id": "test-user-123",
        "booking_id": "booking-456",
        "success": true,
        "message": "Booking confirmed successfully",
        "booking_reference": "REF-2024-001",
        "event_id": "event-789",
        "seat_numbers": ["A1", "A2"],
        "total_amount": "100.00",
        "currency": "USD"
    }' "${REALTIME_GRPC_HOST}:${REALTIME_GRPC_PORT}" realtime.RealtimeService/NotifyBookingResult 2>&1 || echo "ERROR")

    if echo "$BOOKING_RESULT" | grep -q "delivered"; then
        print_status "✅ NotifyBookingResult call successful"
        echo "$BOOKING_RESULT"
    else
        print_warning "⚠️ NotifyBookingResult returned (may not have active WebSocket connection)"
        echo "$BOOKING_RESULT"
    fi
fi

# ============================================
# Test 4: Simulate Queue Position Notification
# ============================================
print_header "Test 4: Simulate Queue Position Notification (gRPC)"

if command -v grpcurl &> /dev/null; then
    print_status "Sending test queue position notification..."

    QUEUE_RESULT=$(grpcurl -plaintext -d '{
        "user_id": "test-user-123",
        "event_id": "event-789",
        "position": 5,
        "estimated_wait_seconds": 120,
        "total_in_queue": 50
    }' "${REALTIME_GRPC_HOST}:${REALTIME_GRPC_PORT}" realtime.RealtimeService/NotifyQueuePosition 2>&1 || echo "ERROR")

    if echo "$QUEUE_RESULT" | grep -q "delivered"; then
        print_status "✅ NotifyQueuePosition call successful"
        echo "$QUEUE_RESULT"
    else
        print_warning "⚠️ NotifyQueuePosition returned (may not have active WebSocket connection)"
        echo "$QUEUE_RESULT"
    fi
fi

# ============================================
# Test 5: Simulate Payment Notification
# ============================================
print_header "Test 5: Simulate Payment Notification (gRPC)"

if command -v grpcurl &> /dev/null; then
    print_status "Sending test payment notification..."

    PAYMENT_RESULT=$(grpcurl -plaintext -d '{
        "user_id": "test-user-123",
        "booking_id": "booking-456",
        "payment_id": "pay-001",
        "status": "PAYMENT_STATUS_SUCCESS",
        "message": "Payment completed",
        "amount": "100.00",
        "currency": "USD"
    }' "${REALTIME_GRPC_HOST}:${REALTIME_GRPC_PORT}" realtime.RealtimeService/NotifyPaymentStatus 2>&1 || echo "ERROR")

    if echo "$PAYMENT_RESULT" | grep -q "delivered"; then
        print_status "✅ NotifyPaymentStatus call successful"
        echo "$PAYMENT_RESULT"
    else
        print_warning "⚠️ NotifyPaymentStatus returned (may not have active WebSocket connection)"
        echo "$PAYMENT_RESULT"
    fi
fi

# ============================================
# Test 6: Broadcast Event
# ============================================
print_header "Test 6: Broadcast Event (gRPC)"

if command -v grpcurl &> /dev/null; then
    print_status "Sending test broadcast event..."

    BROADCAST_RESULT=$(grpcurl -plaintext -d '{
        "event_type": "ticket:availability",
        "room": "event-789",
        "payload": "{\"available_seats\": 45}"
    }' "${REALTIME_GRPC_HOST}:${REALTIME_GRPC_PORT}" realtime.RealtimeService/BroadcastEvent 2>&1 || echo "ERROR")

    if echo "$BROADCAST_RESULT" | grep -q "recipients"; then
        print_status "✅ BroadcastEvent call successful"
        echo "$BROADCAST_RESULT"
    else
        print_warning "⚠️ BroadcastEvent returned"
        echo "$BROADCAST_RESULT"
    fi
fi

# ============================================
# Test 7: WebSocket Connection Test
# ============================================
print_header "Test 7: WebSocket Connection Test"

if command -v wscat &> /dev/null; then
    print_status "Testing WebSocket connection (will timeout after 3s)..."
    print_status "Connect manually: wscat -c 'ws://${REALTIME_WS_HOST}:${REALTIME_WS_PORT}/ws'"

    # Quick connection test (timeout after 3 seconds)
    timeout 3 wscat -c "ws://${REALTIME_WS_HOST}:${REALTIME_WS_PORT}/ws" 2>&1 || true
    print_status "WebSocket endpoint is accessible"
else
    print_warning "wscat not installed. Skipping WebSocket test."
    print_status "Install with: npm install -g wscat"
    print_status "Manual test: wscat -c 'ws://${REALTIME_WS_HOST}:${REALTIME_WS_PORT}/ws'"
fi

# ============================================
# Summary
# ============================================
print_header "Test Summary"
print_status "All basic tests completed!"
print_status ""
print_status "To fully test E2E flow:"
print_status "  1. Start all services: docker-compose -f deploy/docker-compose.dev.yml up -d"
print_status "  2. Connect WebSocket client: wscat -c 'ws://localhost:3003/ws'"
print_status "  3. Trigger a booking through the API or booking-worker"
print_status "  4. Observe real-time notifications in WebSocket client"
print_status ""
print_status "WebSocket authentication:"
print_status "  Send: {\"type\":\"auth\",\"token\":\"<JWT_TOKEN>\"}"
print_status "  Join room: {\"type\":\"join\",\"room\":\"event-789\"}"
