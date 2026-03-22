# go_server

Production-ready REST API built with **Go Fiber**, **PostgreSQL**, dan **pgxpool**.

---

## Tech Stack

- **Runtime**: Go 1.24+
- **Framework**: [Fiber v2](https://gofiber.io/)
- **Database**: PostgreSQL (via [pgx v5](https://github.com/jackc/pgx))
- **Migration**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **Hot Reload**: [Air](https://github.com/air-verse/air)

---

## Project Structure
```
go_server/
├── main.go                          # Entry point — load env, connect DB, start server
├── app.go                           # Bootstrap Fiber app, middleware stack, routes
├── .env                             # Environment variables (git-ignored)
├── .env.example                     # Environment template
├── Makefile                         # Task runner
│
├── config/
│   └── db.go                        # PostgreSQL connection pool (pgxpool)
│
├── handlers/                        # HTTP layer — terima request, kembalikan response
│   ├── base/
│   │   └── base_handler.go          # Handler dasar (welcome, 404)
│   ├── health/
│   │   └── health_handler.go        # Health check + DB ping
│   ├── users/
│   │   └── users_handler.go         # User endpoints
│   └── errors/
│       └── error_handler.go         # Global error handler (Fiber ErrorHandler)
│
├── middlewares/                     # Middleware stack
│   ├── auth_middleware.go           # JWT authentication (Protected())
│   ├── cors_middleware.go           # CORS (dev: *, prod: ALLOWED_ORIGINS)
│   ├── db_middleware.go             # Inject *pgxpool.Pool ke fiber.Ctx
│   ├── logger_middleware.go         # Request logging + JSON pretty print
│   ├── pagination_middleware.go     # Query param: page, limit
│   ├── rate_limit_middleware.go     # Rate limiting (dev: 1000/min, prod: 100/min)
│   └── recovery_middleware.go       # Panic recovery (stack trace di dev only)
│
├── routes/
│   └── routes.go                    # Definisi semua route & grouping
│
├── utils/
│   └── response.go                  # Standard response helpers (Success, Error, dll)
│
├── types/
│   └── types.go                     # Shared struct types (StandardResponse, Meta, dll)
│
└── migrations/                      # SQL migration files (golang-migrate)
    ├── 000001_create_users.up.sql
    └── 000001_create_users.down.sql
```

---

## Getting Started

### Prerequisites

- Go 1.24+
- PostgreSQL
- `make` (untuk task runner)

### 1. Clone & Install
```bash
git clone https://github.com/yourname/go_server.git
cd go_server
make install
```

### 2. Setup Environment
```bash
cp .env.example .env
```

Edit `.env`:
```properties
PORT=3001
APP_ENV=development
JWT_SECRET=ganti-dengan-secret-yang-kuat

DB_USER=postgres
DB_PASS=postgres
DB_HOST=127.0.0.1
DB_PORT=5432
DB_NAME=db_go_server
DB_SSLMODE=disable
DATABASE_URL=postgres://postgres:postgres@127.0.0.1:5432/db_go_server?sslmode=disable

ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

### 3. Jalankan Migrasi
```bash
make migrate-up
```

### 4. Jalankan Server
```bash
# Development — hot reload
make dev

# Production
make run
```

Server berjalan di `http://localhost:3001`

---

## Available Commands
```bash
make install                      # Install dependencies
make build                        # Build binary (optimized)
make run                          # Build & run
make dev                          # Hot reload — auto install Air jika belum ada
make test                         # Run tests dengan -race flag
make test-cover                   # Run tests + generate coverage.html
make clean                        # Hapus binary & artifacts
make fmt                          # Format kode
make lint                         # Lint kode — auto install golangci-lint jika belum ada
make vet                          # Run go vet
```
```bash
make migrate-create name=<name>   # Buat file migration baru
make migrate-up                   # Jalankan semua pending migrations
make migrate-down                 # Rollback 1 migration
make migrate-force version=<ver>  # Force set versi migration
make db-status                    # Cek versi migration saat ini
```
```bash
make docker-build                 # Build Docker image
make docker-run                   # Run di Docker container
make compose-up                   # Start dengan Docker Compose
make compose-down                 # Stop Docker Compose
make compose-build                # Build Docker Compose services
```

---

## API Endpoints

### Public

| Method | Endpoint      | Deskripsi                        |
|--------|---------------|----------------------------------|
| GET    | `/`           | Welcome message                  |
| GET    | `/health`     | Health check + DB pool stats     |

### Private (perlu JWT token)

| Method | Endpoint         | Deskripsi                       |
|--------|------------------|---------------------------------|
| GET    | `/api/v1/users`  | Get users (support pagination)  |

### Headers untuk Private Endpoint
```
Authorization: Bearer <token>
```

### Contoh Response
```json
{
  "success": true,
  "message": "Data berhasil diambil",
  "data": [...],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 100
  }
}
```
```json
{
  "success": false,
  "message": "Token diperlukan untuk mengakses resource ini",
  "data": null,
  "meta": null
}
```

---

## Environment Variables

| Variable          | Default       | Keterangan                                   |
|-------------------|---------------|----------------------------------------------|
| `PORT`            | `3001`        | Port server                                  |
| `APP_ENV`         | `development` | Mode: `development` atau `production`        |
| `JWT_SECRET`      | —             | Secret key untuk signing JWT token           |
| `DATABASE_URL`    | —             | PostgreSQL connection string                 |
| `ALLOWED_ORIGINS` | —             | CORS origins (prod), `*` otomatis saat dev   |

---

## License

MIT
