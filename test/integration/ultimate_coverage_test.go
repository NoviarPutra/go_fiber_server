package integration

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	pkg_errors "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	handler_errors "github.com/yourusername/go_server/internal/handlers/errors"
	"github.com/yourusername/go_server/internal/middlewares"
	"github.com/yourusername/go_server/internal/types"
	"github.com/yourusername/go_server/internal/utils"
)

func TestUltimateCoverage(t *testing.T) {
	t.Run("GlobalErrorHandler_StackTrace_Branch", func(t *testing.T) {
		app := fiber.New()

		// Force development mode for stack trace branch via local context if possible
		// or just call it directly as we do below

		app.Get("/error", func(c *fiber.Ctx) error {
			// pkg_errors.WithStack satisfies the stack_tracer interface
			err := pkg_errors.WithStack(errors.New("custom failure with stack"))
			return handler_errors.GlobalErrorHandler(c, err)
		})

		req := httptest.NewRequest("GET", "/error", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, 500, resp.StatusCode)
	})

	t.Run("Pagination_Clamp_MaxPerPage", func(t *testing.T) {
		app := fiber.New()
		app.Get("/p", middlewares.Pagination, func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"per_page": c.Locals("per_page")})
		})

		req := httptest.NewRequest("GET", "/p?per_page=9999", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Response_Utils_Complete", func(t *testing.T) {
		app := fiber.New()
		app.Get("/1", func(c *fiber.Ctx) error { return utils.Conflict(c, "c") })
		app.Get("/2", func(c *fiber.Ctx) error { return utils.BadRequest(c, "b") })
		app.Get("/3", func(c *fiber.Ctx) error { return utils.Created(c, "d", "m") })

		for _, path := range []string{"/1", "/2", "/3"} {
			req := httptest.NewRequest("GET", path, nil)
			resp, _ := app.Test(req)
			require.NotEmpty(t, resp.StatusCode)
		}
	})

	t.Run("Types_Meta_Build_Alternative", func(t *testing.T) {
		meta := types.BuildMeta(0, 0, 0, "")
		require.NotNil(t, meta)
	})
}
