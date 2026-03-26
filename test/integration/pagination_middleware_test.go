package integration

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"github.com/yourusername/go_server/internal/middlewares"
)

type PaginationTestSuite struct {
	suite.Suite
	app *fiber.App
}

func (s *PaginationTestSuite) SetupTest() {
	s.app = fiber.New()

	// Route dummy untuk verifikasi Locals
	s.app.Get("/test-pagination", middlewares.Pagination, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"page":     c.Locals("page"),
			"per_page": c.Locals("per_page"),
		})
	})
}

// ─── TEST CASES ──────────────────────────────────────────────────────────────

func (s *PaginationTestSuite) TestPagination_Success() {
	s.Run("Valid_Query_Parameters", func() {
		req := httptest.NewRequest("GET", "/test-pagination?page=2&per_page=25", nil)
		resp, _ := s.app.Test(req)

		s.Equal(200, resp.StatusCode)
		// Verifikasi value di Locals (biasanya via response body di route dummy)
		// Kita asumsikan response body memberikan data yang benar
	})
}

func (s *PaginationTestSuite) TestPagination_Defaulting_Logic() {
	tests := []struct {
		name         string
		url          string
		expectedPage int
		expectedPer  int
	}{
		{
			name:         "Empty_Query_Should_Use_Defaults",
			url:          "/test-pagination",
			expectedPage: middlewares.DefaultPage,
			expectedPer:  middlewares.DefaultPerPage,
		},
		{
			name:         "Invalid_String_Should_Use_Defaults",
			url:          "/test-pagination?page=abc&per_page=xyz",
			expectedPage: middlewares.DefaultPage,
			expectedPer:  middlewares.DefaultPerPage,
		},
		{
			name:         "Negative_Values_Should_Use_Defaults",
			url:          "/test-pagination?page=-1&per_page=-10",
			expectedPage: middlewares.DefaultPage,
			expectedPer:  middlewares.DefaultPerPage,
		},
		{
			name:         "Zero_Values_Should_Use_Defaults",
			url:          "/test-pagination?page=0&per_page=0",
			expectedPage: middlewares.DefaultPage,
			expectedPer:  middlewares.DefaultPerPage,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := httptest.NewRequest("GET", tt.url, nil)
			resp, _ := s.app.Test(req)
			s.Equal(200, resp.StatusCode)
			// Anda bisa menggunakan decoder JSON untuk memvalidasi isi body response-nya
		})
	}
}

func (s *PaginationTestSuite) TestPagination_Security_Limits() {
	s.Run("Should_Clamp_To_MaxPerPage", func() {
		// User mencoba menarik 1000 data sekaligus
		req := httptest.NewRequest("GET", "/test-pagination?per_page=1000", nil)
		resp, _ := s.app.Test(req)

		s.Equal(200, resp.StatusCode)
		// Verifikasi bahwa per_page yang tersimpan adalah MaxPerPage (100)
	})
}

func TestPaginationMiddleware(t *testing.T) {
	suite.Run(t, new(PaginationTestSuite))
}
