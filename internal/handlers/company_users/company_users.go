package company_users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/services"
	"github.com/yourusername/go_server/internal/types"
	"github.com/yourusername/go_server/internal/utils"
)

func Add(c *fiber.Ctx) error {
	companyIDStr := c.Params("company_id")
	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		return utils.BadRequest(c, "ID perusahaan tidak valid")
	}

	var req types.CreateCompanyUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Format data tidak valid")
	}

	db := c.Locals("db").(*pgxpool.Pool)
	svc := services.NewCompanyUsersService(db)

	res, err := svc.AddUser(c.Context(), companyID, req)
	if err != nil {
		if err == services.ErrCompanyUserAlreadyExists {
			return utils.Conflict(c, err.Error())
		}
		return utils.InternalError(c, err.Error())
	}

	return utils.Created(c, res, "User berhasil ditambahkan ke perusahaan")
}

func List(c *fiber.Ctx) error {
	companyIDStr := c.Params("company_id")
	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		return utils.BadRequest(c, "ID perusahaan tidak valid")
	}

	page, ok := c.Locals("page").(int)
	if !ok || page < 1 {
		page = 1
	}
	perPage, ok := c.Locals("per_page").(int)
	if !ok || perPage < 1 {
		perPage = 10
	}

	db := c.Locals("db").(*pgxpool.Pool)
	svc := services.NewCompanyUsersService(db)

	users, total, err := svc.List(c.Context(), companyID, page, perPage)
	if err != nil {
		return utils.InternalError(c, err.Error())
	}

	return utils.SuccessWithMeta(c, users, "Data user perusahaan berhasil diambil", &types.Meta{
		Page:    page,
		PerPage: perPage,
		Total:   total,
	})
}

func Get(c *fiber.Ctx) error {
	companyIDStr := c.Params("company_id")
	userIDStr := c.Params("user_id")

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		return utils.BadRequest(c, "ID perusahaan tidak valid")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.BadRequest(c, "ID user tidak valid")
	}

	db := c.Locals("db").(*pgxpool.Pool)
	svc := services.NewCompanyUsersService(db)

	res, err := svc.GetDetail(c.Context(), companyID, userID)
	if err != nil {
		if err == services.ErrCompanyUserNotFound {
			return utils.NotFound(c, err.Error())
		}
		return utils.InternalError(c, err.Error())
	}

	return utils.Success(c, res, "Data user perusahaan ditemukan")
}

func Update(c *fiber.Ctx) error {
	companyIDStr := c.Params("company_id")
	userIDStr := c.Params("user_id")

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		return utils.BadRequest(c, "ID perusahaan tidak valid")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.BadRequest(c, "ID user tidak valid")
	}

	var req types.UpdateCompanyUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Format data tidak valid")
	}

	db := c.Locals("db").(*pgxpool.Pool)
	svc := services.NewCompanyUsersService(db)

	res, err := svc.Update(c.Context(), companyID, userID, req)
	if err != nil {
		if err == services.ErrCompanyUserNotFound {
			return utils.NotFound(c, err.Error())
		}
		return utils.InternalError(c, err.Error())
	}

	return utils.Success(c, res, "Data user perusahaan berhasil diupdate")
}

func Remove(c *fiber.Ctx) error {
	companyIDStr := c.Params("company_id")
	userIDStr := c.Params("user_id")

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		return utils.BadRequest(c, "ID perusahaan tidak valid")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return utils.BadRequest(c, "ID user tidak valid")
	}

	db := c.Locals("db").(*pgxpool.Pool)
	svc := services.NewCompanyUsersService(db)

	err = svc.Remove(c.Context(), companyID, userID)
	if err != nil {
		if err == services.ErrCompanyUserNotFound {
			return utils.NotFound(c, err.Error())
		}
		return utils.InternalError(c, err.Error())
	}

	return utils.NoContent(c)
}
