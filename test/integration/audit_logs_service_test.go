package integration

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/yourusername/go_server/internal/handlers/audit_logs"
	"github.com/yourusername/go_server/internal/middlewares"
	"github.com/yourusername/go_server/internal/services"
	"github.com/yourusername/go_server/internal/types"
	"github.com/yourusername/go_server/internal/utils"
)

type AuditLogsIntegrationTestSuite struct {
	suite.Suite
	app       *fiber.App
	pool      *pgxpool.Pool
	token     string
	companyID string
	userID    string
}

func (s *AuditLogsIntegrationTestSuite) SetupSuite() {
	s.app = fiber.New()
	s.pool = testDBPool

	ctx := context.Background()
	authSvc := services.NewAuthService(s.pool)
	email := "audit-test@officecore.id"
	pass := "Pass123!"

	_, _ = authSvc.Register(ctx, &types.RegisterRequest{
		Email:    email,
		Username: "audittest",
		Password: pass,
	})

	loginRes, err := authSvc.Login(ctx, &types.LoginRequest{
		Email:    email,
		Password: pass,
	})
	s.Require().NoError(err)
	s.token = loginRes.AccessToken

	// Find the user ID dynamically
	var uid string
	err = s.pool.QueryRow(ctx, "SELECT id::text FROM users WHERE email = $1", email).Scan(&uid)
	s.Require().NoError(err)
	s.userID = uid

	compSvc := services.NewCompaniesService(s.pool)
	code := "AUDIT_CO"
	_ = compSvc.Delete(ctx, code) // Clean up incase
	
	comp, err := compSvc.Create(ctx, types.CreateCompanyRequest{
		Name:    "Audit Company",
		Code:    "AUDITCO",
	})
	s.Require().NoError(err)
	s.companyID = comp.ID

	api := s.app.Group("/api/v1")
	api.Use(middlewares.DBMiddleware(s.pool))
	
	audit_group := api.Group("/audit-logs")
	audit_group.Get("/", middlewares.Protected(), middlewares.Pagination, audit_logs.GetAll)
	audit_group.Get("/:id", middlewares.Protected(), audit_logs.GetByID)
}

func (s *AuditLogsIntegrationTestSuite) Test1_UpdateAndLogAction() {
	compSvc := services.NewCompaniesService(s.pool)

	info := utils.AuditInfo{
		UserID:    s.userID,
		IPAddress: "127.0.0.1",
		UserAgent: "Go-HTTP-Client",
	}
	ctx := utils.ContextWithAuditInfo(context.Background(), info)

	// Ubah nama company untuk memicu trigger UPDATE
	newName := "Updated Audit Company"
	_, err := compSvc.Update(ctx, s.companyID, types.UpdateCompanyRequest{
		Name: &newName,
	})

	s.Require().NoError(err)
}

func (s *AuditLogsIntegrationTestSuite) Test2_GetAllLogs() {
	req := httptest.NewRequest("GET", "/api/v1/audit-logs?table_name=companies", nil)
	req.Header.Set("Authorization", "Bearer "+s.token)

	resp, err := s.app.Test(req, -1)
	s.Require().NoError(err)
	s.Equal(200, resp.StatusCode)

	var resBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&resBody)
	s.Require().NoError(err)

	s.True(resBody["success"].(bool))
	data := resBody["data"].([]interface{})
	s.GreaterOrEqual(len(data), 1)

	// Verifikasi bahwa IP address dan Action terekam
	firstLog := data[0].(map[string]interface{})
	s.Equal("UPDATE", firstLog["action"])
	s.Equal("companies", firstLog["table_name"])
	s.Equal("127.0.0.1", firstLog["ip_address"])

	// Test more filters for coverage
	logRecID := firstLog["record_id"].(string)
	
	// Filter by RecordID
	req2 := httptest.NewRequest("GET", "/api/v1/audit-logs?record_id="+logRecID, nil)
	req2.Header.Set("Authorization", "Bearer "+s.token)
	resp2, _ := s.app.Test(req2, -1)
	s.Equal(200, resp2.StatusCode)
	resp2.Body.Close()
	
	// Filter by UserID
	req3 := httptest.NewRequest("GET", "/api/v1/audit-logs?user_id="+s.userID, nil)
	req3.Header.Set("Authorization", "Bearer "+s.token)
	resp3, _ := s.app.Test(req3, -1)
	s.Equal(200, resp3.StatusCode)
	resp3.Body.Close()

	// Filter by CompanyID
	req4 := httptest.NewRequest("GET", "/api/v1/audit-logs?company_id="+s.companyID, nil)
	req4.Header.Set("Authorization", "Bearer "+s.token)
	resp4, _ := s.app.Test(req4, -1)
	s.Equal(200, resp4.StatusCode)
	resp4.Body.Close()
}

func (s *AuditLogsIntegrationTestSuite) Test3_GetLogByID() {
	// First fetch the log from prev test to get ID
	req := httptest.NewRequest("GET", "/api/v1/audit-logs?table_name=companies", nil)
	req.Header.Set("Authorization", "Bearer "+s.token)
	resp, _ := s.app.Test(req, -1)

	var resBody map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&resBody)
	logs := resBody["data"].([]interface{})
	logID := logs[0].(map[string]interface{})["id"].(string)

	// Fetch specific ID
	req2 := httptest.NewRequest("GET", "/api/v1/audit-logs/"+logID, nil)
	req2.Header.Set("Authorization", "Bearer "+s.token)

	resp2, err := s.app.Test(req2, -1)
	s.Require().NoError(err)
	s.Equal(200, resp2.StatusCode)

	var detailBody map[string]interface{}
	_ = json.NewDecoder(resp2.Body).Decode(&detailBody)
	s.True(detailBody["success"].(bool))

	logData := detailBody["data"].(map[string]interface{})
	s.Equal(logID, logData["id"])
	s.NotNil(logData["new_data"])
}

func (s *AuditLogsIntegrationTestSuite) Test4_GetLogByID_NotFound() {
	req := httptest.NewRequest("GET", "/api/v1/audit-logs/00000000-0000-0000-0000-000000000000", nil)
	req.Header.Set("Authorization", "Bearer "+s.token)

	resp, err := s.app.Test(req, -1)
	s.Require().NoError(err)
	s.Equal(fiber.StatusNotFound, resp.StatusCode)
}

func (s *AuditLogsIntegrationTestSuite) Test5_GetLogByID_Invalid_Format() {
	req := httptest.NewRequest("GET", "/api/v1/audit-logs/invalid-uuid", nil)
	req.Header.Set("Authorization", "Bearer "+s.token)

	resp, err := s.app.Test(req, -1)
	s.Require().NoError(err)
	// Bisa internal server error atau bad request, kita pastikan != 200
	s.NotEqual(200, resp.StatusCode)
}

func (s *AuditLogsIntegrationTestSuite) Test6_GetAllLogs_Empty() {
	req := httptest.NewRequest("GET", "/api/v1/audit-logs?table_name=unknown_table", nil)
	req.Header.Set("Authorization", "Bearer "+s.token)

	resp, err := s.app.Test(req, -1)
	s.Require().NoError(err)
	s.Equal(200, resp.StatusCode)

	var resBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&resBody)
	s.Require().NoError(err)

	if dataRaw, ok := resBody["data"]; ok && dataRaw != nil {
		if data, ok := dataRaw.([]interface{}); ok {
			s.Equal(0, len(data))
		}
	}
}

func TestAuditLogsIntegration(t *testing.T) {
	suite.Run(t, new(AuditLogsIntegrationTestSuite))
}
