package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/yourusername/go_server/internal/handlers/divisions"
	"github.com/yourusername/go_server/internal/middlewares"
	"github.com/yourusername/go_server/internal/services"
	"github.com/yourusername/go_server/internal/types"
)

type DivisionsIntegrationTestSuite struct {
	suite.Suite
	app       *fiber.App
	pool      *pgxpool.Pool
	token     string
	companyID string
}

func (s *DivisionsIntegrationTestSuite) SetupSuite() {
	s.app = fiber.New()
	s.pool = testDBPool

	ctx := context.Background()
	authSvc := services.NewAuthService(s.pool)
	compSvc := services.NewCompaniesService(s.pool)

	// Seed User
	email := "div-test@officecore.id"
	pass := "Pass123!"
	_, _ = authSvc.Register(ctx, &types.RegisterRequest{
		Email:    email,
		Username: "div_test",
		Password: pass,
	})
	loginRes, _ := authSvc.Login(ctx, &types.LoginRequest{
		Email:    email,
		Password: pass,
	})
	s.token = loginRes.AccessToken

	// Seed Company
	comp, _ := compSvc.Create(ctx, types.CreateCompanyRequest{
		Name: "Div Test Company",
		Code: "DIVTEST01",
	})
	s.companyID = comp.ID

	// Setup routes
	api := s.app.Group("/api/v1")
	api.Use(middlewares.DBMiddleware(s.pool))
	api.Use(middlewares.Pagination)

	div := api.Group("/divisions")
	div.Post("/", divisions.Create)
	div.Get("/", divisions.GetAll)
	div.Get("/:id", divisions.GetByID)
	div.Put("/:id", divisions.Update)
	div.Delete("/:id", divisions.Delete)
}

