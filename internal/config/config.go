package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
}

// AppConfig holds application-level configuration
type AppConfig struct {
	Name     string
	Env      string
	Location string
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host          string
	Port          string
	Password      string
	Username      string
	DB            int
	InitTimeout   int
	ArticleTTL    int
	ArticlePrefix string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:     getEnv("APP_NAME", "article_api"),
			Env:      getEnv("APP_ENV", "dev"),
			Location: getEnv("SERVER_LOCATION", "Asia/Jakarta"),
		},
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnv("HTTP_SERVER_PORT", "8080"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 120*time.Second),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_DATABASE", "article_db"),
			User:     getEnv("DB_USERNAME", "article_user"),
			Password: getEnv("DB_PASSWORD", "article_password"),
		},
		Redis: RedisConfig{
			Host:          getEnv("REDIS_HOST", "localhost"),
			Port:          getEnv("REDIS_PORT", "6379"),
			Password:      getEnv("REDIS_PASSWORD", ""),
			Username:      getEnv("REDIS_USERNAME", ""),
			DB:            getIntEnv("REDIS_DB", 3),
			InitTimeout:   getIntEnv("REDIS_INIT_TIMEOUT", 5),
			ArticleTTL:    getIntEnv("REDIS_ARTICLE_TTL", 86400),
			ArticlePrefix: getEnv("REDIS_ARTICLE_PREFIX", "article-"),
		},
	}
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getIntEnv gets an integer environment variable with a fallback default value
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getDurationEnv gets a duration environment variable with a fallback default value
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
