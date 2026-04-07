package services

import "errors"

var (
	ErrEmailAlreadyExists    = errors.New("email sudah terdaftar")
	ErrUsernameAlreadyExists = errors.New("username sudah digunakan")
	ErrUserNotFound          = errors.New("user tidak ditemukan")
	ErrInvalidCredentials    = errors.New("email atau password salah")
	ErrAccountInactive       = errors.New("akun tidak aktif")
	ErrRefreshTokenInvalid   = errors.New("refresh token tidak valid")
	ErrRefreshTokenExpired   = errors.New("refresh token kadaluarsa")
	ErrRefreshTokenRevoked   = errors.New("refresh token sudah digunakan")

	ErrCompanyNotFound       = errors.New("perusahaan tidak ditemukan")
	ErrCompanyCodeExists     = errors.New("kode perusahaan sudah digunakan")
)
