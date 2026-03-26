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

	// Setup database satu kali untuk semua file di package ini
	pool, cleanup := SetupTestContainer(ctx)
	testDBPool = pool

	code := m.Run()

	cleanup()
	os.Exit(code)
}
