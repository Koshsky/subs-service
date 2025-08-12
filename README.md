# Subscription Service

Микросервисная архитектура для управления подписками пользователей, построенная на Go с использованием gRPC, PostgreSQL и Docker.

## 🚀 Обзор проекта

Этот проект демонстрирует современные подходы к разработке backend-систем и включает в себя:

- **Микросервисную архитектуру** с разделением на auth-service и core-service
- **Отдельные базы данных** для каждого сервиса (Database per Service pattern)
- **gRPC API** для высокопроизводительного межсервисного взаимодействия
- **JWT-аутентификацию** для безопасности
- **PostgreSQL** с миграциями базы данных
- **Docker Compose** для оркестрации сервисов
- **Health checks** для мониторинга состояния сервисов

## 🏗️ Архитектура

```
┌─────────────────┐    ┌─────────────────┐
│   Auth Service  │    │   Core Service  │
│   (gRPC:50051)  │◄──►│   (HTTP:8080)   │
└─────────┬───────┘    └─────────┬───────┘
          │                      │
          ▼                      ▼
┌─────────────────┐    ┌─────────────────┐
│   Auth DB       │    │   Core DB       │
│ (users, auth)   │    │ (subscriptions) │
│ Port: 5433      │    │ Port: 5434      │
└─────────────────┘    └─────────────────┘
```

### Сервисы

- **Auth Service**: Аутентификация пользователей, JWT токены, gRPC API
- **Core Service**: Основная бизнес-логика управления подписками, REST API
- **Auth Database**: PostgreSQL база для пользователей и аутентификации
- **Core Database**: PostgreSQL база для подписок и бизнес-данных

## 🛠️ Технологический стек

### Backend
- **Go 1.24** - основной язык программирования
- **gRPC** - межсервисное взаимодействие
- **PostgreSQL** - отдельные базы данных для каждого сервиса
- **JWT** - аутентификация и авторизация
- **Docker & Docker Compose** - контейнеризация

### DevOps
- **Dockerfile** - мультистадийная сборка
- **Docker Compose** - оркестрация сервисов
- **Health Checks** - мониторинг состояния
- **Database Migrations** - автоматическое управление схемой БД через контейнеры

## 🚀 Быстрый старт

### Предварительные требования
- Docker & Docker Compose
- Git

### Запуск проекта

1. **Клонирование репозитория**
```bash
git clone <repository-url>
cd subs-service
```

2. **Настройка переменных окружения**
```bash
cp .env.example .env
# Отредактируйте .env файл под ваши нужды
```

3. **Запуск всех сервисов**
```bash
docker-compose up -d
```

4. **Проверка состояния сервисов**
```bash
docker-compose ps
```

### Доступные endpoints

- **Auth Service Health**: http://localhost:8081/health
- **Core Service Health**: http://localhost:8080/health
- **Core Service API**: http://localhost:8080/api/v1/

### Базы данных

- **Auth Database**: localhost:5433 (auth_user/auth_pass/auth_db)
- **Core Database**: localhost:5434 (core_user/core_pass/core_db)

## 📁 Структура проекта

```
.
├── auth-service/           # Сервис аутентификации
│   ├── cmd/               # Точки входа
│   ├── internal/          # Внутренняя логика
│   │   ├── models/        # Модели данных auth-service
│   │   ├── db/            # Подключение к auth БД
│   │   └── ...
│   ├── migrations/        # SQL миграции для auth БД
│   └── Dockerfile         # Контейнер auth-service
├── core-service/          # Основной сервис
│   ├── cmd/               # Точки входа
│   ├── internal/          # Внутренняя логика
│   │   ├── models/        # Модели данных core-service
│   │   ├── db/            # Подключение к core БД
│   │   └── ...
│   ├── migrations/        # SQL миграции для core БД
│   └── Dockerfile         # Контейнер core-service
├── docs/                  # Документация
├── docker-compose.yaml    # Оркестрация сервисов
└── README.md             # Этот файл
```

## 🔧 Разработка

### Локальная разработка

1. **Установка зависимостей**
```bash
cd auth-service && go mod download
cd ../core-service && go mod download
```

2. **Запуск отдельного сервиса**
```bash
# Только auth-service
docker-compose up -d auth-db auth-migrator
docker-compose up auth-service

# Только core-service
docker-compose up -d core-db core-migrator auth-service
docker-compose up core-service
```

### Миграции базы данных

Миграции выполняются автоматически при запуске через Docker Compose. Для ручного управления:

```bash
# Применить миграции auth-service
docker-compose up auth-migrator

# Применить миграции core-service
docker-compose up core-migrator

# Откатить миграции
# Auth service (используйте переменные из .env)
docker-compose run --rm auth-migrator -path=/migrations -database="postgres://${AUTH_DB_USER}:${AUTH_DB_PASSWORD}@auth-db:5432/${AUTH_DB_NAME}?sslmode=${AUTH_DB_SSLMODE}" down

# Core service (используйте переменные из .env)
docker-compose run --rm core-migrator -path=/migrations -database="postgres://${CORE_DB_USER}:${CORE_DB_PASSWORD}@core-db:5432/${CORE_DB_NAME}?sslmode=${CORE_DB_SSLMODE}" down
```

## 🧪 Тестирование

### Health Checks
```bash
# Проверка auth-service
curl http://localhost:8081/health

# Проверка core-service
curl http://localhost:8080/health
```

### API тестирование
```bash
# Пример регистрации пользователя (если реализовано)
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

## 🔐 Безопасность

- JWT токены для аутентификации
- Переменные окружения для конфиденциальных данных
- Непривилегированные пользователи в контейнерах
- Изолированные базы данных для каждого сервиса
- Health checks для мониторинга

## 📊 Мониторинг

Проект включает health checks для всех сервисов:
- Database health check с pg_isready для каждой БД
- HTTP health endpoints для сервисов
- Настроенные retry политики

## 🚀 Деплой

Проект готов к деплою в любой среде, поддерживающей Docker:
- Kubernetes
- Docker Swarm
- Cloud platforms (AWS, GCP, Azure)

## 📝 Лицензия

Этот проект создан в учебных целях и для демонстрации навыков разработки.

---

## 💡 Технические решения

### Почему микросервисы?
- Демонстрация навыков проектирования распределенных систем
- Возможность независимого масштабирования компонентов
- Разделение ответственности между сервисами

### Почему отдельные БД?
- **Database per Service pattern** - каждый сервис владеет своими данными
- Независимое масштабирование баз данных
- Изоляция данных и безопасность
- Возможность использования разных типов БД для разных сервисов

### Почему gRPC?
- Высокая производительность бинарного протокола
- Строгая типизация через Protocol Buffers
- Встроенная поддержка streaming

### Почему PostgreSQL?
- ACID-совместимость для критически важных данных
- Мощные возможности для сложных запросов
- Надежность и производительность

### Почему контейнер-миграторы?
- Отдельный этап миграций в CI/CD pipeline
- Безопасность - миграции не влияют на работу сервисов
- Возможность отката миграций
- Production-ready подход

Этот проект демонстрирует знание современных практик разработки backend-систем и готовность к работе в production-окружении.