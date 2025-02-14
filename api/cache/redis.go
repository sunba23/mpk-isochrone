package cache

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	rdb  *redis.Client
	once sync.Once
)

func GetClient() *redis.Client {
	password := getRedisPassword()

	once.Do(func() {
		rdb = redis.NewClient(&redis.Options{
			Addr:     getEnvOrDefault("REDIS_URL", "redis:6379"),
			Password: password,
			DB:       0,
		})

		ctx := context.Background()
		if err := rdb.Ping(ctx).Err(); err != nil {
			panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
		}
	})
	return rdb
}

// to support both env vars and docker secrets
func getRedisPassword() string {
	secretPath := "/run/secrets/REDIS_PASSWORD"

  _, err := os.Stat(secretPath)

  if os.IsNotExist(err){
		password := getEnvOrDefault("REDIS_PASSWORD", "")
		return password
	} else {
		passwordBytes, err := os.ReadFile(secretPath)
		if err != nil {
			panic(fmt.Sprintf("Failed to read Redis password from secret file: %v", err))
		}
		password := string(passwordBytes)
		return password
	}
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
