package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"runtime"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2 Config (Sesuai rekomendasi OWASP)
type argonConfig struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
}

var config = argonConfig{
	time:    1,                       // Iterasi
	memory:  64 * 1024,               // 64MB RAM
	threads: uint8(runtime.NumCPU()), // Menggunakan jumlah core CPU yang tersedia
	keyLen:  32,
	saltLen: 16,
}

// HashPasswordArgon2 menghasilkan hash string format $argon2id$v=19$m=65536,t=1,p=4$...
func HashPassword(password string) (string, error) {
	// 1. Generate random salt
	salt := make([]byte, config.saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// 2. Generate Hash
	hash := argon2.IDKey([]byte(password), salt, config.time, config.memory, config.threads, config.keyLen)

	// 3. Encode ke format string agar mudah disimpan di DB (Base64)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, config.memory, config.time, config.threads, b64Salt, b64Hash)

	return encodedHash, nil
}

// CheckPasswordHash memvalidasi plain password dengan hash dari DB
func CheckPasswordHash(password, encodedHash string) (bool, error) {
	if len(password) > 128 {
		return false, errors.New("password too long")
	}

	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}

	// PERBAIKAN: Ambil version dari parts[2] (v=19)
	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return false, err
	}

	if version != argon2.Version {
		return false, errors.New("incompatible argon2 version")
	}

	var memory, time uint32
	var threads uint8
	// Parse parameter m, t, p dari parts[3]
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false, err
	}

	// ... (sisa kode salt decode dan comparison tetap sama)
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	comparisonHash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(decodedHash)))

	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1, nil
}
