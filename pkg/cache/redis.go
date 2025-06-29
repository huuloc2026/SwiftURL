package cache

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  InitRedis:: .env file not found, using system env")
	}
}
func InitRedis() *redis.Client {

	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	username := os.Getenv("REDIS_USERNAME")
	password := os.Getenv("REDIS_PASSWORD")
	db := 0

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       db,
	})

	// Ping to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}

	log.Println("✅ Redis connected:", addr)
	return rdb
}
