#!/bin/bash

# Script to test notification service integration
set -e

echo "🧪 Testing Notification Service Integration"
echo "=========================================="

# Check if services are running
echo "🔍 Checking if services are running..."

# Check RabbitMQ
if ! curl -s http://localhost:15672/api/overview > /dev/null 2>&1; then
    echo "❌ RabbitMQ is not running or not accessible"
    echo "   Start services with: docker-compose up -d"
    exit 1
fi
echo "✅ RabbitMQ is running"

# Check notification service
if ! curl -s http://localhost:8082/health > /dev/null 2>&1; then
    echo "❌ Notification service is not running or not accessible"
    echo "   Start services with: docker-compose up -d"
    exit 1
fi
echo "✅ Notification service is running"

# Check auth service
if ! netstat -tuln | grep :50051 > /dev/null 2>&1; then
    echo "❌ Auth service is not running or not accessible"
    echo "   Start services with: docker-compose up -d"
    exit 1
fi
echo "✅ Auth service is running"

echo
echo "📊 RabbitMQ Management Interface:"
echo "   URL: http://localhost:15672"
echo "   Username: guest"
echo "   Password: guest"

echo
echo "🔍 Notification Service Health Check:"
curl -s http://localhost:8082/health | jq .

echo
echo "📝 To test the integration:"
echo "   1. Register a new user via auth service"
echo "   2. Check RabbitMQ management interface for messages"
echo "   3. Check notification service logs for event processing"
echo "   4. Check notification database for created records"

echo
echo "📋 Useful commands:"
echo "   # View notification service logs:"
echo "   docker-compose logs -f notification-service"
echo ""
echo "   # View RabbitMQ logs:"
echo "   docker-compose logs -f rabbitmq"
echo ""
echo "   # Check notification database:"
echo "   docker-compose exec notify-db psql -U notify_user -d notify_db -c 'SELECT * FROM notifications;'"

echo
echo "✅ Integration test setup complete!"
