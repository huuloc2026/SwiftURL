# 🚀 SwiftURL - URL Shortener in Go

SwiftURL is a clean, modular URL shortener service built with:

- [Fiber](https://github.com/gofiber/fiber) web framework
- [SQLX](https://github.com/jmoiron/sqlx) for PostgreSQL
- Clean Architecture principles
- Simple SQL migrations

---

## 📦 Features

- Generate short links for any valid URL
- Resolve short links with redirect
- Track click count and metadata
- User registration, login, password reset (with OTP)
- Modular, testable codebase

---

## 📁 Project Structure

```
SwiftURL/
├── cmd/server/                  # Main entry point
├── config/                      # Environment/config loading
├── internal/
│   ├── entity/                  # Shared entity definitions
│   ├── shorturl/                # ShortURL module (delivery, repository, usecase)
│   ├── user/                    # User module (delivery, repository, usecase)
│   └── auth/                    # Auth module (delivery, usecase)
├── pkg/
│   ├── cache/                   # Redis cache
│   ├── database/                # DB initialization/migration
│   ├── response/                # Standard API responses
│   └── utils/                   # Utilities (short code generator, etc.)
├── migrations/                  # SQL schema
├── tests/                       # Integration/unit tests
├── Dockerfile
├── docker-compose.yml
├── .env.example
└── README.md
```

---

## ⚙️ Requirements

- Go 1.20+
- PostgreSQL 13+
- Redis 7.2+

---

## 🛠️ Setup

1. **Clone the repository**

   ```bash
   git clone https://github.com/huuloc2026/SwiftURL.git
   cd SwiftURL
   ```

2. **Configure environment**

   Copy `.env.example` to `.env` and update values as needed.

   ```bash
   cp .env.example .env
   ```

3. **Start dependencies (Postgres, Redis) with Docker Compose**

   ```bash
   docker compose up -d
   ```

4. **Run migrations and start the app**

   ```bash
   go run cmd/server/main.go
   ```

   Or use Air for live reload:

   ```bash
   air
   ```

---

## 🧪 API Endpoints

### Health Check

```
GET /healthz
```

### Shorten URL

```
POST /api/shorten
```
**Body:**
```json
{
  "long_url": "https://example.com"
}
```
**Response:**
```json
{
  "short_code": "aB12Cd"
}
```

### Redirect Short URL

```
GET /:code
```
Redirects to the original URL.

### User Authentication

- `POST /api/auth/register`
- `POST /api/auth/login`
- `POST /api/auth/forget-password`
- `POST /api/auth/verify-otp`
- `POST /api/auth/change-password`

See the [API usage](#api-usage) section for example requests.

---

## 🧰 Tips

- All configuration is via `.env`
- Migrations run automatically on startup
- Use `/api/auth/*` for authentication endpoints

---

## 🧪 API Usage (Postman Examples)

**Register:**
```
POST /api/auth/register
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "yourpassword"
}
```

**Login:**
```
POST /api/auth/login
{
  "email": "test@example.com",
  "password": "yourpassword"
}
```

**Forget Password:**
```
POST /api/auth/forget-password
{
  "email": "test@example.com"
}
```

**Verify OTP:**
```
POST /api/auth/verify-otp
{
  "email": "test@example.com",
  "otp": "123456"
}
```

**Change Password:**
```
POST /api/auth/change-password
{
  "email": "test@example.com",
  "otp": "123456",
  "new_password": "newpassword"
}
```

---

## 📄 License

MIT © 2025 by [huuloc2026](https://github.com/huuloc2026)



