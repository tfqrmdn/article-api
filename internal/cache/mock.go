package cache

import "fmt"

// MockCacheService is a mock implementation of CacheService for testing
type MockCacheService struct {
	data map[string]interface{}
}

// NewMockCacheService creates a new mock cache service
func NewMockCacheService() *MockCacheService {
	return &MockCacheService{
		data: make(map[string]interface{}),
	}
}

// Set stores a value in the mock cache
func (m *MockCacheService) Set(key string, value interface{}) error {
	m.data[key] = value
	return nil
}

// SetWithTTL stores a value in the mock cache with TTL (ignored in mock)
func (m *MockCacheService) SetWithTTL(key string, value interface{}, ttlSeconds int) error {
	m.data[key] = value
	return nil
}

// Get retrieves a value from the mock cache
func (m *MockCacheService) Get(key string, dest interface{}) error {
	value, exists := m.data[key]
	if !exists {
		return fmt.Errorf("key not found")
	}

	// Simple type assertion for testing
	if destPtr, ok := dest.(*[]interface{}); ok {
		if val, ok := value.([]interface{}); ok {
			*destPtr = val
			return nil
		}
	}

	return fmt.Errorf("type assertion failed")
}

// Delete removes a key from the mock cache
func (m *MockCacheService) Delete(key string) error {
	delete(m.data, key)
	return nil
}

// Close closes the mock cache (no-op)
func (m *MockCacheService) Close() error {
	return nil
}
