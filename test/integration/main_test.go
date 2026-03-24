package integration

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Variable global agar bisa diakses semua file test di package ini
var testDBPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Setup database satu kali
	pool, cleanup := SetupTestContainer(ctx)
	testDBPool = pool

	// Jalankan semua test
	code := m.Run()

	// Hapus container setelah selesai
	cleanup()

	os.Exit(code)
}
