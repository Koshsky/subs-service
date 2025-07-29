package main

import (
	"log"

	"github.com/Koshsky/subs-service/internal/config"
	"github.com/Koshsky/subs-service/internal/repositories"
	"github.com/Koshsky/subs-service/internal/repositories/db"
	"github.com/Koshsky/subs-service/internal/router"
	"github.com/Koshsky/subs-service/internal/utils"
	_ "github.com/lib/pq"
)

func main() {
	utils.RegisterCustomValidations()

	appConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database := db.ConnectDatabase(appConfig.DB)
	// defer database.Close()
	repo := repositories.New(database)

	r := router.SetupRouter(repo, appConfig.Router)

	log.Println("Starting server on :8080")
	r.Run()
}
