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
	"github.com/yourusername/go_server/internal/handlers/company_branches"
	"github.com/yourusername/go_server/internal/middlewares"
	"github.com/yourusername/go_server/internal/services"
	"github.com/yourusername/go_server/internal/types"
)

type CompanyBranchesIntegrationTestSuite struct {
	suite.Suite
	app       *fiber.App
	pool      *pgxpool.Pool
	token     string
	companyID string // Seeded company for branches
}

func (s *CompanyBranchesIntegrationTestSuite) SetupSuite() {
	s.app = fiber.New()
	s.pool = testDBPool

	ctx := context.Background()
	authSvc := services.NewAuthService(s.pool)
	email := "branch-test@officecore.id"
	pass := "Pass123!"

	_, _ = authSvc.Register(ctx, &types.RegisterRequest{
		Email:    email,
		Username: "branchtest",
		Password: pass,
	})

	loginRes, err := authSvc.Login(ctx, &types.LoginRequest{
		Email:    email,
		Password: pass,
	})
	s.Require().NoError(err)
	s.token = loginRes.AccessToken

	api := s.app.Group("/api/v1")
	api.Use(middlewares.DBMiddleware(s.pool))
	api.Use(middlewares.Pagination)

	// Register companies route so we can seed one or we just seed it using service
	compSvc := services.NewCompaniesService(s.pool)
	company, err := compSvc.Create(ctx, types.CreateCompanyRequest{
		Name: "Branch Master Company",
		Code: "BMC001",
	})
	s.Require().NoError(err)
	s.companyID = company.ID

	// Register branches routes
	branches := api.Group("/company-branches")
	branches.Post("/", company_branches.Create)
	branches.Get("/", company_branches.GetAll)
	branches.Get("/:id", company_branches.GetByID)
	branches.Put("/:id", company_branches.Update)
	branches.Delete("/:id", company_branches.Delete)
}