func (s *DivisionsIntegrationTestSuite) TestDivisionsCRUD() {
	var createdID string

	s.Run("Create_Success", func() {
		code := "ITD"
		payload := types.CreateDivisionRequest{
			CompanyID: s.companyID,
			Name:      "IT Division",
			Code:      &code,
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/v1/divisions", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, _ := s.app.Test(req)
		s.Equal(fiber.StatusCreated, resp.StatusCode)

		var res types.StandardResponse[types.DivisionRow]
		_ = json.NewDecoder(resp.Body).Decode(&res)
		s.Equal("IT Division", res.Data.Name)
		s.Equal("ITD", *res.Data.Code)
		createdID = res.Data.ID
	})

	s.Run("Create_DuplicateCode", func() {
		code := "ITD"
		payload := types.CreateDivisionRequest{
			CompanyID: s.companyID,
			Name:      "IT Division Duplicate",
			Code:      &code,
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/v1/divisions", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, _ := s.app.Test(req)
		s.Equal(fiber.StatusConflict, resp.StatusCode)
	})

	s.Run("Create_InvalidCompany", func() {
		payload := types.CreateDivisionRequest{
			CompanyID: uuid.NewString(),
			Name:      "Invalid Company Division",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/v1/divisions", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, _ := s.app.Test(req)
		s.Equal(fiber.StatusNotFound, resp.StatusCode)
	})

	s.Run("GetAll_Success", func() {
		req := httptest.NewRequest("GET", "/api/v1/divisions?company_id="+s.companyID, nil)
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, _ := s.app.Test(req)
		s.Equal(fiber.StatusOK, resp.StatusCode)

		var res types.StandardResponse[[]types.DivisionRow]
		_ = json.NewDecoder(resp.Body).Decode(&res)
		s.GreaterOrEqual(len(res.Data), 1)
	})

	s.Run("GetByID_Success", func() {
		req := httptest.NewRequest("GET", "/api/v1/divisions/"+createdID, nil)
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, _ := s.app.Test(req)
		s.Equal(fiber.StatusOK, resp.StatusCode)

		var res types.StandardResponse[types.DivisionRow]
		_ = json.NewDecoder(resp.Body).Decode(&res)
		s.Equal(createdID, res.Data.ID)
	})

	s.Run("Update_Success", func() {
		newName := "IT Division Updated"
		payload := types.UpdateDivisionRequest{
			Name: &newName,
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("PUT", "/api/v1/divisions/"+createdID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, _ := s.app.Test(req)
		s.Equal(fiber.StatusOK, resp.StatusCode)

		var res types.StandardResponse[types.DivisionRow]
		_ = json.NewDecoder(resp.Body).Decode(&res)
		s.Equal("IT Division Updated", res.Data.Name)
	})

	s.Run("Update_Code", func() {
		newCode := "ITD-MOD"
		payload := types.UpdateDivisionRequest{
			Code: &newCode,
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("PUT", "/api/v1/divisions/"+createdID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, _ := s.app.Test(req)
		s.Equal(fiber.StatusOK, resp.StatusCode)

		var res types.StandardResponse[types.DivisionRow]
		_ = json.NewDecoder(resp.Body).Decode(&res)
		s.Equal("ITD-MOD", *res.Data.Code)
	})

	s.Run("Update_DuplicateCode", func() {
		// Create another one first
		code := "OTHER"
		payloadC := types.CreateDivisionRequest{
			CompanyID: s.companyID,
			Name:      "Other Division",
			Code:      &code,
		}
		bodyC, _ := json.Marshal(payloadC)
		reqC := httptest.NewRequest("POST", "/api/v1/divisions", bytes.NewBuffer(bodyC))
		reqC.Header.Set("Content-Type", "application/json")
		reqC.Header.Set("Authorization", "Bearer "+s.token)
		respC, _ := s.app.Test(reqC)
		s.Equal(fiber.StatusCreated, respC.StatusCode)

		// Try to update ITD-MOD to OTHER
		otherCode := "OTHER"
		payload := types.UpdateDivisionRequest{
			Code: &otherCode,
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("PUT", "/api/v1/divisions/"+createdID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, _ := s.app.Test(req)
		s.Equal(fiber.StatusConflict, resp.StatusCode)
	})

	s.Run("Delete_Success", func() {
		req := httptest.NewRequest("DELETE", "/api/v1/divisions/"+createdID, nil)
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, _ := s.app.Test(req)
		s.Equal(fiber.StatusNoContent, resp.StatusCode)

		// Verify not found after delete
		reqGet := httptest.NewRequest("GET", "/api/v1/divisions/"+createdID, nil)
		reqGet.Header.Set("Authorization", "Bearer "+s.token)
		respGet, _ := s.app.Test(reqGet)
		s.Equal(fiber.StatusNotFound, respGet.StatusCode)
	})

	s.Run("GetByID_NotFound", func() {
		req := httptest.NewRequest("GET", "/api/v1/divisions/"+uuid.NewString(), nil)
		req.Header.Set("Authorization", "Bearer "+s.token)
		resp, _ := s.app.Test(req)
		s.Equal(fiber.StatusNotFound, resp.StatusCode)
	})

	s.Run("Update_NotFound", func() {
		payload := types.UpdateDivisionRequest{}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("PUT", "/api/v1/divisions/"+uuid.NewString(), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)
		resp, _ := s.app.Test(req)
		s.Equal(fiber.StatusNotFound, resp.StatusCode)
	})

	s.Run("Update_AfterDelete_NotFound", func() {
		payload := types.UpdateDivisionRequest{}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("PUT", "/api/v1/divisions/"+createdID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)
		resp, _ := s.app.Test(req)
		s.Equal(fiber.StatusNotFound, resp.StatusCode)
	})

	s.Run("Delete_NotFound", func() {
		req := httptest.NewRequest("DELETE", "/api/v1/divisions/"+uuid.NewString(), nil)
		req.Header.Set("Authorization", "Bearer "+s.token)
		resp, _ := s.app.Test(req)
		s.Equal(fiber.StatusNotFound, resp.StatusCode)
	})

	s.Run("GetAll_LargeOffset", func() {
		req := httptest.NewRequest("GET", "/api/v1/divisions?company_id="+s.companyID+"&page=100", nil)
		req.Header.Set("Authorization", "Bearer "+s.token)
		resp, _ := s.app.Test(req)
		s.Equal(fiber.StatusOK, resp.StatusCode)
	})
}

func TestDivisionsIntegration(t *testing.T) {
	suite.Run(t, new(DivisionsIntegrationTestSuite))
}
