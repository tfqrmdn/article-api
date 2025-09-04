package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"article-api/internal/config"

	"github.com/redis/go-redis/v9"
)

// CacheService handles Redis caching operations
type CacheService struct {
	client *redis.Client
	ctx    context.Context
}

// NewCacheService creates a new Redis cache service
func NewCacheService() (*CacheService, error) {
	cfg := config.LoadConfig()

	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx := context.Background()

	// Test connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &CacheService{
		client: rdb,
		ctx:    ctx,
	}, nil
}

// Set stores a value in cache with 10-minute TTL
func (c *CacheService) Set(key string, value interface{}) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	// Set with 10-minute expiration
	err = c.client.Set(c.ctx, key, jsonData, 10*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// SetWithTTL stores a value in cache with custom TTL
func (c *CacheService) SetWithTTL(key string, value interface{}, ttlSeconds int) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	// Set with custom expiration
	err = c.client.Set(c.ctx, key, jsonData, time.Duration(ttlSeconds)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// Get retrieves a value from cache
func (c *CacheService) Get(key string, dest interface{}) error {
	val, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found")
		}
		return fmt.Errorf("failed to get from cache: %w", err)
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

// Delete removes a key from cache
func (c *CacheService) Delete(key string) error {
	err := c.client.Del(c.ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete from cache: %w", err)
	}
	return nil
}

// Close closes the Redis connection
func (c *CacheService) Close() error {
	return c.client.Close()
}
