// Package rediscache represents caching in application
package rediscache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/malkev1ch/first-task/internal/model"
	"github.com/sirupsen/logrus"
	"sync"

	"github.com/go-redis/redis/v8"
)

// CatCache type represents redis object cat structure and behavior
type CatCache struct {
	redisClient *redis.Client
	mutex       sync.Mutex
}

func NewCatCache(redisClient *redis.Client) *CatCache {
	return &CatCache{
		redisClient: redisClient,
		mutex:       sync.Mutex{},
	}
}

//Save method saves cat
func (cache *CatCache) Save(ctx context.Context, input *model.Cat) error {
	inputJSON, err := json.Marshal(input)
	if err != nil {
		logrus.Error(err, "service: error occurred while marshaling data")
		return fmt.Errorf("service: error occurred while marshaling data - %w", err)
	}
	if err := cache.redisClient.Set(ctx, input.ID, inputJSON, 0).Err(); err != nil {
		logrus.Error(err, "service: error occurred while saving cat in cache")
		return fmt.Errorf("service: error occurred while aving cat in cache - %w", err)
	}
	return nil
}

// Get method return cat instance from redis
func (cache *CatCache) Get(ctx context.Context, id string) (*model.Cat, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	val, err := cache.redisClient.Get(ctx, id).Result()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(val)
	return nil, true
}

// Update method updates cat
func (cache *CatCache) Update() error {
	return nil
}

// Delete method removes cat
func (cache *CatCache) Delete() error {
	return nil
}
