# go_server

Production-ready REST API built with **Go Fiber**, **PostgreSQL (Raw SQL)**, dan **pgxpool**. Mengikuti prinsip **Clean Architecture** untuk skalabilitas dan kemudahan pengujian.

---

## 🚀 Fitur & Keunggulan

- **Security First**: Menggunakan Raw SQL dengan *prepared statements* (mencegah SQL Injection).
- **Password Hashing**: Menggunakan **Argon2id** (lebih aman dari bcrypt).
- **Security Audit**: Integrasi **gosec** untuk pemindaian kerentanan kode secara otomatis.
- **Robust Testing**: Integration tests menggunakan **Testcontainers** (Docker-based) untuk database yang terisolasi.
- **Strict Layering**: Pemisahan logika yang jelas antara `handlers`, `services`, dan `types` di dalam direktori `internal/`.

---

## 🛠️ Tech Stack

- **Runtime**: Go 1.26.x
- **Framework**: [Fiber v2](https://gofiber.io/)
- **Database**: PostgreSQL (via [pgx v5](https://github.com/jackc/pgx))
- **Migration**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **Quality**: [golangci-lint](https://golangci-lint.run/), [gosec](https://securego.io/)
- **Testing**: [testify](https://github.com/stretchr/testify), [testcontainers-go](https://golang.testcontainers.org/)

---

## 📁 Project Structure

```
go_server/
├── main.go                          # Entry point (Bootstrap & Config)
├── Makefile                         # Unified Task Runner
├── internal/                        # Private code (Clean Architecture)
│   ├── app.go                       # Fiber Application Setup
│   ├── config/                      # Database Pool & Environment
│   ├── handlers/                    # HTTP Layer (Payload validation & response)
│   ├── services/                    # Business Logic Layer (SQL operations)
│   ├── types/                       # DTOs & Domain Models
│   ├── middlewares/                 # Auth, Logging, Rate Limit, Recovery
│   ├── routes/                      # Route Definitions
│   └── utils/                       # Security & Response Helpers
├── test/                            # Integration Tests (Testcontainers)
└── migrations/                      # SQL Migrations (.up.sql & .down.sql)
```

---

## 🚀 Getting Started

### Prerequisites

- Go 1.26+
- Docker (untuk menjalankan database lokal & integration tests)
- `make`

### 1. Setup Environment
```bash
cp .env.example .env
```

### 2. Jalankan Infrastruktur & Migrasi
```bash
# Start DB via Docker Compose (Opsional)
make compose-up

# Jalankan migrasi
make migrate-up
```

### 3. Jalankan Server
```bash
# Development (Hot Reload via Air)
make dev

# Production Build
make build && make run
```

---

## 📜 Available Commands

| Command | Deskripsi |
|---------|-----------|
| `make dev` | Jalankan server dengan hot-reload (Air) |
| `make security` | **Audit keamanan** menggunakan gosec (Issues: 0 baseline) |
| `make lint` | Analisis statis kode (golangci-lint) |
| `make test` | Jalankan unit & integration tests |
| `make test-cover` | Jalankan test + buka laporan coverage (Browser) |
| `make migrate-up` | Sinkronisasi skema database ke versi terbaru |
| `make fmt` | Format kode sesuai standar Go |

---

## 🛡️ Security Audit

Proyek ini mewajibkan **0 security issues**. Sebelum melakukan push, pastikan Anda menjalankan:
```bash
make security
```
Kami mengecualikan `.go_cache` untuk menghindari noise dari dependensi eksternal, fokus hanya pada keamanan kode internal.

---

## 📡 API Endpoints

### 🔐 Authentication
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| POST | `/api/v1/auth/register` | Baris pengguna baru (Argon2id) |
| POST | `/api/v1/auth/login` | Login & dapatkan JWT |

### 🏢 Companies (CRUD)
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/api/v1/companies` | List companies (Pagination) |
| POST | `/api/v1/companies` | Buat company baru |
| GET | `/api/v1/companies/:id` | Detail company |
| PUT | `/api/v1/companies/:id` | Update data company |
| DELETE | `/api/v1/companies/:id` | Hapus company |

### 🔍 Health & Meta
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/health` | Status API & Database Pool |

---

## 🧪 Testing Strategy

Kami menggunakan **Integration Testing** yang berat untuk menjamin kualitas:
- **Testcontainers**: Menjalankan PostgreSQL asli dalam container sementara untuk setiap session test.
- **Bullet-Proof**: Test mencakup skenario sukses, error validasi, hingga duplikasi data.

Jalankan test dengan:
```bash
make test
```

---

## License

MIT
