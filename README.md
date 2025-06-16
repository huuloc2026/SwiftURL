# 🚀 SwiftURL - URL Shortener in Golang

SwiftURL is a simple and clean URL shortener service written in **Go**, built with:
- [Fiber](https://github.com/gofiber/fiber) as the web framework
- [SQLX](https://github.com/jmoiron/sqlx) for PostgreSQL database access
- Clean Architecture principles
- Lightweight migration using raw SQL file (`init.sql`)

---

## 📁 Project Structure

```
shortener-app/
├── cmd/ # Main entry point
│ └── main.go
├── migrations/ # SQL file for DB schema initialization
│ └── init.sql
├── internal/ # Main application logic
│ ├── entity/ # Shared entity definitions
│ ├── shorturl/ # ShortURL module
│ │ ├── delivery/http/ # HTTP handlers
│ │ ├── repository/ # SQLX implementation
│ │ ├── usecase/ # Business logic
│ │ └── model.go
├── pkg/ # Shared packages (db, utils, etc.)
│ ├── database/postgres.go # DB initialization and migration
│ └── utils/generator.go # Short code generator
├── go.mod
└── README.md

```


---

## 📦 Features

- ✅ Generate short links for any valid long URL
- 🔗 Resolve short links with redirect support
- 📊 Track basic click count (optional)
- 🧱 SQL-based schema using `init.sql`
- ♻️ Modular and testable Clean Architecture

---

## ⚙️ Requirements

- Go 1.20+
- PostgreSQL 13+
- Git

---

## 🛠️ Setup Instructions

### 1. Clone the repository

```bash
git clone https://github.com/huuloc2026/SwiftURL.git
cd SwiftURL
```

Here is a detailed `README.md` for your **URL Shortener project** using **Go (Fiber + SQLX + PostgreSQL)** with a **Clean Architecture** structure and file-based migration using `init.sql`.

---

## 📘 README.md — SwiftURL: URL Shortener in Go


# 🚀 SwiftURL - URL Shortener in Golang

SwiftURL is a simple and clean URL shortener service written in **Go**, built with:
- [Fiber](https://github.com/gofiber/fiber) as the web framework
- [SQLX](https://github.com/jmoiron/sqlx) for PostgreSQL database access
- Clean Architecture principles
- Lightweight migration using raw SQL file (`init.sql`)



## 📁 Project Structure

```

shortener-app/
├── cmd/                       # Main entry point
│   └── main.go
├── migrations/               # SQL file for DB schema initialization
│   └── init.sql
├── internal/                 # Main application logic
│   ├── entity/               # Shared entity definitions
│   ├── shorturl/             # ShortURL module
│   │   ├── delivery/http/    # HTTP handlers
│   │   ├── repository/       # SQLX implementation
│   │   ├── usecase/          # Business logic
│   │   └── model.go
├── pkg/                      # Shared packages (db, utils, etc.)
│   ├── database/postgres.go  # DB initialization and migration
│   ├── cache/redis.go  # DB initialization and migration
│   ├── jwt/jwt.go  # DB initialization and migration
│   ├── response/response.go  # DB initialization and migration
│   └── utils/generator.go    # Short code generator
├── go.mod
├── .env.example.mod
├── Dockerfile.yml
├── docker-compose.yml
└── README.md
```

---

## 📦 Features

- ✅ Generate short links for any valid long URL
- 🔗 Resolve short links with redirect support
- 📊 Track basic click count (optional)
- 🧱 SQL-based schema using `init.sql`
- ♻️ Modular and testable Clean Architecture

---

## ⚙️ Requirements

- Go 1.20+
- PostgreSQL 13+
- Redis 7.2

---

## 🛠️ Setup Instructions

### 1. Clone the repository

```bash
git clone https://github.com/huuloc2026/SwiftURL.git
cd SwiftURL
```

---

### 2. Configure your PostgreSQL connection

Edit the connection string in `pkg/database/postgres.go`:

```go
dsn := "postgres://<user>:<password>@localhost:5432/<your_db>?sslmode=disable"
```

Alternatively, extract the DSN into environment variables.

---

### 3. Create database and run migration

Make sure your PostgreSQL instance is running and create the database:

```bash
createdb shorturl
```

The first time you run the app, it will automatically execute `migrations/init.sql`:

```bash
go run cmd/server/main.go
```

```bash
air
```

You'll see:

```
📦 Running init.sql migration...
✅ Database initialized.
✅ Redis connected:
```

---

## 🧪 API Endpoints

### 🔍 Health check

```
GET /healthz
```

Returns:

```json
{ "status": "ok" }
```

---

### ✂️ Shorten URL

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
  "short_code": "aB12Cd",
  "long_url": "https://example.com"
}
```

---

### 🔁 Redirect short URL

```
GET /api/:code
```

Example:

```
GET /api/aB12Cd → 301 Redirect → https://example.com
```

---

## 🧪 Testing with Postman

* Set method to `POST`
* URL: `/api/shorten`
* Headers:

  * `Content-Type: application/json`
* Body:

```json
{
  "url": "https://example.com"
}
```

---

## 🧰 Useful Tips

* Use `sqlx` with struct tags (`db:"field"`) for cleaner code.
* Migrations are applied only once at app startup (`init.sql`).
* Customize `short_code` length or charset in `pkg/utils/generator.go`.

---

## 🧪 Coming Soon (Ideas)

* [x] Expiry time for short links
* [ ] Admin panel with stats
* [ ] QR code generation
* [ ] Auth module (already scaffolded)

---

## 📄 License

MIT © 2025 by [huuloc2026](https://github.com/huuloc2026)



