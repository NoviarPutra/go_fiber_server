package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/yourusername/go_server/internal/handlers/companies"
	"github.com/yourusername/go_server/internal/middlewares"
	"github.com/yourusername/go_server/internal/services"
	"github.com/yourusername/go_server/internal/types"
)

type CompaniesIntegrationTestSuite struct {
	suite.Suite
	app   *fiber.App
	pool  *pgxpool.Pool
	token string
}

func (s *CompaniesIntegrationTestSuite) SetupSuite() {
	s.app = fiber.New()
	s.pool = testDBPool

	// Seed user for auth
	ctx := context.Background()
	authSvc := services.NewAuthService(s.pool)
	email := "comp-test@officecore.id"
	pass := "Pass123!"

	_, err := authSvc.Register(ctx, &types.RegisterRequest{
		Email:    email,
		Username: "comptest",
		Password: pass,
	})
	s.Require().NoError(err)

	loginRes, err := authSvc.Login(ctx, &types.LoginRequest{
		Email:    email,
		Password: pass,
	})
	s.Require().NoError(err)
	s.token = loginRes.AccessToken

	// Setup routes
	api := s.app.Group("/api/v1")
	api.Use(middlewares.DBMiddleware(s.pool))
	api.Use(middlewares.Pagination) // Tanpa tanda kurung

	comp := api.Group("/companies")
	comp.Post("/", companies.Create)
	comp.Get("/", companies.GetAll)
	comp.Get("/:id", companies.GetByID)
	comp.Put("/:id", companies.Update)
	comp.Delete("/:id", companies.Delete)
}

