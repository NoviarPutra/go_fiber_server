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

	// Explicitly exit with error code if tests fail, but do not exit if success.
	// Actually, Go 1.15+ testing generates a main that uses os.Exit(m.Run()), 
	// but if we return, it automatically passes the exit code. We should just 
	// set os.Exit(code) only if code != 0, to let defer run for success.
	// Wait, TestMain return is implicitly handled? Yes!
	if code != 0 {
		os.Exit(code)
	}
}
