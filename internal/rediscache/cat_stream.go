package rediscache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/malkev1ch/first-task/internal/config"
	"github.com/malkev1ch/first-task/internal/model"
	"github.com/sirupsen/logrus"
	"strconv"
	"sync"
)

// CatStreamCache type represents cache object structure and behavior
type CatStreamCache struct {
	cats        map[string]*model.Cat
	redisClient *redis.Client
	streamName  string
	workersNum  int
	group       string
	mutex       sync.Mutex
}

func NewCatStreamCache(cfg *config.Config, redisClient *redis.Client) *CatStreamCache {
	CatStreamCache := &CatStreamCache{
		redisClient: redisClient,
		cats:        make(map[string]*model.Cat),
		mutex:       sync.Mutex{},
		streamName:  cfg.CatsStreamName,
		workersNum:  cfg.CacheWorkersNum,
		group:       cfg.CatsStreamGroupName,
	}
	CatStreamCache.StartWorkers(CatStreamCache.redisClient.Context())
	return CatStreamCache
}

// Get method return cat instance from cache
func (cache *CatStreamCache) Get(ctx context.Context, id string) (*model.Cat, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	cat, found := cache.cats[id]
	return cat, found
}

//Set method send message to redis stream for saving cat
func (cache *CatStreamCache) Set(ctx context.Context, input *model.Cat) error {
	catJSON, err := json.Marshal(input)
	if err != nil {
		logrus.Errorf("redis: can't marshall body of message - %e", err)
		return fmt.Errorf("redis: can't marshall body of message - %w", err)
	}
	result := cache.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: cache.streamName,
		Values: map[string]interface{}{
			"method": "set",
			"data":   catJSON,
		},
	})

	if _, err := result.Result(); err != nil {
		logrus.Errorf("redis: can't send message with method set to stream - %e", err)
		return fmt.Errorf("cache: can't send message with method set to stream - %w", err)
	}

	return nil
}

// Update method send message to redis stream for updating cat
func (cache *CatStreamCache) Update(ctx context.Context, input *model.Cat) error {
	catJSON, err := json.Marshal(input)
	if err != nil {
		logrus.Errorf("redis: can't marshall body of message - %e", err)
		return fmt.Errorf("redis: can't marshall body of message - %w", err)
	}
	result := cache.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: cache.streamName,
		Values: map[string]interface{}{
			"method": "update",
			"data":   catJSON,
		},
	})

	if _, err := result.Result(); err != nil {
		logrus.Errorf("redis: can't send message with method update to stream - %e", err)
		return fmt.Errorf("redis: can't send message with method update to stream - %w", err)
	}

	return nil
}

// Delete method removes cats from cache
func (cache *CatStreamCache) Delete(ctx context.Context, id string) error {
	catJSON, err := json.Marshal(&model.Cat{ID: id})
	if err != nil {
		logrus.Errorf("redis: can't marshall body of message - %e", err)
		return fmt.Errorf("redis: can't marshall body of message - %w", err)
	}
	result := cache.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: cache.streamName,
		Values: map[string]interface{}{
			"method": "update",
			"data":   catJSON,
		},
	})

	if _, err := result.Result(); err != nil {
		logrus.Errorf("redis: can't send message with method delete to stream - %e", err)
		return fmt.Errorf("redis: can't send message with method delete to stream - %w", err)
	}
	return nil
}

func (cache *CatStreamCache) StartWorkers(ctx context.Context) {
	_, err := cache.redisClient.XGroupCreateMkStream(ctx, cache.streamName, cache.group, "$").Result()
	if err != nil {
		logrus.Errorf("redis: error occurred while creating redis stream group - %e", err)
	}

	for i := 0; i < cache.workersNum; i++ {
		go func() {
			for {
				result, err := cache.redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
					Group:    cache.group,
					Consumer: "worker-" + strconv.Itoa(i),
					Streams:  []string{cache.streamName, ">"},
					Count:    1,
					Block:    0,
				}).Result()

				if err != nil {
					logrus.Errorf("redis: reading stream message failed- %e", err)
				}
				xMsg := result[0].Messages[0]
				msg := xMsg.Values
				msgString, ok := msg["data"].(string)
				if ok {
					cat := model.Cat{}
					err := json.Unmarshal([]byte(msgString), &cat)
					if err != nil {
						logrus.Errorf("redis: error occurred while deserialization message - %e", err)
					}
					func() {
						cache.mutex.Lock()
						defer cache.mutex.Unlock()
						switch msg["method"].(string) {
						case "set", "update":
							cache.cats[cat.ID] = &cat
							logrus.Infof("successfully cached cat - %+v", cat)
						case "delete":
							delete(cache.cats, cat.ID)
							logrus.Infof("successfully deleted cat - %+v", cat)
						default:
							logrus.Infof("redis: invalid method's name - %s", msg["method"].(string))
						}
					}()
					if _, err := cache.redisClient.XAck(ctx, cache.streamName, cache.group, xMsg.ID).Result(); err != nil {
						logrus.Errorf("redis: acknowledgement stream message failed- %e", err)
					}
				}

			}
		}()
	}
}
