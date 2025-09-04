package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Set test environment variables
	os.Setenv("SERVER_HOST", "test-host")
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("SERVER_READ_TIMEOUT", "60s")
	os.Setenv("DB_HOST", "test-db-host")
	os.Setenv("REDIS_HOST", "test-redis-host")
	os.Setenv("API_KEY", "test-api-key")

	// Clean up after test
	defer func() {
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("SERVER_READ_TIMEOUT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("REDIS_HOST")
		os.Unsetenv("API_KEY")
	}()

	cfg := LoadConfig()

	// Test server configuration
	if cfg.Server.Host != "test-host" {
		t.Errorf("Expected server host 'test-host', got '%s'", cfg.Server.Host)
	}
	if cfg.Server.Port != "9090" {
		t.Errorf("Expected server port '9090', got '%s'", cfg.Server.Port)
	}
	if cfg.Server.ReadTimeout != 60*time.Second {
		t.Errorf("Expected read timeout 60s, got %v", cfg.Server.ReadTimeout)
	}

	// Test database configuration
	if cfg.Database.Host != "test-db-host" {
		t.Errorf("Expected database host 'test-db-host', got '%s'", cfg.Database.Host)
	}

	// Test Redis configuration
	if cfg.Redis.Host != "test-redis-host" {
		t.Errorf("Expected Redis host 'test-redis-host', got '%s'", cfg.Redis.Host)
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	// Clear all environment variables
	envVars := []string{
		"SERVER_HOST", "SERVER_PORT", "SERVER_READ_TIMEOUT", "SERVER_WRITE_TIMEOUT", "SERVER_IDLE_TIMEOUT",
		"DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASSWORD",
		"REDIS_HOST", "REDIS_PORT", "REDIS_PASSWORD", "REDIS_DB",
		"API_KEY",
	}

	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}

	cfg := LoadConfig()

	// Test default values
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Expected default server host '0.0.0.0', got '%s'", cfg.Server.Host)
	}
	if cfg.Server.Port != "8080" {
		t.Errorf("Expected default server port '8080', got '%s'", cfg.Server.Port)
	}
	if cfg.Database.Host != "localhost" {
		t.Errorf("Expected default database host 'localhost', got '%s'", cfg.Database.Host)
	}
	if cfg.Redis.Host != "localhost" {
		t.Errorf("Expected default Redis host 'localhost', got '%s'", cfg.Redis.Host)
	}
}
