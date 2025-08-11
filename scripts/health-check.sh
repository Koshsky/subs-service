#!/bin/bash

# Script to check health of all services
set -e

echo "🏥 Health Check for All Services"
echo "================================"
echo

# Function to check service health
check_service_health() {
    local service_name=$1
    local url=$2
    local endpoint=$3
    
    echo -n "   $service_name: "
    
    if curl -s --max-time 5 "$url$endpoint" > /dev/null 2>&1; then
        echo "✅ Healthy"
        return 0
    else
        echo "❌ Unhealthy"
        return 1
    fi
}

# Function to check container status
check_container_status() {
    local container_name=$1
    local service_name=$2
    
    echo -n "   $service_name: "
    
    if docker ps --format "table {{.Names}}\t{{.Status}}" | grep -q "$container_name.*Up"; then
        echo "✅ Running"
        return 0
    else
        echo "❌ Not running"
        return 1
    fi
}

# Function to get container uptime
get_container_uptime() {
    local container_name=$1
    docker ps --format "table {{.Names}}\t{{.RunningFor}}" | grep "$container_name" | awk '{print $2}' 2>/dev/null || echo "Unknown"
}

# Check container status
echo "📦 Container Status:"
check_container_status "auth_service" "Auth Service"
check_container_status "core_service" "Core Service"
check_container_status "notification_service" "Notification Service"
check_container_status "auth_db" "Auth Database"
check_container_status "core_db" "Core Database"
check_container_status "notify_db" "Notification Database"
check_container_status "rabbitmq" "RabbitMQ"

echo
echo "⏱️  Container Uptime:"
echo "   Auth Service: $(get_container_uptime auth_service)"
echo "   Core Service: $(get_container_uptime core_service)"
echo "   Notification Service: $(get_container_uptime notification_service)"
echo "   Auth Database: $(get_container_uptime auth_db)"
echo "   Core Database: $(get_container_uptime core_db)"
echo "   Notification Database: $(get_container_uptime notify_db)"
echo "   RabbitMQ: $(get_container_uptime rabbitmq)"

echo
echo "🌐 Service Health Endpoints:"
check_service_health "Core Service" "http://localhost:8080" "/health"
check_service_health "Notification Service" "http://localhost:8082" "/health"
check_service_health "Auth Service Health" "http://localhost:8081" "/health"

echo
echo "📊 Database Connections:"
# Check if databases are accessible (simplified check)
echo -n "   Auth Database: "
if docker exec auth_db pg_isready -U auth_user -d auth_db > /dev/null 2>&1; then
    echo "✅ Connected"
else
    echo "❌ Not accessible"
fi

echo -n "   Core Database: "
if docker exec core_db pg_isready -U core_user -d core_db > /dev/null 2>&1; then
    echo "✅ Connected"
else
    echo "❌ Not accessible"
fi

echo -n "   Notification Database: "
if docker exec notify_db pg_isready -U notify_user -d notify_db > /dev/null 2>&1; then
    echo "✅ Connected"
else
    echo "❌ Not accessible"
fi

echo
echo "🐰 RabbitMQ Status:"
echo -n "   RabbitMQ Management: "
if curl -s --max-time 5 "http://localhost:15672" > /dev/null 2>&1; then
    echo "✅ Available"
else
    echo "❌ Not available"
fi

echo -n "   RabbitMQ AMQP: "
if timeout 5 bash -c "</dev/tcp/localhost/5672" 2>/dev/null; then
    echo "✅ Available"
else
    echo "❌ Not available"
fi

echo
echo "🔗 Port Status:"
echo -n "   Core Service (8080): "
if netstat -tuln 2>/dev/null | grep -q ":8080 "; then
    echo "✅ Listening"
else
    echo "❌ Not listening"
fi

echo -n "   Auth Service Health (8081): "
if netstat -tuln 2>/dev/null | grep -q ":8081 "; then
    echo "✅ Listening"
else
    echo "❌ Not listening"
fi

echo -n "   Notification Service (8082): "
if netstat -tuln 2>/dev/null | grep -q ":8082 "; then
    echo "✅ Listening"
else
    echo "❌ Not listening"
fi

echo -n "   RabbitMQ Management (15672): "
if netstat -tuln 2>/dev/null | grep -q ":15672 "; then
    echo "✅ Listening"
else
    echo "❌ Not listening"
fi

echo
echo "📈 Quick Performance Check:"
echo -n "   Core Service Response Time: "
CORE_RESPONSE_TIME=$(curl -s -w "%{time_total}" -o /dev/null "http://localhost:8080/health" 2>/dev/null || echo "0")
if [ "$CORE_RESPONSE_TIME" != "0" ]; then
    echo "✅ ${CORE_RESPONSE_TIME}s"
else
    echo "❌ No response"
fi

echo -n "   Notification Service Response Time: "
NOTIFY_RESPONSE_TIME=$(curl -s -w "%{time_total}" -o /dev/null "http://localhost:8082/health" 2>/dev/null || echo "0")
if [ "$NOTIFY_RESPONSE_TIME" != "0" ]; then
    echo "✅ ${NOTIFY_RESPONSE_TIME}s"
else
    echo "❌ No response"
fi

echo
echo "🔍 Quick Log Check (last 5 lines):"
echo "   Auth Service:"
docker logs --tail 5 auth_service 2>/dev/null | sed 's/^/     /' || echo "     No logs available"
echo
echo "   Core Service:"
docker logs --tail 5 core_service 2>/dev/null | sed 's/^/     /' || echo "     No logs available"
echo
echo "   Notification Service:"
docker logs --tail 5 notification_service 2>/dev/null | sed 's/^/     /' || echo "     No logs available"

echo
echo "✅ Health check completed!"
echo
echo "💡 Useful commands:"
echo "   docker-compose ps                    # Show all containers"
echo "   docker-compose logs -f [service]     # Follow logs for specific service"
echo "   docker stats                         # Show resource usage"
echo "   ./scripts/test-api.sh               # Run full API tests"
