package main

import (
	"log"
	"os"

	"Go-Exercise/config"
	"Go-Exercise/handler"
	"Go-Exercise/pkg/repository"
	"Go-Exercise/pkg/service"

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

	db := config.InitDB()

	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo, os.Getenv("APP_SECRET_KEY"))
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(authService, authRepo)

	app.Post("/auth/register", authHandler.Register)
	app.Post("/auth/login", authHandler.Login)
	app.Post("/auth/refresh", authHandler.RefreshToken)

	app.Post("/users", userHandler.AuthMiddleware, userHandler.CreateUser)
	app.Get("/users", userHandler.AuthMiddleware, userHandler.GetAllUsers)
	app.Get("/users/:id", userHandler.AuthMiddleware, userHandler.GetUser)
	app.Put("/users/:id", userHandler.AuthMiddleware, userHandler.UpdateUser)
	app.Delete("/users/:id", userHandler.AuthMiddleware, userHandler.DeleteUser)
	app.Get("/auth/me", userHandler.AuthMiddleware, userHandler.GetCurrentUser)
	app.Post("/auth/logout", userHandler.AuthMiddleware, authHandler.Logout)

	log.Fatal(app.Listen(":3002"))
}
