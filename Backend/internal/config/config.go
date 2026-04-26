package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port        string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	JWTSecret   string
	JWTExpHours int
	Env         string
}

func Load() *Config {
	jwtExp, _ := strconv.Atoi(getEnv("JWT_EXP_HOURS", "24"))
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "cbqa"),
		DBPassword:  getEnv("DB_PASSWORD", "cbqa123"),
		DBName:      getEnv("DB_NAME", "cbqa_db"),
		JWTSecret:   getEnv("JWT_SECRET", "cbqa-super-secret-key-change-in-production"),
		JWTExpHours: jwtExp,
		Env:         getEnv("ENV", "development"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
