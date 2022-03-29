// Package rediscache represents caching in application
package rediscache

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/malkev1ch/first-task/internal/model"
	"github.com/sirupsen/logrus"
)

// CatCache type represents redis object cat structure and behavior.
type CatCache struct {
	redisClient      *redis.Client
	expirationPeriod time.Duration
}

func NewCatCache(redisClient *redis.Client, expirationPeriod time.Duration) *CatCache {
	return &CatCache{
		redisClient:      redisClient,
		expirationPeriod: expirationPeriod,
	}
}

// Set method saves cat in redis.
func (cache *CatCache) Set(ctx context.Context, input *model.Cat) error {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(input); err != nil {
		logrus.Error(err, "redis: error occurred encoding cat")
		return fmt.Errorf("redis: error occurred encoding cat - %w", err)
	}
	if err := cache.redisClient.Set(ctx, input.ID, b.Bytes(), cache.expirationPeriod).Err(); err != nil {
		logrus.Error(err, "redis: error occurred while saving cat in cache")
		return fmt.Errorf("redis: error occurred while saving cat in cache - %w", err)
	}

	logrus.Infof("redis: saved object cat %s in redis", input.ID)
	return nil
}

// Get method return cat instance from redis.
func (cache *CatCache) Get(ctx context.Context, id string) (*model.Cat, bool) {
	val := cache.redisClient.Get(ctx, id)

	valBytes, err := val.Bytes()
	if err != nil {
		logrus.Infof("redis: key %s doesn't exists - %e", id, err)
		return nil, false
	}

	b := bytes.NewReader(valBytes)

	var cat model.Cat

	if err := gob.NewDecoder(b).Decode(&cat); err != nil {
		logrus.Errorf("redis: error occurred while decoding cat %s - %e", id, err)
		return nil, false
	}

	logrus.Infof("redis: returned cat %s from redis", id)
	return &cat, true
}

func (cache *CatCache) Update(ctx context.Context, input *model.Cat) error {
	return nil
}

func (cache *CatCache) Delete(ctx context.Context, id string) error {
	return nil
}
