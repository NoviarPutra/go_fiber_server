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
	"github.com/joho/godotenv"
	"github.com/yourusername/go_server/config"
	"github.com/yourusername/go_server/handlers/errors"
	"github.com/yourusername/go_server/routes"
)

func main() {
	// 1. Load env & Database
	if err := godotenv.Load(); err != nil {
		log.Println("Info: File .env tidak ditemukan, menggunakan env system")
	}

	config.ConnectDB()

	app := fiber.New(fiber.Config{
		ErrorHandler: errors.GlobalErrorHandler,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Prefork:      false, // Biarkan false jika menggunakan pool database agar tidak konflik
	})

	// 2. Middleware Stack (Urutan Menentukan Performa)
	app.Use(recover.New()) // 1. Tangkap panic dulu
	app.Use(logger.New())  // 2. Catat log request
	app.Use(cors.New())    // 3. Atur akses cross-origin

	app.Use(limiter.New(limiter.Config{
		Max:        100, // Tingkatkan sedikit jika trafik tinggi
		Expiration: 1 * time.Minute,
	}))

	// 3. Routes
	routes.SetupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// 4. Graceful Shutdown Logic
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("⚠️ Shutting down server...")

		// Beri waktu 10 detik untuk menyelesaikan request yang ada
		if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
			log.Printf("❌ Fiber shutdown error: %v", err)
		}

		// Tutup koneksi database setelah server berhenti
		log.Println("🗄️ Closing database connection pool...")
		config.DB.Close()

		close(idleConnsClosed)
	}()

	// 5. Start Server
	log.Printf("🚀 Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Printf("❌ Server stop reason: %v", err)
	}

	<-idleConnsClosed
	log.Println("✅ Server cleanup completed. Goodbye!")
}
