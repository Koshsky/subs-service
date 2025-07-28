package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Koshsky/subs-service/config"
	"github.com/Koshsky/subs-service/router"
	"github.com/Koshsky/subs-service/utils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	appConfig, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := sqlx.Connect("postgres", appConfig.DB.ConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := applyMigrations(appConfig.DB); err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	utils.RegisterCustomValidations()

	r := router.SetupRouter(db, appConfig.Router)

	log.Println("Starting server on :8080")
	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func applyMigrations(cfg *config.DBConfig) error {
	// Формируем DSN для миграций
	dsn := cfg.MigrationDSN()

	// Проверяем существует ли директория с миграциями
	if _, err := os.Stat("migrations"); os.IsNotExist(err) {
		log.Println("No migrations directory found, skipping migrations")
		return nil
	}

	m, err := migrate.New(
		"file://migrations",
		dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Migrations applied successfully")
	return nil
}
