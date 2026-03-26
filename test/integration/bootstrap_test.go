package integration

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/yourusername/go_server/internal"
)

type BootstrapIntegrationTestSuite struct {
	suite.Suite
	// Tidak perlu cleanup func di sini karena dihandle TestMain
}

func (s *BootstrapIntegrationTestSuite) TestAppConfiguration() {
	// Gunakan testDBPool dari global variable di suite_test.go / TestMain
	s.Require().NotNil(testDBPool, "Database pool harus sudah terinisialisasi")

	app := internal.Bootstrap(testDBPool)

	s.Run("Fiber_Settings_Validation", func() {
		cfg := app.Config()
		s.Equal("Office Core API v1.0", cfg.AppName)
		s.Equal(60*time.Second, cfg.IdleTimeout)
		s.Equal(4*1024*1024, cfg.BodyLimit)
	})

	s.Run("Middleware_Execution_Integrity", func() {
		req := httptest.NewRequest("GET", "/health", nil)
		resp, err := app.Test(req, 5000) // 5 detik timeout cukup

		s.NoError(err)
		s.Equal(200, resp.StatusCode)
	})
}

func TestBootstrapIntegration(t *testing.T) {
	suite.Run(t, new(BootstrapIntegrationTestSuite))
}