func (s *CompaniesIntegrationTestSuite) TestCompaniesCRUD() {
	var companyID string
	companyCode := "COMP001"

	s.Run("Create_Company_Success", func() {
		payload := types.CreateCompanyRequest{
			Name: "Test Company",
			Code: companyCode,
		}
		body, err := json.Marshal(payload)
		s.Require().NoError(err)
		req := httptest.NewRequest("POST", "/api/v1/companies", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusCreated, resp.StatusCode)

		var result types.StandardResponse[types.CompanyRow]
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(s.T(), err)
		companyID = result.Data.ID
		s.NotEmpty(companyID)
	})

	s.Run("Create_Company_Duplicate_Code", func() {
		payload := types.CreateCompanyRequest{
			Name: "Duplicate",
			Code: companyCode,
		}
		body, err := json.Marshal(payload)
		s.Require().NoError(err)
		req := httptest.NewRequest("POST", "/api/v1/companies", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusConflict, resp.StatusCode)
	})

	s.Run("Create_Company_Validation_Error", func() {
		payload := types.CreateCompanyRequest{
			Name: "", // Required but empty
			Code: "INVALID",
		}
		body, err := json.Marshal(payload)
		s.Require().NoError(err)
		req := httptest.NewRequest("POST", "/api/v1/companies", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusBadRequest, resp.StatusCode)
	})

	s.Run("GetAll_Companies_Pagination", func() {
		req := httptest.NewRequest("GET", "/api/v1/companies?page=1&limit=5", nil)
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusOK, resp.StatusCode)
	})

	s.Run("GetByID_Success", func() {
		req := httptest.NewRequest("GET", "/api/v1/companies/"+companyID, nil)
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusOK, resp.StatusCode)

		var result types.StandardResponse[types.CompanyRow]
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(s.T(), err)
		s.Equal(companyID, result.Data.ID)
	})

	s.Run("GetByID_NotFound", func() {
		req := httptest.NewRequest("GET", "/api/v1/companies/00000000-0000-0000-0000-000000000000", nil)
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusNotFound, resp.StatusCode)
	})

	s.Run("GetByID_Invalid_Format", func() {
		req := httptest.NewRequest("GET", "/api/v1/companies/not-a-uuid", nil)
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusInternalServerError, resp.StatusCode)
	})

	s.Run("Update_Company_Success", func() {
		newName := "Updated Company Name"
		payload := types.UpdateCompanyRequest{
			Name: ptr(newName),
		}
		body, err := json.Marshal(payload)
		s.Require().NoError(err)
		req := httptest.NewRequest("PUT", "/api/v1/companies/"+companyID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusOK, resp.StatusCode)

		var result types.StandardResponse[types.CompanyRow]
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(s.T(), err)
		s.Equal(newName, result.Data.Name)
	})

	s.Run("GetAll_Companies_Empty", func() {
		req := httptest.NewRequest("GET", "/api/v1/companies?page=100&limit=10", nil)
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusOK, resp.StatusCode)

		var result types.StandardResponse[[]types.CompanyRow]
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(s.T(), err)
		s.Equal(0, len(result.Data))
	})

	s.Run("Update_Company_NotFound", func() {
		payload := types.UpdateCompanyRequest{
			Name: ptr("Not Found"),
		}
		body, err := json.Marshal(payload)
		s.Require().NoError(err)
		req := httptest.NewRequest("PUT", "/api/v1/companies/00000000-0000-0000-0000-000000000000", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusNotFound, resp.StatusCode)
	})

	s.Run("Update_Company_Duplicate_Code", func() {
		// Create another company first
		reqC := httptest.NewRequest("POST", "/api/v1/companies", bytes.NewBufferString(`{"name":"Other","code":"OTHER001"}`))
		reqC.Header.Set("Content-Type", "application/json")
		reqC.Header.Set("Authorization", "Bearer "+s.token)
		respC, err := s.app.Test(reqC)
		require.NoError(s.T(), err)
		defer respC.Body.Close()

		payload := types.UpdateCompanyRequest{
			Code: ptr("OTHER001"),
		}
		body, err := json.Marshal(payload)
		s.Require().NoError(err)
		req := httptest.NewRequest("PUT", "/api/v1/companies/"+companyID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusBadRequest, resp.StatusCode)
	})

	s.Run("Update_Company_Partial", func() {
		payload := types.UpdateCompanyRequest{
			LogoUrl: ptr("https://new-logo.com"),
		}
		body, err := json.Marshal(payload)
		s.Require().NoError(err)
		req := httptest.NewRequest("PUT", "/api/v1/companies/"+companyID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusOK, resp.StatusCode)
	})

	s.Run("Delete_Company_NotFound", func() {
		req := httptest.NewRequest("DELETE", "/api/v1/companies/00000000-0000-0000-0000-000000000000", nil)
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusNotFound, resp.StatusCode)
	})

	s.Run("Update_Company_All_Nil", func() {
		payload := types.UpdateCompanyRequest{}
		body, err := json.Marshal(payload)
		s.Require().NoError(err)
		req := httptest.NewRequest("PUT", "/api/v1/companies/"+companyID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusOK, resp.StatusCode)
	})

	s.Run("Delete_Company_Success", func() {
		req := httptest.NewRequest("DELETE", "/api/v1/companies/"+companyID, nil)
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusNoContent, resp.StatusCode)
	})

	s.Run("GetByID_After_Delete_NotFound", func() {
		req := httptest.NewRequest("GET", "/api/v1/companies/"+companyID, nil)
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusNotFound, resp.StatusCode)
	})

	s.Run("Edge_Cases_Robustness", func() {
		reqC := httptest.NewRequest("POST", "/api/v1/companies", bytes.NewBufferString("{invalid}"))
		reqC.Header.Set("Content-Type", "application/json")
		reqC.Header.Set("Authorization", "Bearer "+s.token)
		respC, err := s.app.Test(reqC)
		require.NoError(s.T(), err)
		defer respC.Body.Close()
		s.Equal(400, respC.StatusCode)

		reqU := httptest.NewRequest("PUT", "/api/v1/companies/"+companyID, bytes.NewBufferString("{invalid}"))
		reqU.Header.Set("Content-Type", "application/json")
		reqU.Header.Set("Authorization", "Bearer "+s.token)
		respU, err := s.app.Test(reqU)
		require.NoError(s.T(), err)
		defer respU.Body.Close()
		s.Equal(400, respU.StatusCode)

		reqP := httptest.NewRequest("GET", "/api/v1/companies?page=-1&limit=-5", nil)
		reqP.Header.Set("Authorization", "Bearer "+s.token)
		respP, err := s.app.Test(reqP)
		require.NoError(s.T(), err)
		defer respP.Body.Close()
		s.Equal(200, respP.StatusCode)
	})
}

func TestCompaniesIntegration(t *testing.T) {
	suite.Run(t, new(CompaniesIntegrationTestSuite))
}

func ptr(s string) *string {
	return &s
}
