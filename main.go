package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	app "github.com/yourusername/go_server/internal"
	"github.com/yourusername/go_server/internal/config"
)

func main() {
	// 1. Load env
	if err := godotenv.Load(); err != nil {
		log.Println("Info: File .env tidak ditemukan, menggunakan env system")
	}

	// 2. Connect DB
	config.ConnectDB()

	// 3. Bootstrap app
	app := app.Bootstrap(config.DB)

	// 4. Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// 5. Graceful shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("⚠️  Shutting down server...")

		if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
			log.Printf("❌ Fiber shutdown error: %v", err)
		}

		log.Println("🗄️  Closing database connection pool...")
		config.DB.Close()

		close(idleConnsClosed)
	}()

	// 6. Start server
	log.Printf("🚀 Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Printf("❌ Server stop reason: %v", err)
		close(idleConnsClosed) // prevent deadlock jika port sudah dipakai
	}

	<-idleConnsClosed
	log.Println("✅ Server cleanup completed. Goodbye!")
}
