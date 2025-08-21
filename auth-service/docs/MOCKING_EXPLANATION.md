# Моки в тестировании: объяснение и примеры

## 🎯 Для чего используется IUserRepository мок?

### 📊 Схема архитектуры

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   AuthService   │───▶│ UserRepository   │───▶│   Database      │
│   (Business     │    │   (Data Access)  │    │   (PostgreSQL)  │
│    Logic)       │    │                  │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

### 🧪 Тестирование БЕЗ моков (Integration Tests)

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   AuthService   │───▶│ UserRepository   │───▶│   Real DB       │
│   Test          │    │   (Real)         │    │   (Setup/Teardown)
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

**Проблемы:**
- 🐌 Медленно (настройка БД, миграции)
- 🔗 Зависит от внешних ресурсов
- 🛠️ Сложная настройка
- 💥 Хрупкие тесты

### 🚀 Тестирование С моками (Unit Tests)

```
┌─────────────────┐    ┌──────────────────┐
│   AuthService   │───▶│ Mock Repository  │
│   Test          │    │   (Fake)         │
└─────────────────┘    └──────────────────┘
```

**Преимущества:**
- ⚡ Быстро (нет реальных запросов)
- 🎯 Изолированно (только бизнес-логика)
- 🔧 Контролируемо (любые сценарии)
- 🛡️ Надежно (предсказуемые результаты)

## 🎭 Практические примеры использования

### 1. Тестирование успешной регистрации

```go
func TestRegisterSuccess(t *testing.T) {
    // Arrange
    mockRepo := new(mocks.IUserRepository)
    mockRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil)

    authService := &AuthService{UserRepo: mockRepo}

    // Act
    user, err := authService.Register(ctx, "test@example.com", "Password123!")

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "test@example.com", user.Email)
    mockRepo.AssertExpectations(t)
}
```

### 2. Тестирование ошибки базы данных

```go
func TestRegisterDatabaseError(t *testing.T) {
    // Arrange
    mockRepo := new(mocks.IUserRepository)
    mockRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(
        errors.New("connection failed"),
    )

    authService := &AuthService{UserRepo: mockRepo}

    // Act
    user, err := authService.Register(ctx, "test@example.com", "Password123!")

    // Assert
    assert.Error(t, err)
    assert.Nil(t, user)
    assert.Contains(t, err.Error(), "connection failed")
}
```

### 3. Тестирование валидации логина

```go
func TestLoginValidation(t *testing.T) {
    // Arrange
    mockRepo := new(mocks.IUserRepository)
    mockRepo.On("ValidateUser", "test@example.com", "Password123!").Return(
        &models.User{Email: "test@example.com"}, nil,
    )

    authService := &AuthService{UserRepo: mockRepo, JWTSecret: []byte("secret")}

    // Act
    token, user, err := authService.Login(ctx, "test@example.com", "Password123!")

    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, token)
    assert.Equal(t, "test@example.com", user.Email)
}
```

## 🔄 Разница между IDatabase и IUserRepository моками

### IDatabase мок
- **Используется для:** тестирования UserRepository
- **Что мокает:** низкоуровневые операции БД (Create, Where, First, Error)
- **Когда нужен:** когда тестируем логику репозитория

### IUserRepository мок
- **Используется для:** тестирования AuthService
- **Что мокает:** высокоуровневые операции (CreateUser, GetUserByEmail, ValidateUser)
- **Когда нужен:** когда тестируем бизнес-логику сервиса

## 🎯 Когда использовать моки?

### ✅ Используйте моки для:
- Unit тестов бизнес-логики
- Тестирования обработки ошибок
- Изоляции компонентов
- Быстрых тестов

### ❌ НЕ используйте моки для:
- Integration тестов
- Тестирования реальных интеграций
- End-to-end тестов

## 🏗️ Архитектурные преимущества

### 1. Dependency Injection
```go
type AuthService struct {
    UserRepo IUserRepository // Интерфейс вместо конкретного типа
}
```

### 2. Testability
```go
// В продакшене
authService := NewAuthService(realRepo, rabbitmq, secret)

// В тестах
authService := NewAuthService(mockRepo, nil, secret)
```

### 3. Loose Coupling
- Сервис не знает о деталях реализации репозитория
- Легко заменить реализацию
- Простое тестирование

## 📈 Результаты использования моков

| Метрика | Без моков | С моками |
|---------|-----------|----------|
| Время выполнения | 2-5 секунд | 0.1-0.5 секунд |
| Настройка | Сложная | Простая |
| Надежность | Средняя | Высокая |
| Изоляция | Низкая | Высокая |
| Покрытие edge cases | Сложно | Легко |
