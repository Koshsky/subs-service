# API для управления подписками

## 1. Создать подписку

```bash
curl -X POST http://localhost:8080/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 450,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }' | jq
```

## 2. Получить все подписки

```bash
curl http://localhost:8080/subscriptions | jq
```

## 3. Получить подписку по id

```bash
curl http://localhost:8080/subscriptions/1 | jq
```

## 4. Обновить подписку

```bash
curl -X PUT http://localhost:8080/subscriptions/1 \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 450,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025",
    "end_date": "12-2026"
  }' | jq
```

## 5. Удалить подписку

```bash
curl -X DELETE http://localhost:8080/subscriptions/1 | jq
```

## 6. Получить сумму по фильтрам

```bash
curl "http://localhost:8080/subscriptions/total?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&service=Yandex%20Plus&start_month=07-2025&end_month=07-2025" | jq
```