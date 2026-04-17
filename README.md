# 🚀 Backend Core Server (Go Fiber + PostgreSQL)

Selamat datang di repositori Backend Core Server! Proyek ini dibangun menggunakan **Go Fiber** yang super cepat dan
**PostgreSQL (pgxpool)**. Arsitekturnya dirancang menggunakan **Clean Architecture** berlapis yang **sangat mudah
dipahami**, **skalabel**, dan siap untuk level skala _Enterprise_ maupun _Production_.

---

## 📋 Persyaratan Sistem (Requirements)

Sebelum memulai, pastikan Anda telah menginstal beberapa perangkat lunak berikut di dalam mesin Anda:

1. **[Go v1.26+](https://go.dev/dl/)**: Bahasa pemrograman utama yang digunakan.
2. **[Docker Desktop / Colima](https://www.docker.com/) (Opsional)**: Sangat direkomendasikan untuk eksekusi
   _Integration Testing_ otomatis (`testcontainers`) dan menghidupkan Database lokal tanpa instalasi manual.
   - _Alternatif tanpa Docker_: Anda **tetap bisa menjalankan proyek ini tanpa Docker!** Cukup instal PostgreSQL (cth:
     Postgres.app / pgAdmin) di komputer Anda sendiri atau gunakan Cloud DB (seperti Supabase/Neon), lalu ubah link
     koneksi `DATABASE_URL` pada file `.env`.
3. **[Make](https://www.gnu.org/software/make/)**: _Task runner_ yang menyederhanakan perintah-perintah terminal
   kompleks. (Sudah bawaan di MacOS/Linux. Windows bisa menggunakan Git Bash atau WSL).

---

## 📥 Cara Cloning Project

Langkah pertama adalah menyalin repositori ini ke dalam mesin lokal Anda:

```bash
# 1. Clone repository
git clone https://github.com/yourusername/go_server.git

# 2. Masuk ke dalam direktori project
cd go_server
```

---

## 🏃‍♂️ Step-by-Step Menjalankan Project

Ikuti 3 langkah mudah di bawah ini untuk menjalankan server dari nol.

### Langkah 1: Persiapan Environment

Aplikasi ini membutuhkan file `.env` sebagai sumber konfigurasi utama (port, koneksi database, secret key, dll).

```bash
# Gandakan file template .env.example menjadi .env
cp .env.example .env
```

_(Opsional: Buka file `.env` dan sesuaikan nilainya jika diperlukan, namun secara default konfigurasi ini sudah siap
digunakan)._

### Langkah 2: Unduh Dependencies & Nyalakan Database

Aplikasi berjalan ditopang dengan PostgreSQL.

```bash
# Unduh library dan dependencies Go
make install
```

> 🔔 **Catatan Database (Tanpa Docker):** Jika Anda **tidak menggunakan Docker**, pastikan PostgreSQL di komputer Anda
> (atau Cloud) sudah menyala, dan link di dalam `.env` (`DATABASE_URL=postgres://...`) sudah mengarah ke database
> tersebut. **Abaikan perintah compose-up di bawah.**

Jika Anda **menggunakan Docker**, Anda cukup memutar perintah otomatis ini untuk menghidupkan PostgreSQL di background:

```bash
# (Khusus Pengguna Docker) Nyalakan PostgreSQL via Docker Compose
make compose-up
```

### Langkah 3: Migrasi Tabel & Jalankan Server!

Setelah database menyala, kita perlu men-generate struktur tabelnya (Migrasi). Setelah itu, nyalakan _development
server_.

```bash
# Buat tabel-tabel di database secara otomatis
make migrate-up

# Jalankan server mode development (dilengkapi Hot-Reload!)
make dev
```

🎉 **Selesai!** Server kini berjalan di `http://127.0.0.1:3000`. Coba akses endpoint `http://127.0.0.1:3000/health` dari
browser/Postman Anda!

---

## 💡 Arsitektur yang Mudah Dikelola (Penjelasan Perintah)

Untuk mempertahankan kualitas kode `(clean code)` dan kemudahan developer (_Developer Experience_), kami telah mengemas
semua _workflow_ kompleks ke dalam satu-kata ajaib yaitu **`make`**.

Anda tidak perlu menghafal panjangnya perintah terminal Go. Silakan lihat betapa rapinya pengelolaan infrastruktur
melalui perintah-perintah berikut:

### 👨‍💻 Development & Proses Build

- **`make dev`** — Menyalakan server dengan kemampuan _hot-reload_. Anda cukup klik _Save_ (`Ctrl+S`) di code editor,
  dan server akan otomatis _restart_ (memakai modul Air).
- **`make build`** — Melakukan kompilasi kode (`build`) yang telah di-optimisasi khusus untuk di-deploy ke environment
  sistem operasi Production.
- **`make run`** — Melakukan build dan langsung mengeksekusi file _binary_ server aplikasi rilis.

### 🗄️ Database & Migrasi

- **`make migrate-create name=tambah_tabel_x`** — Membuat file SQL migrasi _up & down_ baru (kosong).
- **`make migrate-up`** — Menerapkan/Menjalankan file SQL migrasi terbaru ke dalam database.
- **`make migrate-down`** — Mengembalikan (_Rollback_) satu versi langkah migrasi ke belakang. Sangat berguna jika salah
  membuat tabel saat ngoding.

### 🧪 Kualitas dan Keamanan Kode (Tests & Security)

Struktur kode ini menerapkan konsep _Bullet-Proof_ atau tak-tertembus. Anda dapat memverifikasi kualitas penulisan kode
atau mengecek celah kerentanan dalam sekejap mata:

- **`make test-cover`** — Menjalankan _Integration Test Suite_ berlapis menggunakan Docker-Testcontainers + menampilkan
  **persentase laporan cakupan kode (Coverage)**. Target proyek kita selalu mempertahankan cakupan keamanan di angka
  rata-rata `85%` ke-atas.
- **`make lint`** — Merapikan dan menganalisis kesalahan-kesalahan penulisan (typo, import tidak terpakai, deklarasi
  redundant) standar baku golang via _golangci-lint_.
- **`make security`** — 🛡️ **FITUR UNGGULAN.** Meng-audit seluruh repositori baris per baris secara otomatis untuk
  mendeteksi _Hardcoded-password_, potensi kerentanan SQL Injection, dan celah otentikasi.

### 🐳 Manajemen Ekosistem Docker

- **`make compose-up`** — Menghidupkan seluruh layanan infra (Database Postgres) background.
- **`make compose-down`** — Menyapu bersih dan mematikan infrastruktur.

---

## ⚖️ Docker vs Tanpa Docker (Pros & Cons)

Tim bebas mengadopsi lingkungan mereka sendiri. Berikut bahan pertimbangan untuk menentukan alur mana yang lebih cocok
untuk Anda:

### 🟢 Menggunakan Docker (Direkomendasikan)

Cocok untuk lingkungan tim kolaboratif dan standar industri:

- **Positif (+):**
  - _Zero-Config Setup_: Tidak perlu mengotori sistem mac/windows Anda dengan instalasi database yang berat.
  - _Lingkungan Identik_: Semua _developer_ di tim memiliki spesifikasi Postgres yang sama tanpa takut bentrok versi.
  - **Akses Penuh Fitur Testing**: Bisa menjalankan `make test-cover` secara utuh karena _Testcontainers_ mensyaratkan
    ketersediaan daemon Docker untuk membuang dan membangun container uji secara _on-the-fly_.
- **Negatif (-):**
  - Memakan lebih banyak _RAM_ (karena overhead visualisasi mesin/Container).

### 🟠 Tanpa Docker (Local Native)

Cocok jika laptop Anda memiliki sumber daya _Low-End/Memori terbatas_ atau tidak punya izin akses Root/Admin:

- **Positif (+):**
  - _Performa Super Cepat & Ringan_: Mengetik dan mengeksekusi aplikasi jauh lebih ringan pada memori komputer karena
    koneksi database _native_ secara langsung (tanpa melintasi _virtual network layer_).
  - Nyaman jika Anda sudah sangat familiar dengan _pgAdmin / DBeaver_.
- **Negatif (-):**
  - **Tidak Bisa Menjalankan Testcontainers**: Fitur `make test-cover` akan langsung Error ("_Docker not found_"). Anda
    harus mengakalinya secara repot dengan melakukan sinkronisasi database manual khusus testing dan melakukan _Truncate
    Table_ setiap kali sebelum perintah tes di-klik.

---

## 📁 Struktur Direktori Singkat

Arsitektur kita yang berbasis `Clean Architecture` sangat tegas memisahkan logika lalu lintas agar sistem tidak
berantakan:

```text
go_server/
├── internal/             # Jantung dari aplikasi
│   ├── handlers/         # Titik terima & validasi Body Request HTTP (Hanya mengecek input)
│   ├── services/         # Logika Utama Pemrosesan Data & Operasi Transaksi SQL
│   ├── types/            # DTO (Data Transfer Objects), Struct, Bentuk Skema Data
│   ├── middlewares/      # Interceptor (Autentikasi, Batasan Kuota Rate-Limit, CORS)
│   └── routes/           # Peta Alamat URL Endpoint (Router API)
├── migrations/           # Riwayat pembuatan tabel SQL berseri
├── test/integration/     # Rumah dari segala alat pengetesan otomatis Docker
└── Makefile              # Daftar menu kendali "make" (Command Center)
```

Selamat membangun sistem yang menakjubkan! Jika ragu dengan langkah apa saja, cukup ketik **`make`** atau
**`make help`** di terminal untuk melihat daftar bantuan secara instan.
