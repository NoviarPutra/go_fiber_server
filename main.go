package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
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

	// 4. Port & Strict Sanitization
	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "3000"
	}

	// Konversi ke Int untuk validasi
	portInt, err := strconv.Atoi(portStr)
	if err != nil {
		// FIX G706: JANGAN mencetak portStr (input kotor) ke dalam log.
		// Cukup beri tahu bahwa port tidak valid dan kita menggunakan default.
		log.Println("⚠️  Invalid port detected in environment, falling back to 3000")
		portInt = 3000
	}

	// Gunakan strconv.Itoa untuk membuat string baru yang dijamin "bersih"
	cleanPort := strconv.Itoa(portInt)

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
	// Sekarang baris ini sudah aman karena cleanPort berasal dari int
	log.Printf("🚀 Server starting on port %s", cleanPort)

	if err := app.Listen(":" + cleanPort); err != nil {
		log.Printf("❌ Server stop reason: %v", err)
		select {
		case <-idleConnsClosed:
		default:
			close(idleConnsClosed)
		}
	}

	<-idleConnsClosed
	log.Println("✅ Server cleanup completed. Goodbye!")
}
