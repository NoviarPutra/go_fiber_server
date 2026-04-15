package types

import "time"

type UserDeviceRow struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	DeviceName *string    `json:"device_name"`
	DeviceType *string    `json:"device_type"`
	OS         *string    `json:"os"`
	LastActive time.Time  `json:"last_active"`
	PushToken  *string    `json:"push_token"`
	IsOnline   bool       `json:"is_online"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	RevokedAt  *time.Time `json:"revoked_at"`
}

type RegisterDeviceRequest struct {
	DeviceName string `json:"device_name" validate:"required"`
	DeviceType string `json:"device_type" validate:"required,oneof=ios android web"`
	OS         string `json:"os"          validate:"required"`
	PushToken  string `json:"push_token"`
}

type UpdatePushTokenRequest struct {
	PushToken string `json:"push_token" validate:"required"`
}
