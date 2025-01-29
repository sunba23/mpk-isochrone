package dao

import "os"

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func GetConfig() *Config {
	return &Config{
		Host:     "localhost",
		Port:     5432,
		User:     os.Getenv("PG_USER"),
		Password: os.Getenv("PG_PASSWORD"),
		DBName:   os.Getenv("PG_DBNAME"),
	}
}
