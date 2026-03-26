package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordSuite(t *testing.T) {
	password := "rahasia123"

	t.Run("Hashing & Verification Success", func(t *testing.T) {
		hash, err := HashPassword(password)
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)

		// Menangkap 2 return value sesuai tanda tangan fungsi Anda
		match, err := CheckPasswordHash(password, hash)
		assert.NoError(t, err)
		assert.True(t, match)
	})

	t.Run("Wrong Password Failure", func(t *testing.T) {
		hash, _ := HashPassword(password)
		match, err := CheckPasswordHash("salah_dong", hash)
		assert.NoError(t, err)
		assert.False(t, match)
	})

	t.Run("Invalid Hash Format", func(t *testing.T) {
		// Mengetes skenario error sistem (menaikkan coverage error handling)
		match, err := CheckPasswordHash(password, "bukan-hash-valid")
		assert.Error(t, err)
		assert.False(t, match)
	})
}
