package dao

import (
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func getSecretOrEnv(secretName, envVar string) string {
	secretPath := filepath.Join("/run/secrets", secretName)
	data, err := os.ReadFile(secretPath)
	if err == nil {
		return strings.TrimSpace(string(data))
	}

	return os.Getenv(envVar)
}

func GetConfig() *Config {
	return &Config{
		Host:     getSecretOrEnv("PG_HOST", "PG_HOST"),
		Port:     5432,
		User:     getSecretOrEnv("PG_USER", "PG_USER"),
		Password: getSecretOrEnv("PG_PASSWORD", "PG_PASSWORD"),
		DBName:   getSecretOrEnv("PG_DBNAME", "PG_DBNAME"),
	}
}
