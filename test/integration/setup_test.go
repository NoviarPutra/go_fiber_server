package integration

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	pg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupTestContainer(ctx context.Context) (*pgxpool.Pool, func()) {
	// 1. Jalankan Container Postgres menggunakan pg.Run (API Terbaru)
	container, err := pg.Run(ctx,
		"postgres:16-alpine",
		pg.WithDatabase("go_server_test"),
		pg.WithUsername("postgres"),
		pg.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		log.Fatalf("❌ Failed to start container: %v", err)
	}

	// Ambil Connection String dari instance container
	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("❌ Failed to get connection string: %v", err)
	}

	// 2. Jalankan Migrasi
	runGoMigrations(connStr)

	// 3. Buat pgxpool untuk Fiber
	dbPool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("❌ Failed to create pgxpool: %v", err)
	}

	// Cleanup
	cleanup := func() {
		dbPool.Close()
		if err := container.Terminate(ctx); err != nil {
			log.Printf("⚠️ Failed to terminate container: %v", err)
		}
	}

	return dbPool, cleanup
}

// runGoMigrations tetap sama, pastikan path migrations benar
func runGoMigrations(connStr string) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("❌ Migration connect error: %v", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("❌ Migration driver error: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("❌ Migration init error: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("❌ Migration Up failed: %v", err)
	}
	log.Println("✅ Database migrated successfully")
}
