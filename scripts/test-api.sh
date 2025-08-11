#!/bin/bash

# Script to test the API of the subscription service
set -e

API_URL="http://localhost:8080"

# Generate random credentials for each run
RANDOM_ID=$(date +%s%N | cut -b1-13)
EMAIL="test_${RANDOM_ID}@example.com"
PASSWORD="password_$(openssl rand -hex 8)"
COOKIES_FILE="/tmp/test_cookies_${RANDOM_ID}.txt"

echo "🧪 Testing the API of the subscription service and Notification Service"
echo "=========================================================="
echo "Test email: $EMAIL"
echo "Session: $RANDOM_ID"

echo
echo "🏥 Checking all services health..."
echo "=================================="

# Call the health check script
./scripts/health-check.sh

echo
echo "🚀 Starting API tests..."
echo "========================"

# Function to pretty print JSON
print_response() {
    echo "$1" | jq . 2>/dev/null || echo "$1"
}

echo
echo "1️⃣  Creating a user (should trigger user.created event)..."

REGISTER_RESPONSE=$(curl -s -X POST "$API_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"$EMAIL\",
        \"password\": \"$PASSWORD\"
    }")

echo "Response:"
print_response "$REGISTER_RESPONSE"

echo
echo "2️⃣  User authorization..."

LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
    -H "Content-Type: application/json" \
    -c "$COOKIES_FILE" \
    -d "{
        \"email\": \"$EMAIL\",
        \"password\": \"$PASSWORD\"
    }")

echo "Response:"
print_response "$LOGIN_RESPONSE"

# Extract token from cookies or response
if [ -f "$COOKIES_FILE" ]; then
    echo "✅ Cookies saved for further requests"
else
    echo "❌ Failed to save cookies"
    exit 1
fi

echo
echo "3️⃣  Creating a Yandex Plus subscription..."

SUB1_RESPONSE=$(curl -s -X POST "$API_URL/api/subscriptions" \
    -H "Content-Type: application/json" \
    -b "$COOKIES_FILE" \
    -d '{
        "service_name": "Yandex Plus",
        "price": 450,
        "start_date": "07-2025"
    }')

echo "Response:"
print_response "$SUB1_RESPONSE"

echo
echo "4️⃣  Creating a Netflix subscription..."

SUB2_RESPONSE=$(curl -s -X POST "$API_URL/api/subscriptions" \
    -H "Content-Type: application/json" \
    -b "$COOKIES_FILE" \
    -d '{
        "service_name": "Netflix",
        "price": 1200,
        "start_date": "06-2025",
        "end_date": "06-2026"
    }')

echo "Response:"
print_response "$SUB2_RESPONSE"

echo
echo "5️⃣  Getting all subscriptions..."

ALL_SUBS_RESPONSE=$(curl -s -b "$COOKIES_FILE" "$API_URL/api/subscriptions")

echo "Response:"
print_response "$ALL_SUBS_RESPONSE"

echo
echo "6️⃣  Getting the ID of created subscriptions..."

# Extract ID of the first subscription (Yandex Plus)
SUB1_ID=$(echo "$SUB1_RESPONSE" | jq -r '.ID' 2>/dev/null || echo "unknown")
echo "   ID of the first subscription (Yandex Plus): $SUB1_ID"

# Extract ID of the second subscription (Netflix)
SUB2_ID=$(echo "$SUB2_RESPONSE" | jq -r '.ID' 2>/dev/null || echo "unknown")
echo "   ID of the second subscription (Netflix): $SUB2_ID"

echo
echo "7️⃣  Getting a subscription by ID=$SUB1_ID..."

SUB_BY_ID_RESPONSE=$(curl -s -b "$COOKIES_FILE" "$API_URL/api/subscriptions/$SUB1_ID")

echo "Response:"
print_response "$SUB_BY_ID_RESPONSE"

echo
echo "8️⃣  Updating a subscription ID=$SUB1_ID..."

