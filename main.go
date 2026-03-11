package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/yourusername/go_server/handlers/errors"
	"github.com/yourusername/go_server/routes"
)

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: errors.GlobalErrorHandler,
		IdleTimeout:  60 * time.Second,
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
	app.Use(limiter.New(limiter.Config{
		Max:        50,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	}))

	// Routes
	routes.SetupRoutes(app)

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Buat channel untuk mendengarkan sinyal terminasi
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		// Menangkap sinyal interrupt (Ctrl+C) dan terminate (Docker stop)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("⚠️ Shutting down server...")

		// ShutdownWithTimeout memastikan request yang sedang berjalan diselesaikan dulu
		if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
			log.Printf("❌ Shutdown error: %v", err)
		}

		// Tambahkan di sini: Tutup koneksi database jika ada (misal: db.Close())

		close(idleConnsClosed)
	}()

	log.Printf("🚀 Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("❌ Server failed to start: %v", err)
	}

	<-idleConnsClosed
	log.Println("✅ Server cleanup completed. Goodbye!")
}
