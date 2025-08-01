# API для управления подписками

## 1. Создать пользователя
```bash
curl -X POST http://localhost:8080/register \
     -H "Content-Type: application/json" \
     -d '{
           "email": "meology77@gmail.com",
           "password": "fdsjklasd"
         }' \
     | jq
```

## 2. Получить токен аутентификации
```bash
curl -X POST http://localhost:8080/login \
     -H "Content-Type: application/json" \
     -d '{
           "email": "meology77@gmail.com",
           "password": "fdsjklasd"
         }' \
     -c cookies.txt \
     | jq
```

## 3. Создать подписку
```bash
curl -X POST http://localhost:8080/subscriptions \
     -H "Content-Type: application/json" \
     -b cookies.txt \
     -d '{
       "service_name": "Yandex Plus",
       "price": 450,
       "start_date": "07-2025"
     }' | jq
```

## 4. Получить все подписки

```bash
curl -b cookies.txt \
     http://localhost:8080/subscriptions | jq
```

## 5. Получить подписку по id

```bash
curl -b cookies.txt \
     http://localhost:8080/subscriptions/1 | jq
```

## 6. Обновить подписку

```bash
curl -X PUT http://localhost:8080/subscriptions/3 \
     -H "Content-Type: application/json" \
     -b cookies.txt \
     -d '{
       "service_name": "Yandex Plus",
       "price": 450,
       "start_date": "07-2025",
       "end_date": "12-2026"
     }' | jq
```

## 7. Удалить подписку

```bash
curl -X DELETE http://localhost:8080/subscriptions/1 \
     -b cookies.txt | jq
```

## 8. Получить сумму по фильтрам

```bash
SERVICE="Yandex%20Plus"
START_MONTH="07-2025"
END_MONTH="07-2025"

curl -b cookies.txt \
  "http://localhost:8080/subscriptions/total?service=${SERVICE}&start_month=${START_MONTH}&end_month=${END_MONTH}" | jq
```