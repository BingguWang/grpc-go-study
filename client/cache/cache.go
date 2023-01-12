package cache

import (
	"fmt"
	"github.com/BingguWang/grpc-go-study/server/utils"
	cache "github.com/patrickmn/go-cache"
	"sync"
)

var (
	myonce             sync.Once
	clientCounterCache *ClientCounterCache
)

type ClientCounterCache struct {
	*cache.Cache
}

const (
	succeedKey   = "succeed"
	failedKey    = "failed"
	limitedKey   = "limited"
	callTimesKey = "callTimes"
)

func NewClientCounterCache() *ClientCounterCache {
	// 初始化一个请求计数器
	myonce.Do(func() {
		if clientCounterCache == nil {
			clientCounterCache = &ClientCounterCache{
				cache.New(cache.NoExpiration, cache.DefaultExpiration),
			}
		}
		clientCounterCache.Set(succeedKey, 0, cache.NoExpiration)
		clientCounterCache.Set(failedKey, 0, cache.NoExpiration)
		clientCounterCache.Set(limitedKey, 0, cache.NoExpiration)
		clientCounterCache.Set(callTimesKey, 0, cache.NoExpiration)
	})
	return clientCounterCache
}

func (c *ClientCounterCache) IncrementSucceed(i int64) {
	c.Increment(succeedKey, i)
}
func (c *ClientCounterCache) IncrementFailed(i int64) {
	c.Increment(failedKey, i)
}
func (c *ClientCounterCache) IncrementLimitedKey(i int64) {
	c.Increment(limitedKey, i)
}
func (c *ClientCounterCache) IncrementCallTimesKey(i int64) {
	c.Increment(callTimesKey, i)
}
func (c *ClientCounterCache) CachePrint() {
	items := c.Items()
	for s, item := range items {
		fmt.Printf("key: %v, val: %v \n", s, utils.ToJsonString(item))
	}
}
