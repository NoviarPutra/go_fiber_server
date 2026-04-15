package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"github.com/yourusername/go_server/internal/handlers/user_devices"
	"github.com/yourusername/go_server/internal/middlewares"
	"github.com/yourusername/go_server/internal/services"
	"github.com/yourusername/go_server/internal/types"
)

type UserDevicesIntegrationTestSuite struct {
	suite.Suite
	app   *fiber.App
	pool  *pgxpool.Pool
	token string
}

func (s *UserDevicesIntegrationTestSuite) SetupSuite() {
	s.app = fiber.New()
	s.pool = testDBPool

	// Seed user for auth
	ctx := context.Background()
	authSvc := services.NewAuthService(s.pool)
	email := "device-test@example.com"
	pass := "Pass123!"

	_, err := authSvc.Register(ctx, &types.RegisterRequest{
		Email:    email,
		Username: "devicetest",
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

	devices := api.Group("/user-devices")
	devices.Use(middlewares.Protected())
	devices.Get("/", user_devices.List)
	devices.Post("/", user_devices.Register)
	devices.Post("/:id/revoke", user_devices.Revoke)
	devices.Patch("/:id/push-token", user_devices.UpdatePushToken)
}

func (s *UserDevicesIntegrationTestSuite) TestUserDevicesFlow() {
	var deviceID string
	pushToken := "token_123"

	s.Run("Register_Device_Success", func() {
		payload := types.RegisterDeviceRequest{
			DeviceName: "iPhone 15 Pro",
			DeviceType: "ios",
			OS:         "iOS 17.4",
			PushToken:  pushToken,
		}
		body, err := json.Marshal(payload)
		s.Require().NoError(err)
		req := httptest.NewRequest("POST", "/api/v1/user-devices", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		s.Require().Equal(fiber.StatusOK, resp.StatusCode)

		var result types.StandardResponse[types.UserDeviceRow]
		err = json.NewDecoder(resp.Body).Decode(&result)
		s.Require().NoError(err)
		deviceID = result.Data.ID
		s.NotEmpty(deviceID)
		s.Require().NotNil(result.Data.PushToken)
		s.Equal(*result.Data.PushToken, pushToken)
	})

	s.Run("Register_Device_Duplicate_PushToken_Updates_Existing", func() {
		payload := types.RegisterDeviceRequest{
			DeviceName: "iPhone 15 Pro Updated",
			DeviceType: "ios",
			OS:         "iOS 17.5",
			PushToken:  pushToken,
		}
		body, err := json.Marshal(payload)
		s.Require().NoError(err)
		req := httptest.NewRequest("POST", "/api/v1/user-devices", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		s.Require().Equal(fiber.StatusOK, resp.StatusCode)

		var result types.StandardResponse[types.UserDeviceRow]
		err = json.NewDecoder(resp.Body).Decode(&result)
		s.Require().NoError(err)
		s.Equal(deviceID, result.Data.ID)
		s.Require().NotNil(result.Data.DeviceName)
		s.Equal(*result.Data.DeviceName, "iPhone 15 Pro Updated")
	})

	s.Run("Register_Device_Validation_Error", func() {
		payload := types.RegisterDeviceRequest{
			DeviceName: "", // Required
			DeviceType: "invalid_type", // oneof=ios android web
		}
		body, err := json.Marshal(payload)
		s.Require().NoError(err)
		req := httptest.NewRequest("POST", "/api/v1/user-devices", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		s.Equal(fiber.StatusBadRequest, resp.StatusCode)
	})

	s.Run("Register_Device_Without_PushToken_Success", func() {
		payload := types.RegisterDeviceRequest{
			DeviceName: "Web Browser",
			DeviceType: "web",
			OS:         "macOS",
			PushToken:  "", // Empty push token
		}
		body, err := json.Marshal(payload)
		s.Require().NoError(err)
		req := httptest.NewRequest("POST", "/api/v1/user-devices", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		s.Require().Equal(fiber.StatusOK, resp.StatusCode)

		var result types.StandardResponse[types.UserDeviceRow]
		err = json.NewDecoder(resp.Body).Decode(&result)
		s.Require().NoError(err)
		s.Nil(result.Data.PushToken)
	})

	s.Run("List_Devices_Success", func() {
		req := httptest.NewRequest("GET", "/api/v1/user-devices", nil)
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		s.Require().Equal(fiber.StatusOK, resp.StatusCode)

		var result types.StandardResponse[[]types.UserDeviceRow]
		err = json.NewDecoder(resp.Body).Decode(&result)
		s.Require().NoError(err)
		s.GreaterOrEqual(len(result.Data), 1)
	})

	s.Run("Update_PushToken_Success", func() {
		newToken := "new_token_456"
		payload := types.UpdatePushTokenRequest{
			PushToken: newToken,
		}
		body, err := json.Marshal(payload)
		s.Require().NoError(err)
		req := httptest.NewRequest("PATCH", "/api/v1/user-devices/"+deviceID+"/push-token", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		s.Require().Equal(fiber.StatusOK, resp.StatusCode)
	})

	s.Run("Update_PushToken_NotFound", func() {
		payload := types.UpdatePushTokenRequest{
			PushToken: "any_token",
		}
		body, err := json.Marshal(payload)
		s.Require().NoError(err)
		req := httptest.NewRequest("PATCH", "/api/v1/user-devices/00000000-0000-0000-0000-000000000000/push-token", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		s.Require().Equal(fiber.StatusNotFound, resp.StatusCode)
	})

	s.Run("Revoke_Device_Success", func() {
		req := httptest.NewRequest("POST", "/api/v1/user-devices/"+deviceID+"/revoke", nil)
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		s.Require().Equal(fiber.StatusOK, resp.StatusCode)

		// Verify it's gone from active list
		reqList := httptest.NewRequest("GET", "/api/v1/user-devices", nil)
		reqList.Header.Set("Authorization", "Bearer "+s.token)
		respList, err := s.app.Test(reqList)
		s.Require().NoError(err)
		var resultList types.StandardResponse[[]types.UserDeviceRow]
		err = json.NewDecoder(respList.Body).Decode(&resultList)
		s.Require().NoError(err)
		
		found := false
		for _, d := range resultList.Data {
			if d.ID == deviceID {
				found = true
				break
			}
		}
		s.False(found)
	})

	s.Run("Revoke_Device_NotFound_Or_Already_Revoked", func() {
		req := httptest.NewRequest("POST", "/api/v1/user-devices/"+deviceID+"/revoke", nil)
		req.Header.Set("Authorization", "Bearer "+s.token)

		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		s.Equal(fiber.StatusNotFound, resp.StatusCode)
	})

	s.Run("Unauthorized_Access", func() {
		req := httptest.NewRequest("GET", "/api/v1/user-devices", nil)
		// No token

		resp, err := s.app.Test(req)
		s.Require().NoError(err)
		s.Equal(fiber.StatusUnauthorized, resp.StatusCode)
	})
}

func TestUserDevicesIntegration(t *testing.T) {
	suite.Run(t, new(UserDevicesIntegrationTestSuite))
}
