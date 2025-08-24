# gRPC - Документация

## Обзор

В данной системе **gRPC** используется для синхронного взаимодействия между **Auth Service** и **Core Service**. gRPC обеспечивает высокопроизводительное межсервисное взаимодействие с использованием Protocol Buffers и HTTP/2.

## Архитектура gRPC

### Основные компоненты

```
┌─────────────────┐    ┌─────────────────┐
│   Auth Service  │◄──►│   Core Service  │
│   (gRPC Server) │    │ (gRPC Client)   │
│   Port: 50051   │    │   Port: 8080    │
└─────────────────┘    └─────────────────┘
```

### Тип взаимодействия: Unary RPC

В системе используется **Unary RPC** - простейший тип gRPC вызовов, где клиент отправляет один запрос и получает один ответ.

## Protocol Buffers (protobuf)

### Определение сервиса

```protobuf
syntax = "proto3";

package authpb;

option go_package = "github.com/Koshsky/subs-service/auth-service/internal/authpb";

// Authentication service
service AuthService {
  // Token validation and user information retrieval
  rpc ValidateToken(TokenRequest) returns (UserResponse);

  // New user registration
  rpc Register(RegisterRequest) returns (RegisterResponse);

  // User login
  rpc Login(LoginRequest) returns (LoginResponse);
}
```

### Структуры сообщений

```protobuf
// Token validation request
message TokenRequest {
  string token = 1;
}

// Response with user information
message UserResponse {
  string user_id = 1;
  string email = 2;
  bool valid = 3;
  string error = 4;
}

// Request for user registration
message RegisterRequest {
  string email = 1;
  string password = 2;
}

// Response for user registration
message RegisterResponse {
  string user_id = 1;
  string email = 2;
  bool success = 3;
  string error = 4;
  string message = 5;
}

// Login request
message LoginRequest {
  string email = 1;
  string password = 2;
}

// Login response
message LoginResponse {
  string token = 1;
  string user_id = 2;
  string email = 3;
  bool success = 4;
  string error = 5;
  string message = 6;
}
```

### Особенности определения

1. **Семантические имена**: Все поля имеют понятные имена
2. **Единообразные ответы**: Все ответы содержат поля `success` и `error`
3. **Строгая типизация**: Все поля имеют явные типы
4. **Версионирование**: Используется proto3 для совместимости

## Реализация Server (Auth Service)

### Структура сервера

```go
type AuthServer struct {
    authpb.UnimplementedAuthServiceServer   // Встроенная структура для совместимости
    AuthService *services.AuthService       // Бизнес-логика
}

func NewAuthServer(authService *services.AuthService) *AuthServer {
    return &AuthServer{
        AuthService: authService,
    }
}
```

### Реализация методов

#### ValidateToken

```go
func (s *AuthServer) ValidateToken(ctx context.Context, req *authpb.TokenRequest) (*authpb.UserResponse, error) {
    claims, err := s.AuthService.ValidateToken(ctx, req.Token)
    if err != nil {
        return &authpb.UserResponse{
            Valid: false,
            Error: err.Error(),
        }, nil
    }

    userIDStr, ok := claims["user_id"].(string)
    if !ok {
        return &authpb.UserResponse{
            Valid: false,
            Error: "Invalid user ID in token",
        }, nil
    }

    email, ok := claims["email"].(string)
    if !ok {
        return &authpb.UserResponse{
            Valid: false,
            Error: "Invalid email in token",
        }, nil
    }

    return &authpb.UserResponse{
        UserId: userIDStr,
        Email:  email,
        Valid:  true,
    }, nil
}
```

#### Register

```go
func (s *AuthServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
    startTime := time.Now()
    log.Printf("[AUTH_SERVER] [%s] Starting Register gRPC handler for email: %s",
        startTime.Format("15:04:05.000"), req.Email)

    user, err := s.AuthService.Register(ctx, req.Email, req.Password)

    if err != nil {
        totalDuration := time.Since(startTime)
        log.Printf("[AUTH_SERVER] [%s] Register FAILED after %v (service error: %v)",
            time.Now().Format("15:04:05.000"), totalDuration, err)
        return &authpb.RegisterResponse{
            Success: false,
            Error:   err.Error(),
        }, nil
    }

    response := &authpb.RegisterResponse{
        UserId:  user.ID.String(),
        Email:   user.Email,
        Success: true,
        Message: "User created successfully",
    }

    totalDuration := time.Since(startTime)
    log.Printf("[AUTH_SERVER] [%s] Register SUCCESS in %v",
        time.Now().Format("15:04:05.000"), totalDuration)

    return response, nil
}
```

