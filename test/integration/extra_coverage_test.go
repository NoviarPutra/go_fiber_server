package integration

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/yourusername/go_server/internal/handlers/companies"
	"github.com/yourusername/go_server/internal/types"
	"github.com/yourusername/go_server/internal/utils"
)

func TestExtraCoverage(t *testing.T) {
	t.Run("BuildMeta_Basic", func(t *testing.T) {
		meta := types.BuildMeta(1, 10, 100, "req-123")
		require.Equal(t, 1, meta.Page)
		require.Equal(t, 10, meta.PerPage)
		require.Equal(t, int64(100), meta.Total)
		require.Equal(t, "req-123", meta.RequestID)
	})

	t.Run("Forbidden_Response", func(t *testing.T) {
		app := fiber.New()
		app.Get("/forbidden", func(c *fiber.Ctx) error {
			return utils.Forbidden(c, "access denied")
		})

		req := httptest.NewRequest("GET", "/forbidden", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusForbidden, resp.StatusCode)
	})

	t.Run("InternalError_Response", func(t *testing.T) {
		app := fiber.New()
		app.Get("/error", func(c *fiber.Ctx) error {
			return utils.InternalError(c, "something went wrong")
		})

		req := httptest.NewRequest("GET", "/error", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("CheckPasswordHash_InvalidFormat", func(t *testing.T) {
		match, err := utils.CheckPasswordHash("pass", "invalid-hash-format")
		require.False(t, match)
		require.Error(t, err)
	})

	t.Run("SendResponse_Variety", func(t *testing.T) {
		app := fiber.New()
		app.Get("/var", func(c *fiber.Ctx) error {
			return utils.SendResponse(c, 202, true, "data", "msg", &types.Meta{Page: 1})
		})

		req := httptest.NewRequest("GET", "/var", nil)
		resp, _ := app.Test(req)
		require.Equal(t, 202, resp.StatusCode)
	})
	
	t.Run("SuccessWithMeta_Coverage", func(t *testing.T) {
		app := fiber.New()
		app.Get("/meta", func(c *fiber.Ctx) error {
			return utils.SuccessWithMeta(c, "data", "msg", &types.Meta{Page: 1})
		})

		req := httptest.NewRequest("GET", "/meta", nil)
		resp, _ := app.Test(req)
		require.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Unauthorized_Response", func(t *testing.T) {
		app := fiber.New()
		app.Get("/unauth", func(c *fiber.Ctx) error {
			return utils.Unauthorized(c, "unauth")
		})

		req := httptest.NewRequest("GET", "/unauth", nil)
		resp, _ := app.Test(req)
		require.Equal(t, 401, resp.StatusCode)
	})

	t.Run("BadRequest_Response", func(t *testing.T) {
		app := fiber.New()
		app.Get("/bad", func(c *fiber.Ctx) error {
			return utils.BadRequest(c, "bad")
		})

		req := httptest.NewRequest("GET", "/bad", nil)
		resp, _ := app.Test(req)
		require.Equal(t, 400, resp.StatusCode)
	})

	t.Run("NotFound_Response", func(t *testing.T) {
		app := fiber.New()
		app.Get("/notfound", func(c *fiber.Ctx) error {
			return utils.NotFound(c, "notfound")
		})

		req := httptest.NewRequest("GET", "/notfound", nil)
		resp, _ := app.Test(req)
		require.Equal(t, 404, resp.StatusCode)
	})

	t.Run("Conflict_Response", func(t *testing.T) {
		app := fiber.New()
		app.Get("/conflict", func(c *fiber.Ctx) error {
			return utils.Conflict(c, "conflict")
		})

		req := httptest.NewRequest("GET", "/conflict", nil)
		resp, _ := app.Test(req)
		require.Equal(t, 409, resp.StatusCode)
	})

	t.Run("Created_Response", func(t *testing.T) {
		app := fiber.New()
		app.Get("/created", func(c *fiber.Ctx) error {
			return utils.Created(c, "data", "msg")
		})

		req := httptest.NewRequest("GET", "/created", nil)
		resp, _ := app.Test(req)
		require.Equal(t, 201, resp.StatusCode)
	})

	t.Run("NoContent_Response", func(t *testing.T) {
		app := fiber.New()
		app.Get("/nocontent", func(c *fiber.Ctx) error {
			return utils.NoContent(c)
		})

		req := httptest.NewRequest("GET", "/nocontent", nil)
		resp, _ := app.Test(req)
		require.Equal(t, 204, resp.StatusCode)
	})

	t.Run("GetAll_NoDB_In_Context", func(t *testing.T) {
		app := fiber.New()
		app.Get("/nodes", companies.GetAll)
		req := httptest.NewRequest("GET", "/nodes", nil)
		resp, _ := app.Test(req)
		require.Equal(t, 500, resp.StatusCode)
	})
}
