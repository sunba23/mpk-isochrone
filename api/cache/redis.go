package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"sync"
)

var (
	rdb  *redis.Client
	once sync.Once
)

func GetClient() *redis.Client {
	once.Do(func() {
		rdb = redis.NewClient(&redis.Options{
			Addr:     getEnvOrDefault("REDIS_URL", "redis:6379"),
			Password: getEnvOrDefault("REDIS_PASSWORD", ""),
			DB:       0,
		})

		ctx := context.Background()
		if err := rdb.Ping(ctx).Err(); err != nil {
			panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
		}
	})
	return rdb
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func Close() error {
	if rdb != nil {
		return rdb.Close()
	}
	return nil
}
