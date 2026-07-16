# Go GORM Auto Migration

Project Go sederhana yang menggunakan GORM dan PostgreSQL. Saat aplikasi
dijalankan, GORM otomatis membuat atau memperbarui tabel berdasarkan model yang
terdaftar. REST API menggunakan Gin dengan pemisahan repository, service, DTO,
handler, dan route.

## Menyiapkan PostgreSQL

Pastikan PostgreSQL aktif dan database sudah dibuat:

```sql
CREATE DATABASE payment;
```

Konfigurasi koneksi menggunakan environment berikut:

```env
APP_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=payment
DB_SSLMODE=disable
DB_TIMEZONE=Asia/Singapore
DB_DEBUG=true
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=30m
DB_CONN_MAX_IDLE_TIME=5m
```

Nilai tersebut juga tersedia di `.env.example`. Salin menjadi `.env`, lalu
sesuaikan dengan PostgreSQL lokal. Aplikasi membaca `.env` secara otomatis.

Contoh:

```bash
cp .env.example .env
go run ./cmd/migrate
go run ./cmd/api
```

Jika environment tidak diisi, aplikasi memakai nilai default di atas.

## Menjalankan migration

Migration dipisahkan dari startup API agar beberapa replica API tidak menjalankan
perubahan schema secara bersamaan.

Jalankan migration sebelum menjalankan atau mendeploy API:

```bash
go run ./cmd/migrate
```

Setelah migration berhasil:

```bash
go run ./cmd/api
```

## Log query global

`DB_DEBUG=true` mengaktifkan log untuk semua query yang dijalankan melalui
instance GORM. Log menampilkan SQL, nilai parameter, durasi eksekusi, jumlah row,
dan error database.

Contoh output:

```text
[GORM] 2026/07/16 20:00:00 /path/to/file.go:10
[1.245ms] [rows:1] SELECT * FROM "users" WHERE "users"."id" = 1
```

Matikan log query detail pada environment production:

```env
DB_DEBUG=false
```

Dalam mode tersebut, logger tetap mencatat query lambat dan error database.

## Database connection pool

Connection pool PostgreSQL dapat diatur melalui:

```env
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=30m
DB_CONN_MAX_IDLE_TIME=5m
```

Nilai durasi mengikuti format Go seperti `30s`, `5m`, atau `1h`.

## Struktur project

```text
.
├── cmd
│   ├── api
│   │   └── main.go    # entry point REST API
│   └── migrate
│       └── main.go    # command migration database
├── internal
│   ├── bootstrap      # dependency wiring dan application setup
│   ├── database       # koneksi dan migration
│   ├── dto            # request, response, dan pagination reusable
│   ├── handler        # validasi input dan HTTP response
│   ├── mapper         # konversi model ke DTO
│   ├── model          # entity GORM
│   ├── repository     # operasi database
│   ├── route          # registrasi endpoint Gin
│   └── service        # aturan bisnis
└── go.mod
```

Alur sebuah request:

```text
Client -> Route -> Handler -> Service -> Repository -> PostgreSQL
```

## Endpoint API

| Method | Endpoint | Keterangan |
|---|---|---|
| `GET` | `/health` | Health check (kompatibilitas) |
| `GET` | `/health/live` | Liveness check proses API |
| `GET` | `/health/ready` | Readiness check API dan PostgreSQL |
| `POST` | `/api/product` | Membuat produk |
| `GET` | `/api/product?page=1&limit=10` | Daftar produk |
| `GET` | `/api/product/:id` | Detail produk |
| `PATCH` | `/api/product/:id` | Memperbarui produk |
| `DELETE` | `/api/product/:id` | Menghapus produk |
| `POST` | `/api/category` | Membuat kategori |
| `GET` | `/api/category` | Daftar kategori |
| `GET` | `/api/category/:id` | Detail kategori |
| `GET` | `/api/category/:id/products` | Product dalam kategori |
| `PATCH` | `/api/category/:id` | Memperbarui kategori |
| `DELETE` | `/api/category/:id` | Menghapus kategori |

## Relasi Product

Model Product mempunyai beberapa jenis relasi:

| Relasi | Tipe | Foreign key/junction |
|---|---|---|
| Category ke Product | One-to-many | `products.category_id` |
| Product ke ProductDetail | One-to-one | `product_details.product_id` |
| Product ke ProductVariant | One-to-many | `product_variants.product_id` |
| Product ke ProductImage | One-to-many | `product_images.product_id` |
| Product ke Tag | Many-to-many | `product_tags` |

Category harus dibuat terlebih dahulu, kemudian ID-nya dapat dikirim melalui
`category_id`. Detail, variants, images, dan tags dapat dibuat bersamaan dengan
Product.

Membuat produk:

```bash
curl -X POST http://localhost:8080/api/product \
  -H "Content-Type: application/json" \
  -d '{
    "category_id": 1,
    "name": "Keyboard Mechanical",
    "description": "Keyboard switch linear",
    "price": 750000,
    "stock": 20,
    "detail": {
      "weight_grams": 850,
      "length_cm": 35,
      "width_cm": 13,
      "height_cm": 4,
      "brand": "Enigma",
      "sku": "KEYBOARD-001"
    },
    "variants": [
      {
        "name": "Black",
        "sku": "KEYBOARD-001-BLACK",
        "price": 750000,
        "stock": 10
      },
      {
        "name": "White",
        "sku": "KEYBOARD-001-WHITE",
        "price": 775000,
        "stock": 10
      }
    ],
    "images": [
      {
        "url": "https://example.com/keyboard.jpg",
        "alt_text": "Keyboard Mechanical",
        "is_primary": true,
        "position": 1
      }
    ],
    "tags": [
      {
        "name": "Mechanical",
        "slug": "mechanical"
      },
      {
        "name": "Keyboard",
        "slug": "keyboard"
      }
    ]
  }'
```

Memperbarui sebagian field produk:

```bash
curl -X PATCH http://localhost:8080/api/product/1 \
  -H "Content-Type: application/json" \
  -d '{
    "price": 700000,
    "stock": 15
  }'
```

Response API memakai format konsisten:

```json
{
  "success": true,
  "message": "produk berhasil diambil",
  "data": {
    "id": 1,
    "name": "Keyboard Mechanical"
  }
}
```

## Menambahkan model baru

1. Buat struct baru di `internal/model`.
2. Tambahkan model tersebut ke `db.AutoMigrate` di
   `internal/database/migrate.go`.
3. Buat repository, service, DTO, dan handler sesuai kebutuhan.
4. Daftarkan endpoint baru di `internal/route/route.go`.

Contoh:

```go
type Order struct {
	gorm.Model
	Total int64 `gorm:"not null"`
}
```

Kemudian daftarkan:

```go
return db.AutoMigrate(
	&model.User{},
	&model.Product{},
	&model.Order{},
)
```

`AutoMigrate` akan membuat tabel, kolom, index, dan constraint yang belum ada
tanpa menghapus kolom lama dari database.
