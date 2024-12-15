package services

import (
    "github.com/go-redis/redis/v8"
    "context"
)

type CacheService interface {
    GetFromCache(key string) (string, error)
    SetToCache(key, value string) error
}

// RedisCacheClient wraps *redis.Client to implement CacheService
type RedisCacheClient struct {
    Client *redis.Client
}

var RedisClient *RedisCacheClient

// InitializeCache initializes the Redis client and wraps it into RedisCacheClient
func InitializeCache() {
    RedisClient = &RedisCacheClient{
        Client: redis.NewClient(&redis.Options{
            Addr: "localhost:6379", // Redis server address
        }),
    }
}

// GetFromCache retrieves the value from the cache
func (r *RedisCacheClient) GetFromCache(key string) (string, error) {
    return r.Client.Get(context.Background(), key).Result()
}

// SetToCache sets a value to the cache
func (r *RedisCacheClient) SetToCache(key, value string) error {
    return r.Client.Set(context.Background(), key, value, 0).Err()
}
