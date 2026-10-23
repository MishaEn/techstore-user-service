package config

import (
	"fmt"
	"os"
)

type Config struct {
	HTTPPort      string
	DBHost        string
	DBPort        string
	DBName        string
	DBUser        string
	DBPassword    string
	JWTPrivateKey string
	JWTPublicKey  string
}

func Load() (*Config, error) {
	cfg := &Config{
		HTTPPort:      getEnv("HTTP_PORT", "8080"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBName:        getEnv("DB_NAME", "users"),
		DBUser:        getEnv("DB_USER", "techstore"),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		JWTPrivateKey: getEnv("JWT_PRIVATE_KEY", ""),
		JWTPublicKey:  getEnv("JWT_PUBLIC_KEY", ""),
	}

	if cfg.DBPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable is required")
	}

	// JWT ключи обязательны только для ручек логина/регистрации,
	// но для простоты проверим их наличие сразу
	if cfg.JWTPrivateKey == "" || cfg.JWTPublicKey == "" {
		return nil, fmt.Errorf("JWT_PRIVATE_KEY and JWT_PUBLIC_KEY are required")
	}

	return cfg, nil
}

func (c *Config) DBConnString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
