package localcache

import (
	"fmt"
	"time"

	gocache "github.com/patrickmn/go-cache"
	"github.com/yiran15/api-server/base/constant"
	"github.com/yiran15/api-server/pkg/oauth"
)

type Cacher interface {
	SetCache(key string, value any, expire time.Duration)
	GetCache(key string) (any, error)
}

type Cache struct {
	cache *gocache.Cache
}

func NewCacher(oauth *oauth.OAuth2) Cacher {
	oauth2ProviderList := make([]string, 0, len(oauth.Providers))
	if oauth.Enable && len(oauth.Providers) > 0 {
		for key := range oauth.Providers {
			oauth2ProviderList = append(oauth2ProviderList, key)
		}
	}
	c := gocache.New(5*time.Minute, 10*time.Minute)
	cache := &Cache{
		cache: c,
	}

	cache.SetCache(constant.OAuth2ProviderList, oauth2ProviderList, gocache.NoExpiration)
	return cache
}

func (receive *Cache) SetCache(key string, value any, expire time.Duration) {
	receive.cache.Set(key, value, expire)
}

func (receive *Cache) GetCache(key string) (any, error) {
	item, found := receive.cache.Get(key)
	if !found {
		return nil, fmt.Errorf("cache not found: %s", key)
	}
	return item, nil
}
