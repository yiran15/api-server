package store

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yiran15/api-server/base/conf"
)

type CacheStorer interface {
	DelKey(ctx context.Context, cacheType CacheType, cacheKey any) error
	GetSet(ctx context.Context, cacheType CacheType, cacheKey any) ([]string, error)
	SetSet(ctx context.Context, cacheType CacheType, cacheKey any, cacheValue []any, expireTime *time.Duration) error
}

var (
	NeverExpires time.Duration = 0
)

type CacheType string

const (
	RoleType CacheType = "role"
	TestType CacheType = "test"
)

type CacheStore struct {
	client     *redis.Client
	expireTime time.Duration
	keyPrefix  string
}

func NewCacheStore(redisClient *redis.Client) (*CacheStore, func(), error) {
	expireTime, err := conf.GetRedisExpireTime()
	if err != nil {
		return nil, nil, err
	}
	closeup := func() {
		_ = redisClient.Close()
	}
	prefix, err := conf.GetRedisKeyPrefix()
	if err != nil {
		return nil, nil, err
	}
	return &CacheStore{
		client:     redisClient,
		expireTime: expireTime,
		keyPrefix:  prefix,
	}, closeup, nil
}

// NormalizeCacheKey 将常用类型的 cacheKey 转换为 string
func (c *CacheStore) NormalizeCacheKey(cacheKey any) (string, error) {
	switch v := cacheKey.(type) {
	case string:
		return v, nil
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	default:
		return "", fmt.Errorf("unsupported cacheKey type: %v", cacheKey)
	}
}

func (c *CacheStore) GetSet(ctx context.Context, cacheType CacheType, cacheKey any) ([]string, error) {
	key, err := c.NormalizeCacheKey(cacheKey)
	if err != nil {
		return nil, err
	}

	saveKey := fmt.Sprintf("%s:%v:%s", c.keyPrefix, cacheType, key)
	result, err := c.client.SMembers(ctx, saveKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("get set error: %w", err)
	}
	return result, nil
}

func (c *CacheStore) SetSet(ctx context.Context, cacheType CacheType, cacheKey any, cacheValue []any, expireTime *time.Duration) error {
	key, err := c.NormalizeCacheKey(cacheKey)
	if err != nil {
		return err
	}

	saveKey := fmt.Sprintf("%s:%v:%s", c.keyPrefix, cacheType, key)
	if err := c.client.SAdd(ctx, saveKey, cacheValue...).Err(); err != nil {
		return fmt.Errorf("redis setSet error: %w", err)
	}

	if expireTime != nil {
		if err := c.client.Expire(ctx, saveKey, *expireTime).Err(); err != nil {
			return fmt.Errorf("redis setSet error: %w", err)
		}
	}
	return nil
}

func GetExpireTime(expireTime time.Duration) *time.Duration {
	return &expireTime
}

func (c *CacheStore) DelKey(ctx context.Context, cacheType CacheType, cacheKey any) error {
	key, err := c.NormalizeCacheKey(cacheKey)
	if err != nil {
		return err
	}

	saveKey := fmt.Sprintf("%s:%v:%s", c.keyPrefix, cacheType, key)
	if err := c.client.Del(ctx, saveKey).Err(); err != nil {
		return fmt.Errorf("redis delKey error: %w", err)
	}
	return nil
}
