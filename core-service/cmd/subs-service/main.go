package main

import (
	"log"

	"github.com/Koshsky/subs-service/core-service/internal/config"
	"github.com/Koshsky/subs-service/core-service/internal/controllers"
	"github.com/Koshsky/subs-service/core-service/internal/repositories"
	"github.com/Koshsky/subs-service/core-service/internal/router"
	"github.com/Koshsky/subs-service/core-service/internal/services"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.LoadConfig()

	// Подключаемся к базе данных
	// Миграции выполняются автоматически контейнерами
	database, err := cfg.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to core database: %v", err)
	}
	log.Println("Core database connection established successfully")

	// Создаем gRPC клиент для auth-service
	authClient, err := services.NewAuthClient(cfg.AuthServiceAddr, cfg.EnableTLS, cfg.TLSCertFile)
	if err != nil {
		log.Fatalf("Failed to connect to auth service: %v", err)
	}
	defer authClient.Close()

	// Создаем репозитории
	subRepo := repositories.NewSubscriptionRepository(database)

	// Создаем сервисы
	subService := services.NewSubscriptionService(subRepo)

	// Создаем контроллеры
	authController := controllers.NewAuthController(authClient)
	subController := controllers.NewSubscriptionController(subService)

	// Настраиваем роутер
	r := router.SetupRouter(authController, subController, authClient)

	log.Printf("Core service started on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
