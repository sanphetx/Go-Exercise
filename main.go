package main

import (
	"log"

	"go-exercise-1/config"
	"go-exercise-1/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Bangkok",
	}))

	config.ConnectDB()

	app.Post("/users", handlers.CreateUser)
	app.Get("/users", handlers.GetAllUsers)
	app.Get("/users/:id", handlers.GetUser)
	app.Put("/users/:id", handlers.UpdateUser)
	app.Delete("/users/:id", handlers.DeleteUser)

	log.Fatal(app.Listen(":3002"))
}
