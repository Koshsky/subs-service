# API для управления подписками - Обновленная документация

## Запуск сервисов
```bash
# Запуск всех сервисов
docker-compose up -d

# Проверка статуса
docker ps
```

## Быстрое тестирование
```bash
# Полное тестирование всех endpoint'ов
./scripts/test-api.sh
```

## API Endpoints

### 1. Регистрация пользователя
```bash
curl -X POST http://localhost:8080/auth/register \
     -H "Content-Type: application/json" \
     -d '{
           "email": "user@example.com",
           "password": "password123"
         }' \
     | jq
```

### 2. Авторизация пользователя
```bash
curl -X POST http://localhost:8080/auth/login \
     -H "Content-Type: application/json" \
     -d '{
           "email": "user@example.com",
           "password": "password123"
         }' \
     -c cookies.txt \
     | jq
```

### 3. Создать подписку
```bash
curl -X POST http://localhost:8080/api/subscriptions \
     -H "Content-Type: application/json" \
     -b cookies.txt \
     -d '{
       "service_name": "Yandex Plus",
       "price": 450,
       "start_date": "07-2025"
     }' | jq
```

### 4. Получить все подписки пользователя
```bash
curl -b cookies.txt \
     http://localhost:8080/api/subscriptions | jq
```

### 5. Получить подписку по ID
```bash
curl -b cookies.txt \
     http://localhost:8080/api/subscriptions/1 | jq
```

### 6. Обновить подписку
```bash
curl -X PUT http://localhost:8080/api/subscriptions/1 \
     -H "Content-Type: application/json" \
     -b cookies.txt \
     -d '{
       "service_name": "Yandex Plus Premium",
       "price": 650,
       "start_date": "07-2025",
       "end_date": "12-2026"
     }' | jq
```

### 7. Удалить подписку
```bash
curl -X DELETE http://localhost:8080/api/subscriptions/1 \
     -b cookies.txt | jq
```

## Структура данных

### Пользователь
```json
{
  "id": 1,
  "email": "user@example.com"
}
```

### Подписка
```json
{
  "ID": 1,
  "CreatedAt": "2025-08-04T20:14:12.160093Z",
  "UpdatedAt": "2025-08-04T20:14:12.160093Z",
  "DeletedAt": null,
  "service_name": "Yandex Plus",
  "price": 450,
  "user_id": 1,
  "start_date": "07-2025",
  "end_date": null
}
```

## Отладка

### Профилирование (pprof)
```bash
# Heap профиль
curl http://localhost:8080/internal/debug/pprof/heap

# Goroutines
curl http://localhost:8080/internal/debug/pprof/goroutine

# CPU профиль
curl http://localhost:8080/internal/debug/pprof/profile
```

### Логи контейнеров
```bash
# Логи core-service
docker logs core_service

# Логи auth-service
docker logs auth_service

# Логи базы данных
docker logs subs_db
```

### Остановка сервисов
```bash
# Остановка контейнеров
docker-compose down

# Остановка с удалением volumes (полная очистка)
docker-compose down -v
```