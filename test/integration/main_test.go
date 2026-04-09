package integration

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Cukup deklarasikan di sini saja
var testDBPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Set environment variables for tests
	os.Setenv("JWT_SECRET", "test-secret-key-12345")

	// Setup database satu kali untuk semua file di package ini
	pool, cleanup := SetupTestContainer(ctx)
	testDBPool = pool

	code := m.Run()

	cleanup()
	os.Exit(code)
}