func (s *CompanyBranchesIntegrationTestSuite) TestCompanyBranchesCRUD() {
	var branchID string

	s.Run("Create_Branch_Success", func() {
		payload := types.CreateCompanyBranchRequest{
			CompanyID: s.companyID,
			Name:      "Kantor Pusat JKT",
			Timezone:  "Asia/Jakarta",
		}
		body, err := json.Marshal(payload)
		s.Require().NoError(err)
		req := httptest.NewRequest("POST", "/api/v1/company-branches", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusCreated, resp.StatusCode)

		var result types.StandardResponse[types.CompanyBranchRow]
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(s.T(), err)
		branchID = result.Data.ID
		s.NotEmpty(branchID)
	})

	s.Run("Create_Branch_Duplicate_Name", func() {
		payload := types.CreateCompanyBranchRequest{
			CompanyID: s.companyID,
			Name:      "Kantor Pusat JKT", // Duplicate
			Timezone:  "Asia/Jakarta",
		}
		body, err := json.Marshal(payload)
		s.Require().NoError(err)
		req := httptest.NewRequest("POST", "/api/v1/company-branches", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusConflict, resp.StatusCode)
	})

	s.Run("Create_Branch_Validation_Error", func() {
		payload := types.CreateCompanyBranchRequest{
			CompanyID: "not-a-uuid", // Invalid UUID
			Name:      "",
		}
		body, err := json.Marshal(payload)
		s.Require().NoError(err)
		req := httptest.NewRequest("POST", "/api/v1/company-branches", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusBadRequest, resp.StatusCode)
	})

	s.Run("GetAll_Branches_Pagination", func() {
		req := httptest.NewRequest("GET", "/api/v1/company-branches?page=1&limit=5", nil)
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusOK, resp.StatusCode)
	})

	s.Run("GetByID_Success", func() {
		req := httptest.NewRequest("GET", "/api/v1/company-branches/"+branchID, nil)
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusOK, resp.StatusCode)

		var result types.StandardResponse[types.CompanyBranchRow]
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(s.T(), err)
		s.Equal(branchID, result.Data.ID)
	})

	s.Run("GetByID_NotFound", func() {
		req := httptest.NewRequest("GET", "/api/v1/company-branches/00000000-0000-0000-0000-000000000000", nil)
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusNotFound, resp.StatusCode)
	})

	s.Run("GetByID_Invalid_Format", func() {
		req := httptest.NewRequest("GET", "/api/v1/company-branches/not-a-uuid", nil)
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusInternalServerError, resp.StatusCode)
	})

	s.Run("Update_Branch_Success", func() {
		newName := "Updated Branch Name"
		payload := types.UpdateCompanyBranchRequest{
			Name: ptr(newName),
		}
		body, err := json.Marshal(payload)
		s.Require().NoError(err)
		req := httptest.NewRequest("PUT", "/api/v1/company-branches/"+branchID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusOK, resp.StatusCode)

		var result types.StandardResponse[types.CompanyBranchRow]
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(s.T(), err)
		s.Equal(newName, result.Data.Name)
	})

	s.Run("GetAll_Branches_Empty", func() {
		req := httptest.NewRequest("GET", "/api/v1/company-branches?page=100&limit=10", nil)
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusOK, resp.StatusCode)
	})

	s.Run("Update_Branch_NotFound", func() {
		payload := types.UpdateCompanyBranchRequest{
			Name: ptr("Not Found"),
		}
		body, err := json.Marshal(payload)
		s.Require().NoError(err)
		req := httptest.NewRequest("PUT", "/api/v1/company-branches/00000000-0000-0000-0000-000000000000", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusNotFound, resp.StatusCode)
	})

	s.Run("Update_Branch_Duplicate_Name", func() {
		payloadOther := types.CreateCompanyBranchRequest{
			CompanyID: s.companyID,
			Name:      "Other Branch",
			Timezone:  "Asia/Jakarta",
		}
		bodyOther, _ := json.Marshal(payloadOther)
		reqC := httptest.NewRequest("POST", "/api/v1/company-branches", bytes.NewBuffer(bodyOther))
		reqC.Header.Set("Content-Type", "application/json")
		reqC.Header.Set("Authorization", "Bearer "+s.token)
		respC, err := s.app.Test(reqC)
		require.NoError(s.T(), err)
		defer respC.Body.Close()

		payload := types.UpdateCompanyBranchRequest{
			Name: ptr("Other Branch"),
		}
		body, err := json.Marshal(payload)
		s.Require().NoError(err)
		req := httptest.NewRequest("PUT", "/api/v1/company-branches/"+branchID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusConflict, resp.StatusCode) // Name exists
	})

	s.Run("Update_Branch_Partial", func() {
		payload := types.UpdateCompanyBranchRequest{
			Timezone: ptr("Asia/Makassar"),
		}
		body, err := json.Marshal(payload)
		s.Require().NoError(err)
		req := httptest.NewRequest("PUT", "/api/v1/company-branches/"+branchID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusOK, resp.StatusCode)
	})

	s.Run("Delete_Branch_NotFound", func() {
		req := httptest.NewRequest("DELETE", "/api/v1/company-branches/00000000-0000-0000-0000-000000000000", nil)
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusNotFound, resp.StatusCode)
	})

	s.Run("Update_Branch_All_Nil", func() {
		payload := types.UpdateCompanyBranchRequest{}
		body, err := json.Marshal(payload)
		s.Require().NoError(err)
		req := httptest.NewRequest("PUT", "/api/v1/company-branches/"+branchID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusOK, resp.StatusCode)
	})

	s.Run("Delete_Branch_Success", func() {
		req := httptest.NewRequest("DELETE", "/api/v1/company-branches/"+branchID, nil)
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusNoContent, resp.StatusCode)
	})

	s.Run("GetByID_After_Delete_NotFound", func() {
		req := httptest.NewRequest("GET", "/api/v1/company-branches/"+branchID, nil)
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		require.NoError(s.T(), err)
		defer resp.Body.Close()

		s.Equal(fiber.StatusNotFound, resp.StatusCode)
	})

	s.Run("Edge_Cases_Robustness", func() {
		reqC := httptest.NewRequest("POST", "/api/v1/company-branches", bytes.NewBufferString("{invalid}"))
		reqC.Header.Set("Content-Type", "application/json")
		reqC.Header.Set("Authorization", "Bearer "+s.token)
		respC, err := s.app.Test(reqC)
		require.NoError(s.T(), err)
		defer respC.Body.Close()
		s.Equal(400, respC.StatusCode)

		reqU := httptest.NewRequest("PUT", "/api/v1/company-branches/"+branchID, bytes.NewBufferString("{invalid}"))
		reqU.Header.Set("Content-Type", "application/json")
		reqU.Header.Set("Authorization", "Bearer "+s.token)
		respU, err := s.app.Test(reqU)
		require.NoError(s.T(), err)
		defer respU.Body.Close()
		s.Equal(400, respU.StatusCode)

		reqP := httptest.NewRequest("GET", "/api/v1/company-branches?page=-1&limit=-5", nil)
		reqP.Header.Set("Authorization", "Bearer "+s.token)
		respP, err := s.app.Test(reqP)
		require.NoError(s.T(), err)
		defer respP.Body.Close()
		s.Equal(200, respP.StatusCode)
	})
}

func TestCompanyBranchesIntegration(t *testing.T) {
	suite.Run(t, new(CompanyBranchesIntegrationTestSuite))
}
