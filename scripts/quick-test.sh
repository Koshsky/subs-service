#!/bin/bash

# Быстрый тест основных функций API
API_URL="http://localhost:8080"

echo "🚀 Быстрый тест API"
echo "=================="

# Проверяем, что сервисы запущены
echo "1. Проверка доступности сервисов..."
if curl -s "$API_URL/auth/register" > /dev/null; then
    echo "✅ Core service доступен"
else
    echo "❌ Core service недоступен"
    exit 1
fi

# Регистрация
echo
echo "2. Регистрация тестового пользователя..."
REGISTER_RESP=$(curl -s -X POST "$API_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d '{"email":"quicktest@example.com","password":"password123"}')

if echo "$REGISTER_RESP" | grep -q "successfully\|exists"; then
    echo "✅ Регистрация работает"
else
    echo "❌ Проблема с регистрацией: $REGISTER_RESP"
fi

# Авторизация
echo
echo "3. Авторизация..."
LOGIN_RESP=$(curl -s -X POST "$API_URL/auth/login" \
    -H "Content-Type: application/json" \
    -c /tmp/quick_cookies.txt \
    -d '{"email":"quicktest@example.com","password":"password123"}')

if echo "$LOGIN_RESP" | grep -q "Successful"; then
    echo "✅ Авторизация работает"
else
    echo "❌ Проблема с авторизацией: $LOGIN_RESP"
fi

# Создание подписки
echo
echo "4. Создание подписки..."
SUB_RESP=$(curl -s -X POST "$API_URL/api/subscriptions" \
    -H "Content-Type: application/json" \
    -b /tmp/quick_cookies.txt \
    -d '{"service_name":"Test Service","price":100,"start_date":"01-2025"}')

if echo "$SUB_RESP" | grep -q '"ID":'; then
    echo "✅ Создание подписки работает"
else
    echo "❌ Проблема с созданием подписки: $SUB_RESP"
fi

# Получение подписок
echo
echo "5. Получение списка подписок..."
LIST_RESP=$(curl -s -b /tmp/quick_cookies.txt "$API_URL/api/subscriptions")

if echo "$LIST_RESP" | grep -q '"service_name"'; then
    echo "✅ Получение подписок работает"
else
    echo "❌ Проблема с получением подписок: $LIST_RESP"
fi

# Очистка
rm -f /tmp/quick_cookies.txt

echo
echo "🎉 Быстрый тест завершен!"
echo "Для полного тестирования запустите: ./scripts/test-api.sh"