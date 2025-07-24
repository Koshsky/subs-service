package main

import (
	"log"

	"github.com/Koshsky/subs-service/repositories/db"
	"github.com/Koshsky/subs-service/router"
	"github.com/Koshsky/subs-service/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	database := db.ConnectDatabase()
	// defer database.Close()

	utils.RegisterCustomValidations()

	r := gin.New()
	router.SetupRouter(r, database)

	log.Println("Starting server on :8080")
	r.Run()
}
