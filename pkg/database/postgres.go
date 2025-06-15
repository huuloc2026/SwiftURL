package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDB() *sqlx.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getEnv("DB_USER", "myuser"),
		getEnv("DB_PASS", "mypassword"),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "3306"),
		getEnv("DB_NAME", "mydatabase"),
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to connect to PostgreSQL: %v", err)
	}

	// Chạy init.sql
	if err := runMigration(db); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	DB = db
	return db
}

func runMigration(db *sqlx.DB) error {
	fmt.Println("📦 Running init.sql migration...")
	sql, err := os.ReadFile("migrations/init.sql")
	if err != nil {
		return fmt.Errorf("cannot read init.sql: %w", err)
	}
	_, err = db.Exec(string(sql))
	if err != nil {
		return fmt.Errorf("cannot exec init.sql: %w", err)
	}
	fmt.Println("✅ Database initialized.")
	return nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
