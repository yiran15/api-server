package store

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
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

	saveKey := c.buildCacheKey(cacheType, key)
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
	if cacheValue == nil {
		return fmt.Errorf("cacheValue cannot be nil")
	}

	key, err := c.NormalizeCacheKey(cacheKey)
	if err != nil {
		return err
	}

	saveKey := c.buildCacheKey(cacheType, key)
	if expireTime != nil {
		// 使用事务确保SADD和EXPIRE的原子性
		pipe := c.client.TxPipeline()
		pipe.SAdd(ctx, saveKey, cacheValue...)
		pipe.Expire(ctx, saveKey, *expireTime)

		if _, err := pipe.Exec(ctx); err != nil {
			return fmt.Errorf("redis setSet error: %w", err)
		}
	}

	if err := c.client.SAdd(ctx, saveKey, cacheValue...).Err(); err != nil {
		return fmt.Errorf("redis setSet error: %w", err)
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
	delKey := c.buildCacheKey(cacheType, key)
	if err := c.client.Del(ctx, delKey).Err(); err != nil {
		return fmt.Errorf("redis delKey error: %w", err)
	}
	return nil
}

// 新增辅助方法用于构建缓存key，提高可读性和可测试性
func (c *CacheStore) buildCacheKey(cacheType CacheType, key string) string {
	var sb strings.Builder
	sb.Grow(len(c.keyPrefix) + 1 + len(cacheType) + 1 + len(key))
	sb.WriteString(c.keyPrefix)
	sb.WriteByte(':')
	sb.WriteString(string(cacheType))
	sb.WriteByte(':')
	sb.WriteString(key)
	return sb.String()
}
