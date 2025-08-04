package db

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Драйвер для PostgreSQL
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Драйвер для файловых миграций
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var cfg *DBConfig

func init() {
	cfg = loadDBConfig()
}

// ConnectDatabase connects to the database
func ConnectDatabase(connectionString string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	if err := applyMigrations(); err != nil {
		log.Fatal("Failed to apply migrations: ", err)
	}

	return db
}

// applyMigrations applies the migrations to the database
func applyMigrations() error {
	if _, err := os.Stat("migrations"); os.IsNotExist(err) {
		log.Println("No migrations directory found, skipping migrations")
		return nil
	}

	m, err := migrate.New(
		"file://migrations",
		cfg.MigrationDSN())
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
