package main

import (
	"log"
	"maglo/config"
	"maglo/database"
	"maglo/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	database.ConnectDB()

	router.SetupRoutes(app)
	PORT := config.Config("PORT")
	if PORT == "" {
		PORT = "4000"
	}

	log.Fatal(app.Listen(":" + PORT))

}
