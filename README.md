# Subscription Service

Микросервисная архитектура для управления подписками пользователей, построенная на Go с использованием gRPC, PostgreSQL и Docker.

## 🚀 Обзор проекта

Этот проект демонстрирует современные подходы к разработке backend-систем и включает в себя:

- **Микросервисную архитектуру** с разделением на auth-service и core-service
- **gRPC API** для высокопроизводительного межсервисного взаимодействия
- **JWT-аутентификацию** для безопасности
- **PostgreSQL** с миграциями базы данных
- **Docker Compose** для оркестрации сервисов
- **Health checks** для мониторинга состояния сервисов

## 🏗️ Архитектура

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Auth Service  │    │   Core Service  │    │   PostgreSQL    │
│   (gRPC:50051)  │◄──►│   (HTTP:8080)   │◄──►│   (Port:5432)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Сервисы

- **Auth Service**: Аутентификация пользователей, JWT токены, gRPC API
- **Core Service**: Основная бизнес-логика управления подписками, REST API
- **Database**: PostgreSQL с автоматическими миграциями

## 🛠️ Технологический стек

### Backend
- **Go 1.24** - основной язык программирования
- **gRPC** - межсервисное взаимодействие
- **PostgreSQL** - основная база данных
- **JWT** - аутентификация и авторизация
- **Docker & Docker Compose** - контейнеризация

### DevOps
- **Dockerfile** - мультистадийная сборка
- **Docker Compose** - оркестрация сервисов
- **Health Checks** - мониторинг состояния
- **Database Migrations** - автоматическое управление схемой БД

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

## 📁 Структура проекта

```
.
├── auth-service/           # Сервис аутентификации
│   ├── cmd/               # Точки входа
│   ├── internal/          # Внутренняя логика
│   ├── proto/             # gRPC протоколы
│   └── Dockerfile         # Контейнер auth-service
├── core-service/          # Основной сервис
│   ├── cmd/               # Точки входа
│   ├── internal/          # Внутренняя логика
│   └── Dockerfile         # Контейнер core-service
├── shared/                # Общие библиотеки
├── migrations/            # SQL миграции
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
docker-compose up -d db migrator
docker-compose up auth-service

# Только core-service
docker-compose up -d db migrator auth-service
docker-compose up core-service
```

### Миграции базы данных

Миграции выполняются автоматически при запуске через Docker Compose. Для ручного управления:

```bash
# Применить миграции
docker-compose up migrator

# Откатить миграции
docker-compose run --rm migrator -path=/migrations -database="postgres://..." down
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
- Health checks для мониторинга

## 📊 Мониторинг

Проект включает health checks для всех сервисов:
- Database health check с pg_isready
- HTTP health endpoints для сервисов
- Настроенные retry политики

## 🚀 Деплой

Проект готов к деплою в любой среде, поддерживающей Docker:
- Kubernetes
- Docker Swarm
- Cloud platforms (AWS, GCP, Azure)

## 🤝 Вклад в проект

1. Fork проекта
2. Создайте feature branch (`git checkout -b feature/amazing-feature`)
3. Commit изменения (`git commit -m 'Add amazing feature'`)
4. Push в branch (`git push origin feature/amazing-feature`)
5. Откройте Pull Request

## 📝 Лицензия

Этот проект создан в учебных целях и для демонстрации навыков разработки.

## 👤 Автор

**Ваше имя** - демонстрационный проект для портфолио

---

## 💡 Технические решения

### Почему микросервисы?
- Демонстрация навыков проектирования распределенных систем
- Возможность независимого масштабирования компонентов
- Разделение ответственности между сервисами

### Почему gRPC?
- Высокая производительность бинарного протокола
- Строгая типизация через Protocol Buffers
- Встроенная поддержка streaming

### Почему PostgreSQL?
- ACID-совместимость для критически важных данных
- Мощные возможности для сложных запросов
- Надежность и производительность

Этот проект демонстрирует знание современных практик разработки backend-систем и готовность к работе в production-окружении.