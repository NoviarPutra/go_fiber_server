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
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL tidak ditemukan di environment")
	}

	pool_config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Gagal parsing config database: %v", err)
	}

	pool_config.MaxConns = 25
	pool_config.MinConns = 5
	pool_config.MaxConnIdleTime = 5 * time.Minute
	pool_config.MaxConnLifetime = 60 * time.Minute
	pool_config.HealthCheckPeriod = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(context.Background(), pool_config)
	if err != nil {
		log.Fatalf("Gagal membuat pool database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		log.Fatalf("Database tidak merespons: %v", err)
	}

	DB = pool

	stats := pool.Stat()
	fmt.Printf("✅ Database Postgres Terkoneksi | Pool: %d conns\n", stats.TotalConns())
}
