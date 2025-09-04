package cache

// CacheServiceInterface defines the contract for cache operations
type CacheServiceInterface interface {
	Set(key string, value interface{}) error
	SetWithTTL(key string, value interface{}, ttlSeconds int) error
	Get(key string, dest interface{}) error
	Delete(key string) error
	Close() error
}
