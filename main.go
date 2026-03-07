package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/yourusername/go_server/handlers/errors"
	"github.com/yourusername/go_server/routes"
)

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: errors.GlobalErrorHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Prefork:      false,
	})

	// Middleware
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	app.Use(logger.New())
	app.Use(cors.New())

	// Routes
	routes.SetupRoutes(app)

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down server...")
		_ = app.Shutdown()
	}()

	if err := app.Listen(":" + port); err != nil {
		log.Panic(err)
	}

	log.Fatal(app.Listen(":" + port))
}
