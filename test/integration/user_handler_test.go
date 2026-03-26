package integration

import (
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require" // Gunakan require untuk fail-fast
	"github.com/stretchr/testify/suite"
	"github.com/yourusername/go_server/internal/handlers/users"
)

type UsersIntegrationTestSuite struct {
	suite.Suite
	app  *fiber.App
	pool *pgxpool.Pool
}

func (s *UsersIntegrationTestSuite) SetupSuite() {
	s.app = fiber.New()
	// Langsung pakai testDBPool dari main_test.go
	s.pool = testDBPool

	s.app.Get("/api/v1/users", func(c *fiber.Ctx) error {
		c.Locals("page", 1)
		c.Locals("per_page", 10)
		c.Locals("db", s.pool)
		return users.GetAll(c)
	})
}

func (s *UsersIntegrationTestSuite) TearDownTest() {
	// FIX: Cek error truncate secara eksplisit atau gunakan blank identifier yang jelas
	_, err := s.pool.Exec(context.Background(), "TRUNCATE users RESTART IDENTITY CASCADE")
	require.NoError(s.T(), err, "Gagal membersihkan database")
}

func (s *UsersIntegrationTestSuite) Test_GetAll_Success() {
	_, err := s.pool.Exec(context.Background(),
		"INSERT INTO users (email, username, password_hash) VALUES ($1, $2, $3)",
		"budiawan@fedora.id", "budiawan", "securehash")
	require.NoError(s.T(), err)

	req := httptest.NewRequest("GET", "/api/v1/users", nil)

	// FIX: Tangkap error dari app.Test
	resp, err := s.app.Test(req)
	require.NoError(s.T(), err)
	defer func() { _ = resp.Body.Close() }() // Bullet-proof: Tutup body

	s.Equal(200, resp.StatusCode)

	// FIX: Tangkap error dari io.ReadAll
	body, err := io.ReadAll(resp.Body)
	require.NoError(s.T(), err)

	var res map[string]interface{}

	// FIX UTAMA: Cek hasil return dari json.Unmarshal (Sesuai laporan linter)
	err = json.Unmarshal(body, &res)
	require.NoError(s.T(), err, "Gagal unmarshal JSON response")

	s.Equal(true, res["success"])
	s.NotNil(res["data"])
}

func TestUsersIntegration(t *testing.T) {
	suite.Run(t, new(UsersIntegrationTestSuite))
}