#### Login

```go
func (s *AuthServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
    token, user, err := s.AuthService.Login(ctx, req.Email, req.Password)
    if err != nil {
        return &authpb.LoginResponse{
            Success: false,
            Error:   err.Error(),
        }, nil
    }

    return &authpb.LoginResponse{
        Token:   token,
        UserId:  user.ID.String(),
        Email:   user.Email,
        Success: true,
        Message: "Successful login",
    }, nil
}
```

### Особенности реализации Server

1. **Детальное логирование**: Каждый вызов логируется с временными метками
2. **Обработка ошибок**: Все ошибки возвращаются в структурированном виде
3. **Type assertions**: Безопасное извлечение данных из JWT claims
4. **Context support**: Поддержка контекста для timeout и cancellation

## Реализация Client (Core Service)

### Структура клиента

```go
type AuthClient struct {
    client authpb.AuthServiceClient
    conn   *grpc.ClientConn
}

func NewAuthClient(target string) (*AuthClient, error) {
    conn, err := grpc.Dial(target, grpc.WithInsecure())
    if err != nil {
        return nil, fmt.Errorf("failed to connect to auth service: %v", err)
    }

    client := authpb.NewAuthServiceClient(conn)
    return &AuthClient{
        client: client,
        conn:   conn,
    }, nil
}
```

### Использование клиента

```go
func (a *AuthClient) ValidateToken(ctx context.Context, token string) (*authpb.UserResponse, error) {
    req := &authpb.TokenRequest{
        Token: token,
    }

    response, err := a.client.ValidateToken(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("failed to validate token: %v", err)
    }

    return response, nil
}

func (a *AuthClient) Close() {
    if a.conn != nil {
        a.conn.Close()
    }
}
```

## Инициализация gRPC сервера

### Основная функция

```go
func main() {
    cfg := config.LoadConfig()

    // ... инициализация базы данных и сервисов ...

    userRepo := repositories.NewUserRepository(database)
    authService := services.NewAuthService(userRepo, rabbitmqService, []byte(cfg.JWTSecret))
    authServer := server.NewAuthServer(authService)

    var grpcServer *grpc.Server
    if cfg.EnableTLS {
        creds, err := credentials.NewServerTLSFromFile(cfg.TLSCertFile, cfg.TLSKeyFile)
        if err != nil {
            log.Fatalf("Failed to create TLS credentials: %v", err)
        }
        grpcServer = grpc.NewServer(grpc.Creds(creds))
        log.Printf("Auth service configured with TLS")
    } else {
        grpcServer = grpc.NewServer()
        log.Printf("Auth service configured without TLS (WARNING: Insecure)")
    }

    authpb.RegisterAuthServiceServer(grpcServer, authServer)

    lis, err := net.Listen("tcp", ":"+cfg.Port)
    if err != nil {
        log.Fatalf("Failed to listen on port %s: %v", cfg.Port, err)
    }

    log.Printf("Auth service starting on port %s", cfg.Port)
    if err := grpcServer.Serve(lis); err != nil {
        log.Printf("gRPC server stopped: %v", err)
    }
}
```

### Особенности инициализации

1. **TLS поддержка**: Условная настройка TLS в зависимости от конфигурации
2. **Graceful shutdown**: Правильное завершение работы сервера
3. **Error handling**: Детальная обработка ошибок инициализации
4. **Logging**: Подробное логирование состояния сервера

## Middleware и Interceptors

### Логирование Interceptor

```go
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    start := time.Now()
    
    log.Printf("gRPC call: %s", info.FullMethod)
    
    resp, err := handler(ctx, req)
    
    log.Printf("gRPC call %s completed in %v", info.FullMethod, time.Since(start))
    
    return resp, err
}
```

### Аутентификация Interceptor

```go
func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    // Пропускаем аутентификацию для методов регистрации и логина
    if info.FullMethod == "/authpb.AuthService/Register" ||
       info.FullMethod == "/authpb.AuthService/Login" {
        return handler(ctx, req)
    }

    // Проверяем токен для остальных методов
    // ... логика проверки токена ...

    return handler(ctx, req)
}
```

## Обработка ошибок

### gRPC Status Codes

