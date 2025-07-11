package data

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/yiran15/api-server/base/conf"
	"go.uber.org/zap"
)

func NewRDB() (*redis.Client, error) {
	ctx := context.TODO()
	switch conf.GetRedisMode() {
	case "sentinel":
		return initSentinelRedis(ctx)
	case "single":
		return initSingleRedis(ctx)
	default:
		return nil, fmt.Errorf("redis.mode is not supported: %s", conf.GetRedisMode())
	}
}

func initSingleRedis(ctx context.Context) (*redis.Client, error) {
	host, err := conf.GetRedisHost()
	if err != nil {
		return nil, err
	}
	password, err := conf.GetRedisPassword()
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       conf.GetRedisDB(),
	})
	err = rdb.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("redis connect failed: %w", err)
	}
	zap.S().Info("redis connect success")
	return rdb, nil
}

func initSentinelRedis(ctx context.Context) (*redis.Client, error) {
	sentinelHosts, err := conf.GetRedisSentinelHosts()
	if err != nil {
		return nil, err
	}
	masterName, err := conf.GetRedisMasterName()
	if err != nil {
		return nil, err
	}
	password, err := conf.GetRedisPassword()
	if err != nil {
		return nil, err
	}
	sentPassword, err := conf.GetRedisSentinelPassword()
	if err != nil {
		return nil, err
	}
	rdb := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:       masterName,
		SentinelAddrs:    sentinelHosts,
		Password:         password,
		SentinelPassword: sentPassword,
		RouteByLatency:   true,
		DB:               conf.GetRedisDB(),
	})
	err = rdb.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("redis sentinel connect failed: %w", err)
	}
	zap.S().Info("redis sentinel connect success")
	return rdb, nil
}
