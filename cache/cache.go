package cache

import (
	"encoding/json"
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	defaultExpiration         = 12 * time.Hour
	purgeTime                 = 24 * time.Hour
	ManagementTokenKey string = "managementToken"
)

type tokenCache struct {
	tokens *cache.Cache
}

var (
	tokenCacheInstance *tokenCache
)

func newCache() *tokenCache {
	Cache := cache.New(defaultExpiration, purgeTime)
	return &tokenCache{
		tokens: Cache,
	}
}

func (c *tokenCache) Read(key string) (cachedToken *string, err error) {
	token, ok := c.tokens.Get(key)
	if !ok {
		return nil, nil
	}
	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}
	var res string
	err = json.Unmarshal(tokenBytes, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *tokenCache) Update(id string, token string) {
	c.tokens.Set(id, token, cache.DefaultExpiration)
}

func Get() *tokenCache {
	if tokenCacheInstance == nil {
		tokenCacheInstance = newCache()
	}
	return tokenCacheInstance
}
