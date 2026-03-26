package utils

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestGetRequestIDSuite(t *testing.T) {
	is := assert.New(t)
	app := fiber.New()

	t.Run("Success: Should return correct request ID from Locals", func(t *testing.T) {
		ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
		defer app.ReleaseCtx(ctx)

		expectedID := "req-123-abc"
		ctx.Locals("requestid", expectedID)

		result := GetRequestID(ctx)

		is.Equal(expectedID, result)
	})

	t.Run("Edge Case: Should return 'unknown' when requestid is missing", func(t *testing.T) {
		ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
		defer app.ReleaseCtx(ctx)

		// Tidak set locals apa-apa
		result := GetRequestID(ctx)

		is.Equal("unknown", result)
	})

	t.Run("Edge Case: Should return 'unknown' when requestid is NOT a string", func(t *testing.T) {
		ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
		defer app.ReleaseCtx(ctx)

		// Set locals tapi tipenya integer (kasus salah input data)
		ctx.Locals("requestid", 999)

		result := GetRequestID(ctx)

		is.Equal("unknown", result, "Harus mengembalikan unknown jika tipe data di locals bukan string")
	})

	t.Run("Edge Case: Should return 'unknown' when requestid is nil", func(t *testing.T) {
		ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
		defer app.ReleaseCtx(ctx)

		ctx.Locals("requestid", nil)

		result := GetRequestID(ctx)

		is.Equal("unknown", result)
	})
}
