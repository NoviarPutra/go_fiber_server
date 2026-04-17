package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/go_server/internal/types"
)

type UserDevicesService struct {
	db *pgxpool.Pool
}

func NewUserDevicesService(db *pgxpool.Pool) *UserDevicesService {
	return &UserDevicesService{db: db}
}

// RegisterDevice mendaftarkan perangkat baru atau memperbarui perangkat yang sudah ada
func (s *UserDevicesService) RegisterDevice(ctx context.Context, userID string, req *types.RegisterDeviceRequest) (*types.UserDeviceRow, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var device types.UserDeviceRow

	// Upsert berdasarkan user_id dan device_name jika belum di-revoke
	// Atau jika push_token diberikan, kita bisa pakai itu sebagai identifier unik jika perlu
	// Tapi migration hanya kasih unique index di push_token WHERE revoked_at IS NULL

	query := `
		INSERT INTO user_devices (user_id, device_name, device_type, os, push_token, is_online, last_active)
		VALUES ($1, $2, $3, $4, $5, true, now())
		ON CONFLICT (push_token) WHERE revoked_at IS NULL AND push_token IS NOT NULL 
		DO UPDATE SET 
			device_name = EXCLUDED.device_name,
			device_type = EXCLUDED.device_type,
			os = EXCLUDED.os,
			is_online = true,
			last_active = now(),
			updated_at = now()
		RETURNING id::text, user_id::text, device_name, device_type, os, last_active, push_token, is_online, created_at, updated_at, revoked_at
	`

	// Jika push_token kosong, kita tidak bisa pakai ON CONFLICT push_token
	if req.PushToken == "" {
		query = `
			INSERT INTO user_devices (user_id, device_name, device_type, os, is_online, last_active)
			VALUES ($1, $2, $3, $4, true, now())
			RETURNING id::text, user_id::text, device_name, device_type, os, last_active, push_token, is_online, created_at, updated_at, revoked_at
		`
		err := s.db.QueryRow(ctx, query, userID, req.DeviceName, req.DeviceType, req.OS).Scan(
			&device.ID, &device.UserID, &device.DeviceName, &device.DeviceType, &device.OS,
			&device.LastActive, &device.PushToken, &device.IsOnline, &device.CreatedAt, &device.UpdatedAt, &device.RevokedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal mendaftarkan perangkat: %w", err)
		}
	} else {
		err := s.db.QueryRow(ctx, query, userID, req.DeviceName, req.DeviceType, req.OS, req.PushToken).Scan(
			&device.ID, &device.UserID, &device.DeviceName, &device.DeviceType, &device.OS,
			&device.LastActive, &device.PushToken, &device.IsOnline, &device.CreatedAt, &device.UpdatedAt, &device.RevokedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal mendaftarkan perangkat: %w", err)
		}
	}

	return &device, nil
}

// ListDevices mengambil daftar perangkat aktif milik user
func (s *UserDevicesService) ListDevices(ctx context.Context, userID string) ([]types.UserDeviceRow, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
		SELECT id::text, user_id::text, device_name, device_type, os, last_active, push_token, is_online, created_at, updated_at, revoked_at
		FROM user_devices
		WHERE user_id = $1 AND revoked_at IS NULL
		ORDER BY last_active DESC
	`

	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil daftar perangkat: %w", err)
	}
	defer rows.Close()

	devices := make([]types.UserDeviceRow, 0)
	for rows.Next() {
		var d types.UserDeviceRow
		err := rows.Scan(
			&d.ID, &d.UserID, &d.DeviceName, &d.DeviceType, &d.OS,
			&d.LastActive, &d.PushToken, &d.IsOnline, &d.CreatedAt, &d.UpdatedAt, &d.RevokedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("gagal memproses data perangkat: %w", err)
		}
		devices = append(devices, d)
	}

	return devices, nil
}

// RevokeDevice menonaktifkan perangkat (log out dari perangkat tersebut)
func (s *UserDevicesService) RevokeDevice(ctx context.Context, userID, deviceID string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
		UPDATE user_devices
		SET revoked_at = now(), is_online = false, updated_at = now()
		WHERE id = $1 AND user_id = $2 AND revoked_at IS NULL
	`

	result, err := s.db.Exec(ctx, query, deviceID, userID)
	if err != nil {
		return fmt.Errorf("gagal menonaktifkan perangkat: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrDeviceNotFound
	}

	return nil
}

// UpdatePushToken memperbarui token push notification
func (s *UserDevicesService) UpdatePushToken(ctx context.Context, userID, deviceID string, req *types.UpdatePushTokenRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `
		UPDATE user_devices
		SET push_token = $1, updated_at = now()
		WHERE id = $2 AND user_id = $3 AND revoked_at IS NULL
	`

	result, err := s.db.Exec(ctx, query, req.PushToken, deviceID, userID)
	if err != nil {
		return fmt.Errorf("gagal memperbarui push token: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrDeviceNotFound
	}

	return nil
}
