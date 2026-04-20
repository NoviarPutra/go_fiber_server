package services

import "errors"

var (
	ErrEmailAlreadyExists    = errors.New("email sudah terdaftar")
	ErrUsernameAlreadyExists = errors.New("username sudah terdaftar")
	ErrRefreshTokenInvalid   = errors.New("refresh token tidak valid atau sudah kedaluwarsa")
	ErrRefreshTokenRevoked   = errors.New("refresh token telah dicabut")
	ErrUserNotFound          = errors.New("user tidak ditemukan")
	ErrInvalidCredentials    = errors.New("email atau password salah")
	ErrAccountInactive       = errors.New("akun tidak aktif")
	ErrRefreshTokenExpired   = errors.New("refresh token kadaluarsa")

	ErrCompanyNotFound   = errors.New("perusahaan tidak ditemukan")
	ErrCompanyCodeExists = errors.New("kode perusahaan sudah digunakan")

	ErrCompanyBranchNotFound   = errors.New("cabang perusahaan tidak ditemukan")
	ErrCompanyBranchNameExists = errors.New("nama cabang sudah digunakan di perusahaan ini")

	ErrAuditLogNotFound = errors.New("audit log tidak ditemukan")

	ErrDeviceNotFound       = errors.New("perangkat tidak ditemukan")
	ErrDeviceAlreadyRevoked = errors.New("perangkat sudah tidak aktif")

	ErrCompanyUserNotFound      = errors.New("user perusahaan tidak ditemukan")
	ErrCompanyUserAlreadyExists = errors.New("user sudah terdaftar di perusahaan ini")

	ErrDivisionNotFound   = errors.New("divisi tidak ditemukan")
	ErrDivisionCodeExists = errors.New("kode divisi sudah digunakan di perusahaan ini")
)
