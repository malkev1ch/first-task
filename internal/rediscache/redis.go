package rediscache

import (
	"context"
	"github.com/malkev1ch/first-task/internal/config"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/malkev1ch/first-task/internal/model"
)

var (
	cacheTTL = time.Hour
)

type Cat interface {
	Get(ctx context.Context, id string) (*model.Cat, bool)
	Set(ctx context.Context, input *model.Cat) error
	Update(ctx context.Context, input *model.Cat) error
	Delete(ctx context.Context, id string) error
}

type Cache struct {
	Cat
}

// NewCache returns new cache instance with redisdb client
func NewCache(redisClient *redis.Client) *Cache {
	return &Cache{NewCatCache(redisClient, cacheTTL)}
}

// NewStreamCache returns new cache instance with redisdb client
func NewStreamCache(cfg *config.Config, redisClient *redis.Client) *Cache {
	return &Cache{NewCatStreamCache(cfg, redisClient)}
}
