package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// ─── Format Constants ─────────────────────────────────────────────────────────

const production_log_format = "${time} | ${status} | ${latency_readable} | ${method} ${path} | ${error}\n"

const development_log_format = "" +
	"${cyan}${time}${reset} | ${status_colored} | ${yellow}${latency_readable}${reset} | " +
	"${blue}${method}${reset} ${white}${path}${reset}" +
	"${error_colored}\n"

// ─── Middleware ───────────────────────────────────────────────────────────────

func LoggerMiddleware() fiber.Handler {
	is_dev := os.Getenv("APP_ENV") == "development"

	format := production_log_format
	if is_dev {
		format = development_log_format
	}

	return logger.New(logger.Config{
		Format:     format,
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Jakarta",

		CustomTags: map[string]logger.LogFunc{

			// Status code asli dengan warna
			"status_colored": func(buf logger.Buffer, c *fiber.Ctx, data *logger.Data, _ string) (int, error) {
				status := c.Response().StatusCode()
				color := status_color(status)
				text := http_status_text(status)
				return buf.WriteString(fmt.Sprintf("%s%d %s\x1b[0m", color, status, text))
			},

			// Latency dalam ms + emoji berdasarkan kecepatan
			"latency_readable": func(buf logger.Buffer, c *fiber.Ctx, data *logger.Data, _ string) (int, error) {
				latency := data.Stop.Sub(data.Start)
				ms := float64(latency.Nanoseconds()) / 1e6

				var result string
				switch {
				case latency < 100*time.Millisecond:
					result = fmt.Sprintf("%.2fms ⚡", ms)
				case latency < time.Second:
					result = fmt.Sprintf("%.2fms ⚠️", ms)
				case latency < 10*time.Second:
					result = fmt.Sprintf("%.2fs 🐌", latency.Seconds())
				default:
					result = fmt.Sprintf("%.2fs 🔴", latency.Seconds())
				}

				return buf.WriteString(result)
			},

			// Error hanya muncul kalau ada
			"error_colored": func(buf logger.Buffer, c *fiber.Ctx, data *logger.Data, _ string) (int, error) {
				if data.ChainErr == nil {
					return 0, nil
				}
				return buf.WriteString(" | \x1b[31m⚠ " + data.ChainErr.Error() + "\x1b[0m")
			},
		},

		Done: func(c *fiber.Ctx, logString []byte) {
			// Jika terjadi server error (500+), tulis ke Stderr
			if c.Response().StatusCode() >= 500 {
				// FIX G104: Gunakan blank identifier untuk membuang error secara eksplisit
				// Ini memberitahu linter bahwa kita sadar ada error tapi memilih mengabaikannya
				_, _ = os.Stderr.Write(logString)
			}

			if is_dev {
				body := c.Response().Body()
				content_type := string(c.Response().Header.ContentType())

				if len(body) > 0 && containsStr(content_type, "application/json") {
					pretty := pretty_json(body)
					if pretty != "" {
						// FIX G104: Sama seperti di atas, buang return value error-nya
						_, _ = os.Stdout.WriteString("  ↳ " + pretty + "\n")
					}
				}
			}
		},
	})
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func status_color(status int) string {
	switch {
	case status >= 500:
		return "\x1b[31m" // Merah
	case status >= 400:
		return "\x1b[33m" // Kuning
	case status >= 300:
		return "\x1b[36m" // Cyan
	case status >= 200:
		return "\x1b[32m" // Hijau
	default:
		return "\x1b[0m"
	}
}

func http_status_text(status int) string {
	texts := map[int]string{
		200: "OK",
		201: "Created",
		204: "No Content",
		301: "Moved",
		302: "Found",
		400: "Bad Request",
		401: "Unauthorized",
		403: "Forbidden",
		404: "Not Found",
		405: "Method Not Allowed",
		409: "Conflict",
		422: "Unprocessable",
		429: "Too Many Req",
		500: "Server Error",
		502: "Bad Gateway",
		503: "Unavailable",
	}
	if text, ok := texts[status]; ok {
		return text
	}
	return "Unknown"
}

func pretty_json(body []byte) string {
	var parsed any
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "" // bukan JSON valid, skip
	}

	var buf bytes.Buffer
	colorize_json(&buf, parsed, "    ")
	return buf.String()
}

func colorize_json(buf *bytes.Buffer, v interface{}, indent string) {
	const (
		color_key    = "\x1b[96m" // Cyan terang — key
		color_string = "\x1b[32m" // Hijau — string value
		color_number = "\x1b[35m" // Magenta — number
		color_bool   = "\x1b[33m" // Kuning — bool
		color_null   = "\x1b[31m" // Merah — null
		color_reset  = "\x1b[0m"
	)

	child_indent := indent + "  "

	switch val := v.(type) {

	case map[string]any:
		buf.WriteString("{\n")
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		// Sort keys agar output konsisten
		sort.Strings(keys)
		for i, k := range keys {
			buf.WriteString(child_indent)
			buf.WriteString(color_key + `"` + k + `"` + color_reset + ": ")
			colorize_json(buf, val[k], child_indent)
			if i < len(keys)-1 {
				buf.WriteString(",")
			}
			buf.WriteString("\n")
		}
		buf.WriteString(indent + "}")

	case []any:
		buf.WriteString("[\n")
		for i, item := range val {
			buf.WriteString(child_indent)
			colorize_json(buf, item, child_indent)
			if i < len(val)-1 {
				buf.WriteString(",")
			}
			buf.WriteString("\n")
		}
		buf.WriteString(indent + "]")

	case string:
		buf.WriteString(color_string + `"` + val + `"` + color_reset)

	case float64:
		// JSON number selalu float64 saat di-unmarshal
		if val == float64(int64(val)) {
			buf.WriteString(color_number + fmt.Sprintf("%.0f", val) + color_reset)
		} else {
			buf.WriteString(color_number + fmt.Sprintf("%g", val) + color_reset)
		}

	case bool:
		buf.WriteString(color_bool + fmt.Sprintf("%t", val) + color_reset)

	case nil:
		buf.WriteString(color_null + "null" + color_reset)
	}
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
