package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB() {
	// 1. Ambil DSN dari environment
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL tidak ditemukan di environment")
	}

	// 2. Parse konfigurasi
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Gagal parsing config database: %v", err)
	}

	// 3. Optimasi Pool untuk Scalability & Performance
	config.MaxConns = 25                        // Maksimal koneksi simultan
	config.MinConns = 5                         // Koneksi yang selalu standby
	config.MaxConnIdleTime = 5 * time.Minute    // Tutup koneksi jika tidak dipakai
	config.HealthCheckPeriod = 30 * time.Second // Cek kesehatan koneksi berkala

	// 4. Inisialisasi Pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Gagal membuat pool database: %v", err)
	}

	// 5. Verifikasi koneksi (Ping)
	// Penting: NewWithConfig tidak langsung mencoba konek ke DB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Database tidak merespons: %v", err)
	}

	DB = pool
	fmt.Println("✅ Database Postgres Terkoneksi (via pgxpool)")
}
