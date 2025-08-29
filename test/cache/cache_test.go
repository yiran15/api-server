package cache_test

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/yiran15/api-server/base/conf"
	"github.com/yiran15/api-server/base/data"
	"github.com/yiran15/api-server/base/log"
	"github.com/yiran15/api-server/store"
	"go.uber.org/zap"
)

var (
	cacheStore  store.CacheStorer
	redisClient *redis.Client
	closeup     func()
)

func init() {
	var (
		err error
	)
	conf.LoadConfig("../../config.yaml")
	redisClient, err = data.NewRDB()
	if err != nil {
		panic(err)
	}
	log.NewLogger()
	cacheStore, closeup, err = store.NewCacheStore(redisClient)
	if err != nil {
		panic(err)
	}
}

func TestCacheStore(t *testing.T) {
	defer closeup()
	var (
		roles []string
		err   error
	)
	roleNames := []any{"test"}
	if err := cacheStore.SetSet(context.Background(), store.RoleType, "test", roleNames, nil); err != nil {
		t.Fatal(err)
	}
	if roles, err = cacheStore.GetSet(context.Background(), store.RoleType, "test"); err != nil {
		t.Fatal(err)
	}
	zap.L().Info("roles", zap.Any("roles", roles))
}
