package rediscache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/malkev1ch/first-task/internal/model"
)

type Cat interface {
	Get(ctx context.Context, id string) (*model.Cat, bool)
	Save(ctx context.Context, input *model.Cat) error
	Update() error
	Delete() error
}

type Cache struct {
	Cat
}

func NewCache(redisClient *redis.Client) *Cache {
	return &Cache{NewCatCache(redisClient)}
}
