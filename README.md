# Gamification App — Backend API

Backend service สำหรับเกมสะสมคะแนน พัฒนาด้วย GoLang + Fiber v2

**Live API:** https://gamification-app-qie0.onrender.com
**Swagger UI:** https://gamification-app-qie0.onrender.com/swagger/index.html

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.25 |
| HTTP Framework | Fiber v2 |
| ORM | GORM |
| Database | PostgreSQL |
| Auth | JWT (golang-jwt/jwt/v5) |
| Logger | zerolog |
| Validation | go-playground/validator/v10 |
| API Docs | Swagger (swaggo/swag) |

---

## Architecture

ใช้ **Clean Architecture แบบ Feature-Modular** — แต่ละ feature เป็น folder ของตัวเองใน `internal/`

```
internal/{feature}/
├── domain/        # Entity, Interface, Sentinel Errors (innermost — ไม่ depend ใคร)
├── repository/    # GORM implementation
│   └── models/    # GORM models
├── service/       # Business Logic
└── handler/       # HTTP Handlers + Router
    └── dto/       # Request/Response DTOs
```

### Dependency Direction

```
handler → service → domain ← repository
```

- **Domain** คือ contract กลาง — ทุก layer depend บน interface ที่นี่
- **Repository** implement `domain.XxxRepository`
- **Service** depend บน interface ไม่ใช่ struct concrete
- **Handler** depend บน `domain.XxxService` interface

### Directory Structure

```
cmd/
  app/main.go          # Entry point + graceful shutdown
  server/server.go     # Fiber setup, middleware, DB wiring
config/
  config.go            # Config struct โหลดจาก env
internal/
  users/               # Auth + Profile
  games/               # Spin wheel
  games-histories/     # Spin history
  rewards/             # Rewards + Claim
pkg/
  auth/                # JWT generate/validate
  common/              # Response, Logger helpers
  constants/           # Error codes, Messages
  database/            # PostgreSQL connection
  middleware/          # JWT middleware
migrations/
  schema/              # SQL schema files
  seed/                # CSV import script
```

---

## Features

### Authentication
- `POST /auth/login` — Login ด้วย nickname (ถ้าไม่มีในระบบ สร้างอัตโนมัติ)

### Users
- `GET /api/v1/users` — ดู users ทั้งหมด (filter by nickname, pagination)
- `GET /api/v1/users/me` — ดู profile ตัวเอง + rewards status (parallel fetch)

### Spin Game
- `POST /api/v1/games/spin` — หมุนวงล้อรับคะแนน (300 / 500 / 1000 / 3000)
  - คะแนนสูงสุด 10,000 — เกินไม่ได้

### Rewards
- `GET /api/v1/rewards` — ดู rewards ทั้งหมด (checkpoint 500 / 1000 / 10000)
- `POST /api/v1/rewards/claim` — แลกรางวัล
- `GET /api/v1/user-rewards` — ดูรางวัลที่ claim แล้ว

### Game Histories
- `GET /api/v1/games-histories` — ประวัติการเล่นทั้งหมด (filter by nickname, pagination, sort)

---

## Installation

### Prerequisites
- Go 1.21+
- PostgreSQL
- Docker (optional)

### 1. Clone

```bash
git clone https://github.com/ohmtanawin02/gamification-app.git
cd gamification-app
```

### 2. Setup Environment

```bash
cp .env.example .env
```

แก้ค่าใน `.env`:

```env
PORT=9393
ENVIRONMENT=development
DEBUG=true

DB_WRITE_HOST=localhost
DB_WRITE_PORT=5432
DB_WRITE_USER=postgres
DB_WRITE_PASSWORD=postgres
DB_WRITE_DB_NAME=gamification_db

DB_READ_HOST=localhost
DB_READ_PORT=5432
DB_READ_USER=postgres
DB_READ_PASSWORD=postgres
DB_READ_DB_NAME=gamification_db

JWT_SECRET=your-secret-key
JWT_EXPIRE_HOURS=24
```

### 3. Start Database

```bash
docker compose up -d
```

### 4. Run Migrations

```bash
psql "postgres://postgres:postgres@localhost:5432/gamification_db" \
  -f migrations/schema/001_create_users.sql \
  -f migrations/schema/002_create_user_game_histories.sql \
  -f migrations/schema/003_create_rewards.sql \
  -f migrations/schema/004_create_user_rewards.sql
```

### 5. Import CSV Data (optional)

```bash
psql "postgres://postgres:postgres@localhost:5432/gamification_db" \
  -c "CREATE TABLE csv_staging (nickname VARCHAR, point INT, datetime TIMESTAMPTZ);"

psql "postgres://postgres:postgres@localhost:5432/gamification_db" \
  -c "\copy csv_staging FROM 'mock_data.csv' WITH (FORMAT csv, HEADER true)"

psql "postgres://postgres:postgres@localhost:5432/gamification_db" \
  -f migrations/seed/import_csv.sql
```

### 6. Run Server

```bash
go mod tidy
go run cmd/app/main.go
```

หรือใช้ hot reload:

```bash
go install github.com/air-verse/air@latest
air
```

Server รันที่ `http://localhost:9393`

---

## API Documentation

Swagger UI: `http://localhost:9393/swagger/index.html`

### Authentication

ทุก endpoint ใต้ `/api/v1` ต้องใช้ JWT

```
Authorization: Bearer <token>
```

รับ token จาก `POST /auth/login`

---

## Error Codes

| HTTP | Code | Description |
|------|------|-------------|
| 200 | OK | Success |
| 400 | BAD_REQUEST | Invalid request |
| 401 | UNAUTHORIZED | Missing or invalid token |
| 404 | NOT_FOUND | Resource not found |
| 409 | CONFLICT | Already claimed |
| 422 | UNPROCESSABLE_ENTITY | Business rule violation |
| 500 | INTERNAL_SERVER_ERROR | Server error |