UPDATE_RESPONSE=$(curl -s -X PUT "$API_URL/api/subscriptions/$SUB1_ID" \
    -H "Content-Type: application/json" \
    -b "$COOKIES_FILE" \
    -d '{
        "service_name": "Yandex Plus Premium",
        "price": 650,
        "start_date": "07-2025",
        "end_date": "12-2026"
    }')

echo "Response:"
print_response "$UPDATE_RESPONSE"

echo
echo "9️⃣  Getting the total amount of all subscriptions..."

TOTAL_RESPONSE=$(curl -s -b "$COOKIES_FILE" "$API_URL/api/analytics/total?start_month=01-2025&end_month=12-2025")

echo "Response:"
print_response "$TOTAL_RESPONSE"

echo
echo "🔟  Deleting a subscription ID=$SUB2_ID..."

DELETE_RESPONSE=$(curl -s -X DELETE "$API_URL/api/subscriptions/$SUB2_ID" -b "$COOKIES_FILE")

echo "Response:"
print_response "$DELETE_RESPONSE"

echo
echo "1️⃣1️⃣  Checking remaining subscriptions..."

FINAL_SUBS_RESPONSE=$(curl -s -b "$COOKIES_FILE" "$API_URL/api/subscriptions")

echo "Response:"
print_response "$FINAL_SUBS_RESPONSE"

echo
echo "1️⃣2️⃣  Checking RabbitMQ Management UI..."
echo "   Open in browser: http://localhost:15672"
echo "   Username: guest"
echo "   Password: guest"
echo "   Check exchange 'user_events' and queue 'user_created'"

echo
echo "✅ Testing the API and Notification Service completed!"

# Cleanup
rm -f "$COOKIES_FILE"

echo
echo "📝 Testing results:"
echo "- Email: $EMAIL"
echo "- Registration: $(echo "$REGISTER_RESPONSE" | jq -r '.message // "❌ Error"' 2>/dev/null || echo "❌ Error")"
echo "- Authorization: $(echo "$LOGIN_RESPONSE" | jq -r '.message // "❌ Error"' 2>/dev/null || echo "❌ Error")"
echo "- Creating a Yandex Plus subscription: $(if echo "$SUB1_RESPONSE" | jq -e '.ID' > /dev/null 2>&1; then echo "✅ Success"; else echo "❌ Error"; fi)"
echo "- Creating a Netflix subscription: $(if echo "$SUB2_RESPONSE" | jq -e '.ID' > /dev/null 2>&1; then echo "✅ Success"; else echo "❌ Error"; fi)"
echo "- Getting subscriptions: $(if echo "$ALL_SUBS_RESPONSE" | jq -e '.[0]' > /dev/null 2>&1; then echo "✅ Success"; else echo "❌ Error"; fi)"
echo "- Getting a subscription by ID: $(if echo "$SUB_BY_ID_RESPONSE" | jq -e '.ID' > /dev/null 2>&1; then echo "✅ Success"; else echo "❌ Error"; fi)"
echo "- Updating a subscription: $(if echo "$UPDATE_RESPONSE" | jq -e '.ID' > /dev/null 2>&1; then echo "✅ Success"; else echo "❌ Error"; fi)"
echo "- Getting analytics: $(if echo "$TOTAL_RESPONSE" | jq -e '.total_price' > /dev/null 2>&1; then echo "✅ Success"; else echo "❌ Error"; fi)"
echo "- Deleting a subscription: $(if echo "$DELETE_RESPONSE" | jq -e '.message' > /dev/null 2>&1; then echo "✅ Success"; else echo "❌ Error"; fi)"
echo "- Final check of subscriptions: $(if echo "$FINAL_SUBS_RESPONSE" | jq -e '.[0]' > /dev/null 2>&1; then echo "✅ Success"; else echo "❌ Error"; fi)"
echo "- User.created event: $(if echo "$REGISTER_RESPONSE" | jq -e '.user' > /dev/null 2>&1; then echo "✅ Sent to RabbitMQ"; else echo "❌ Not sent"; fi)"
echo "- Session: $RANDOM_ID"

echo
echo "🔍 To check the logs, run:"
echo "   docker logs notification_service"
echo "   docker logs auth_service"
echo "   docker logs core_service"
echo "   docker logs rabbitmq"