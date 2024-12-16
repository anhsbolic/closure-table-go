package main

import (
	"fmt"
	"github.com/anhsbolic/closure-table-go/config"
	"github.com/anhsbolic/closure-table-go/middleware"
	"github.com/anhsbolic/closure-table-go/pkg"
	"github.com/anhsbolic/closure-table-go/routes"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"time"
)

func main() {
	// Get Config
	env := config.GetEnvConfig()

	// Setup Server
	addr := fmt.Sprintf(":%s", env.Get("APP_PORT"))
	server := fiber.New(fiber.Config{
		IdleTimeout:  time.Second * 30,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
		Prefork:      true,
		ErrorHandler: pkg.NewErrorHandler,
	})

	// Set Global Middleware
	server.Use(middleware.XApiKeyMiddleware)

	// Setup DB
	db := pkg.NewDB()

	// Setup Redis
	// redisClient := pkg.NewRedisClient()

	// Setup Validator
	validate := validator.New()

	// Setup Routes
	routes.InitNodeRoutes(server, db, validate)

	// Start Server
	err := server.Listen(addr)
	pkg.PanicIfError(err)
}
