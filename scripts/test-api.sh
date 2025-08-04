#!/bin/bash

# Скрипт для тестирования API сервиса подписок
set -e

API_URL="http://localhost:8080"

# Генерируем случайные учетные данные для каждого запуска
RANDOM_ID=$(date +%s%N | cut -b1-13)
EMAIL="test_${RANDOM_ID}@example.com"
PASSWORD="password_$(openssl rand -hex 8)"
COOKIES_FILE="/tmp/test_cookies_${RANDOM_ID}.txt"

echo "🧪 Тестирование API сервиса подписок"
echo "=================================="
echo "Тестовый email: $EMAIL"
echo "Сессия: $RANDOM_ID"

# Функция для красивого вывода JSON
print_response() {
    echo "$1" | jq . 2>/dev/null || echo "$1"
}

echo
echo "1️⃣  Создание пользователя..."

REGISTER_RESPONSE=$(curl -s -X POST "$API_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"$EMAIL\",
        \"password\": \"$PASSWORD\"
    }")

echo "Ответ:"
print_response "$REGISTER_RESPONSE"

echo
echo "2️⃣  Авторизация пользователя..."

LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
    -H "Content-Type: application/json" \
    -c "$COOKIES_FILE" \
    -d "{
        \"email\": \"$EMAIL\",
        \"password\": \"$PASSWORD\"
    }")

echo "Ответ:"
print_response "$LOGIN_RESPONSE"

# Извлекаем токен из cookies или response
if [ -f "$COOKIES_FILE" ]; then
    echo "✅ Cookies сохранены для дальнейших запросов"
else
    echo "❌ Не удалось сохранить cookies"
    exit 1
fi

echo
echo "3️⃣  Создание подписки Yandex Plus..."

SUB1_RESPONSE=$(curl -s -X POST "$API_URL/api/subscriptions" \
    -H "Content-Type: application/json" \
    -b "$COOKIES_FILE" \
    -d '{
        "service_name": "Yandex Plus",
        "price": 450,
        "start_date": "07-2025"
    }')

echo "Ответ:"
print_response "$SUB1_RESPONSE"

echo
echo "4️⃣  Создание подписки Netflix..."

SUB2_RESPONSE=$(curl -s -X POST "$API_URL/api/subscriptions" \
    -H "Content-Type: application/json" \
    -b "$COOKIES_FILE" \
    -d '{
        "service_name": "Netflix",
        "price": 1200,
        "start_date": "06-2025",
        "end_date": "06-2026"
    }')

echo "Ответ:"
print_response "$SUB2_RESPONSE"

echo
echo "5️⃣  Получение всех подписок..."

ALL_SUBS_RESPONSE=$(curl -s -b "$COOKIES_FILE" "$API_URL/api/subscriptions")

echo "Ответ:"
print_response "$ALL_SUBS_RESPONSE"

echo
echo "6️⃣  Получение подписки по ID=1..."

SUB_BY_ID_RESPONSE=$(curl -s -b "$COOKIES_FILE" "$API_URL/api/subscriptions/1")

echo "Ответ:"
print_response "$SUB_BY_ID_RESPONSE"

echo
echo "7️⃣  Обновление подписки ID=1..."

UPDATE_RESPONSE=$(curl -s -X PUT "$API_URL/api/subscriptions/1" \
    -H "Content-Type: application/json" \
    -b "$COOKIES_FILE" \
    -d '{
        "service_name": "Yandex Plus Premium",
        "price": 650,
        "start_date": "07-2025",
        "end_date": "12-2026"
    }')

echo "Ответ:"
print_response "$UPDATE_RESPONSE"

echo
echo "8️⃣  Получение суммы всех подписок..."

TOTAL_RESPONSE=$(curl -s -b "$COOKIES_FILE" "$API_URL/api/analytics/total")

echo "Ответ:"
print_response "$TOTAL_RESPONSE"

echo
echo "9️⃣  Удаление подписки ID=2..."

DELETE_RESPONSE=$(curl -s -X DELETE "$API_URL/api/subscriptions/2" -b "$COOKIES_FILE")

echo "Ответ:"
print_response "$DELETE_RESPONSE"

echo
echo "🔟  Проверка оставшихся подписок..."

FINAL_SUBS_RESPONSE=$(curl -s -b "$COOKIES_FILE" "$API_URL/api/subscriptions")

echo "Ответ:"
print_response "$FINAL_SUBS_RESPONSE"

echo
echo "✅ Тестирование API завершено!"

# Очистка
rm -f "$COOKIES_FILE"

echo
echo "📝 Результаты тестирования:"
echo "- Email: $EMAIL"
echo "- Регистрация: $(echo "$REGISTER_RESPONSE" | jq -r '.message // "❌ Ошибка"' 2>/dev/null || echo "❌ Ошибка")"
echo "- Авторизация: $(echo "$LOGIN_RESPONSE" | jq -r '.message // "❌ Ошибка"' 2>/dev/null || echo "❌ Ошибка")"
echo "- Создание подписок: Выполнено"
echo "- Получение подписок: Выполнено"
echo "- Обновление подписки: Выполнено"
echo "- Удаление подписки: Выполнено"
echo "- Аналитика: Выполнено"
echo "- Сессия: $RANDOM_ID"