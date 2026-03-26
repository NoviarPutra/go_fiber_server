package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yourusername/go_server/internal/services"
)

type UsersServiceTestSuite struct {
	suite.Suite
	service *services.UsersService
}

func (s *UsersServiceTestSuite) SetupSuite() {
	s.Require().NotNil(testDBPool, "DB Pool harus sudah siap dari TestMain")
	s.service = services.NewUsersService(testDBPool)
}

// TearDownTest memastikan setiap sub-test mulai dengan DB bersih
func (s *UsersServiceTestSuite) TearDownTest() {
	_, err := testDBPool.Exec(context.Background(), "TRUNCATE users CASCADE")
	s.Require().NoError(err)
}

// Helper Seeding: Menambahkan data dummy
func (s *UsersServiceTestSuite) seedUsers(ctx context.Context, count int) {
	for i := 1; i <= count; i++ {
		// Gunakan formatting yang benar tanpa karakter aneh (non-breaking space)
		query := `INSERT INTO users (email, username, password_hash, is_active) 
				  VALUES ($1, $2, 'dummy_hash', true)`
		_, err := testDBPool.Exec(ctx, query,
			fmt.Sprintf("user%d@officecore.id", i),
			fmt.Sprintf("user%d", i),
		)
		s.Require().NoError(err)
	}
}

// ─── TEST CASES ──────────────────────────────────────────────────────────────

func (s *UsersServiceTestSuite) TestGetUsers_Pagination() {
	ctx := context.Background()
	totalInserted := 15
	s.seedUsers(ctx, totalInserted)

	s.Run("Normal_Pagination_Page_1", func() {
		// Ambil 5 data pertama
		users, total, err := s.service.GetUsers(ctx, 1, 5)

		s.NoError(err)
		s.Equal(int64(totalInserted), total)
		s.Len(users, 5)
		// Karena ORDER BY created_at DESC, user terakhir masuk harusnya muncul pertama
		s.Equal("user15@officecore.id", users[0].Email)
	})

	s.Run("Last_Page_Pagination", func() {
		// Page 3 dengan limit 5 (data 11-15)
		users, total, err := s.service.GetUsers(ctx, 3, 5)

		s.NoError(err)
		s.Equal(int64(totalInserted), total)
		s.Len(users, 5)
		s.Equal("user5@officecore.id", users[0].Email)
	})
}

func (s *UsersServiceTestSuite) TestGetUsers_EdgeCases() {
	ctx := context.Background()

	s.Run("Empty_Database", func() {
		// Pastikan saat DB kosong, total_count adalah 0 dan users adalah slice kosong
		users, total, err := s.service.GetUsers(ctx, 1, 10)

		s.NoError(err)
		s.Equal(int64(0), total)
		s.NotNil(users)
		s.Empty(users)
	})

	s.Run("Offset_Out_Of_Range", func() {
		s.seedUsers(ctx, 5) // Data cuma 5

		// Request Page 2 dengan Limit 10 (Offset 10)
		// Ini akan memicu logic `if len(users) == 0` di service Anda
		users, total, err := s.service.GetUsers(ctx, 2, 10)

		s.NoError(err)
		s.Equal(int64(5), total, "Total data harusnya tetap 5 walaupun halaman ini kosong")
		s.Empty(users, "Daftar user harusnya kosong karena melewati batas data")
		s.NotNil(users, "Slice tidak boleh nil")
	})
}

func TestUsersServiceIntegration(t *testing.T) {
	suite.Run(t, new(UsersServiceTestSuite))
}
