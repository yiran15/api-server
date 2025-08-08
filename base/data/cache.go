package data

import (
	"context"
	"fmt"
	"time"

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
		Addr:            host,
		Password:        password,
		DB:              conf.GetRedisDB(),
		PoolSize:        conf.GetRedisPoolSize(),
		MinIdleConns:    conf.GetRedisMinIdleConns(),
		ConnMaxLifetime: conf.GetRedisConnMaxLifetime(),
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
		PoolSize:         50,               // 最多50个连接
		MinIdleConns:     20,               // 最少20个空闲连接
		ConnMaxIdleTime:  10 * time.Minute, // 空闲 10 分钟关闭
		ConnMaxLifetime:  30 * time.Minute, // 强制重连以避免连接老化
	})
	err = rdb.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("redis sentinel connect failed: %w", err)
	}
	zap.S().Info("redis sentinel connect success")
	return rdb, nil
}
