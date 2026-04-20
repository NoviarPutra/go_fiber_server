package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/yourusername/go_server/internal/handlers/company_users"
	"github.com/yourusername/go_server/internal/middlewares"
	"github.com/yourusername/go_server/internal/services"
	"github.com/yourusername/go_server/internal/types"
)

type CompanyUsersIntegrationTestSuite struct {
	suite.Suite
	app        *fiber.App
	pool       *pgxpool.Pool
	token      string
	companyID  uuid.UUID
	testUserID uuid.UUID // User to be added
}

func (s *CompanyUsersIntegrationTestSuite) SetupSuite() {
	s.app = fiber.New()
	s.pool = testDBPool

	ctx := context.Background()
	authSvc := services.NewAuthService(s.pool)
	compSvc := services.NewCompaniesService(s.pool)

	// 1. Seed Manager User (for token)
	email := "cu-manager@officecore.id"
	pass := "Pass123!"
	_, _ = authSvc.Register(ctx, &types.RegisterRequest{
		Email:    email,
		Username: "cu_manager",
		Password: pass,
	})
	loginRes, _ := authSvc.Login(ctx, &types.LoginRequest{
		Email:    email,
		Password: pass,
	})
	s.token = loginRes.AccessToken

	// 2. Seed Company
	comp, _ := compSvc.Create(ctx, types.CreateCompanyRequest{
		Name: "CU Test Company",
		Code: "CUTest01",
	})
	compID, _ := uuid.Parse(comp.ID)
	s.companyID = compID

	// 3. Seed Target User (to be added to company)
	targetEmail := "cu-target@officecore.id"
	targetRes, _ := authSvc.Register(ctx, &types.RegisterRequest{
		Email:    targetEmail,
		Username: "cu_target",
		Password: pass,
	})
	targetID, _ := uuid.Parse(targetRes.ID)
	s.testUserID = targetID

	// Setup routes
	api := s.app.Group("/api/v1")
	api.Use(middlewares.DBMiddleware(s.pool))
	api.Use(middlewares.Pagination)

	cu := api.Group("/companies/:company_id/users")
	cu.Post("/", company_users.Add)
	cu.Get("/", company_users.List)
	cu.Get("/:user_id", company_users.Get)
	cu.Put("/:user_id", company_users.Update)
	cu.Delete("/:user_id", company_users.Remove)
}

func (s *CompanyUsersIntegrationTestSuite) TestCompanyUsersCRUD() {
	s.Run("AddUser_Success", func() {
		payload := types.CreateCompanyUserRequest{
			UserID:       s.testUserID,
			EmployeeCode: ptrStr("EMP-001"),
			IsOwner:      false,
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/companies/%s/users", s.companyID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, _ := s.app.Test(req)
		s.Equal(fiber.StatusCreated, resp.StatusCode)

		var result types.StandardResponse[types.CompanyUser]
		_ = json.NewDecoder(resp.Body).Decode(&result)
		s.Equal(s.testUserID, result.Data.UserID)
		s.Equal("EMP-001", *result.Data.EmployeeCode)
	})

	s.Run("AddUser_Duplicate", func() {
		payload := types.CreateCompanyUserRequest{
			UserID: s.testUserID,
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/companies/%s/users", s.companyID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, _ := s.app.Test(req)
		s.Equal(fiber.StatusConflict, resp.StatusCode)
	})

	s.Run("List_CompanyUsers", func() {
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/companies/%s/users?page=1&limit=10", s.companyID), nil)
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, _ := s.app.Test(req)
		s.Equal(fiber.StatusOK, resp.StatusCode)

		var result types.StandardResponse[[]types.CompanyUser]
		_ = json.NewDecoder(resp.Body).Decode(&result)
		s.GreaterOrEqual(len(result.Data), 1)
	})

	s.Run("GetDetail_Success", func() {
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/companies/%s/users/%s", s.companyID, s.testUserID), nil)
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, _ := s.app.Test(req)
		s.Equal(fiber.StatusOK, resp.StatusCode)

		var result types.StandardResponse[types.CompanyUser]
		_ = json.NewDecoder(resp.Body).Decode(&result)
		s.Equal(s.testUserID, result.Data.UserID)
		s.Equal("cu-target@officecore.id", *result.Data.UserEmail)
	})

	s.Run("Update_CompanyUser", func() {
		payload := types.UpdateCompanyUserRequest{
			EmployeeCode: ptrStr("EMP-001-MOD"),
			IsActive:     ptrBool(false),
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/companies/%s/users/%s", s.companyID, s.testUserID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, _ := s.app.Test(req)
		s.Equal(fiber.StatusOK, resp.StatusCode)

		var result types.StandardResponse[types.CompanyUser]
		_ = json.NewDecoder(resp.Body).Decode(&result)
		s.Equal("EMP-001-MOD", *result.Data.EmployeeCode)
		s.Equal(false, result.Data.IsActive)
	})

	s.Run("Remove_CompanyUser", func() {
		req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/companies/%s/users/%s", s.companyID, s.testUserID), nil)
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, _ := s.app.Test(req)
		s.Equal(fiber.StatusNoContent, resp.StatusCode)
	})

	s.Run("GetDetail_After_Remove_NotFound", func() {
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/companies/%s/users/%s", s.companyID, s.testUserID), nil)
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, _ := s.app.Test(req)
		s.Equal(fiber.StatusNotFound, resp.StatusCode)
	})
}

func TestCompanyUsersIntegration(t *testing.T) {
	suite.Run(t, new(CompanyUsersIntegrationTestSuite))
}

func ptrStr(s string) *string {
	return &s
}

func ptrBool(b bool) *bool {
	return &b
}
