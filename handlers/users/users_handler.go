package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/go_server/types"
	"github.com/yourusername/go_server/utils"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var dummyUsers = []User{
	{ID: 1, Name: "Andi"},
	{ID: 2, Name: "Budi"},
	{ID: 3, Name: "Citra"},
	{ID: 4, Name: "Dewi"},
	{ID: 5, Name: "Eko"},
}

func UsersHandler(ctx *fiber.Ctx) error {
	// 1. Safe Type Assertion (Mencegah Panic)
	page, ok := ctx.Locals("page").(int)
	if !ok || page < 1 {
		page = 1
	}

	perPage, ok := ctx.Locals("per_page").(int)
	if !ok || perPage < 1 {
		perPage = 10
	}

	// 2. Logic Boundary Safety
	total := len(dummyUsers)
	start := (page - 1) * perPage

	// Jika start lebih besar dari data, kembalikan list kosong (bukan panic)
	if start >= total {
		return utils.SuccessWithMeta(ctx, []User{}, "No more data", &types.Meta{
			Page: page, PerPage: perPage, Total: int64(total),
		})
	}

	end := min(start+perPage, total)

	users := dummyUsers[start:end]

	// 3. Dynamic Request ID (Gunakan Request ID dari context/header)
	// reqID := ctx.Locals("requestid").(string)
	reqID := utils.GetRequestID(ctx)

	meta := &types.Meta{
		Page:      page,
		PerPage:   perPage,
		Total:     int64(total),
		RequestID: reqID,
	}

	return utils.SuccessWithMeta(ctx, users, "Users retrieved successfully", meta)
}