```go
import (
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

func (s *AuthServer) ValidateToken(ctx context.Context, req *authpb.TokenRequest) (*authpb.UserResponse, error) {
    if req.Token == "" {
        return nil, status.Error(codes.InvalidArgument, "token is required")
    }

    claims, err := s.AuthService.ValidateToken(ctx, req.Token)
    if err != nil {
        return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
    }

    // ... остальная логика ...
}
```

### Кастомные ошибки

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for field %s: %s", e.Field, e.Message)
}

func (s *AuthServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
    if req.Email == "" {
        return nil, status.Error(codes.InvalidArgument, "email is required")
    }

    if req.Password == "" {
        return nil, status.Error(codes.InvalidArgument, "password is required")
    }

    // ... остальная логика ...
}
```

## Конфигурация

### Переменные окружения

```bash
# Auth Service
AUTH_SERVICE_PORT=50051
ENABLE_TLS=false
TLS_CERT_FILE=certs/server-cert.pem
TLS_KEY_FILE=certs/server-key.pem

# Core Service (Client)
AUTH_SERVICE_HOST=auth-service
AUTH_SERVICE_PORT=50051
```

### Docker Compose

```yaml
auth-service:
  build:
    context: .
    dockerfile: ./auth-service/Dockerfile
  ports:
    - "${AUTH_SERVICE_PORT}:${AUTH_SERVICE_PORT}"
  environment:
    - AUTH_SERVICE_PORT=${AUTH_SERVICE_PORT}
    - ENABLE_TLS=${ENABLE_TLS}
  networks:
    - subs_net

core-service:
  build:
    context: .
    dockerfile: ./core-service/Dockerfile
  depends_on:
    - auth-service
  environment:
    - AUTH_SERVICE_HOST=auth-service
    - AUTH_SERVICE_PORT=${AUTH_SERVICE_PORT}
  networks:
    - subs_net
```

## Мониторинг и отладка

### Health Checks

```go
import (
    "google.golang.org/grpc/health"
    "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
    // ... инициализация сервера ...

    healthServer := health.NewServer()
    grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
    healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

    // ... запуск сервера ...
}
```

### Метрики

```go
import (
    "github.com/grpc-ecosystem/go-grpc-prometheus"
    "github.com/prometheus/client_golang/prometheus"
)

func main() {
    // Регистрируем метрики
    grpc_prometheus.EnableHandlingTimeHistogram()

    // Создаем сервер с метриками
    grpcServer := grpc.NewServer(
        grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
    )

    // Регистрируем метрики в Prometheus
    prometheus.MustRegister(grpc_prometheus.DefaultServerMetrics)
}
```

## Лучшие практики

### 1. Обработка ошибок
- Используйте стандартные gRPC status codes
- Предоставляйте детальную информацию об ошибках
- Логируйте все ошибки для отладки

### 2. Производительность
- Используйте connection pooling для клиентов
- Настройте timeout для всех вызовов
- Мониторьте время выполнения методов

### 3. Безопасность
- Используйте TLS в production
- Реализуйте аутентификацию через interceptors
- Валидируйте все входящие данные

### 4. Мониторинг
- Добавьте health checks
- Настройте метрики для мониторинга
- Логируйте все важные операции

## Расширение API

### Добавление нового метода

1. **Обновить proto файл:**
```protobuf
service AuthService {
  // ... существующие методы ...

  // New method
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
}

message UpdateUserRequest {
  string user_id = 1;
  string email = 2;
}

message UpdateUserResponse {
  bool success = 1;
  string error = 2;
  string message = 3;
}
```

2. **Реализовать метод в сервере:**
```go
func (s *AuthServer) UpdateUser(ctx context.Context, req *authpb.UpdateUserRequest) (*authpb.UpdateUserResponse, error) {
    // Реализация метода
    return &authpb.UpdateUserResponse{
        Success: true,
        Message: "User updated successfully",
    }, nil
}
```

3. **Обновить клиент:**
```go
func (a *AuthClient) UpdateUser(ctx context.Context, userID, email string) (*authpb.UpdateUserResponse, error) {
    req := &authpb.UpdateUserRequest{
        UserId: userID,
        Email:  email,
    }

    return a.client.UpdateUser(ctx, req)
}
```

## Заключение

Реализация gRPC в данной системе обеспечивает:

- **Высокую производительность**: HTTP/2 и Protocol Buffers
- **Строгую типизацию**: compile-time проверки
- **Надежность**: встроенная обработка ошибок
- **Масштабируемость**: легко добавлять новые методы
- **Мониторинг**: встроенные health checks и метрики

Система готова к расширению и может легко адаптироваться к новым требованиям API.
